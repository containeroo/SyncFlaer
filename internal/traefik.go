package sf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
)

var traefikRegex = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)

// TraefikRouter is a struct to store a router object of Traefik
type TraefikRouter struct {
	EntryPoints []string `json:"entryPoints"`
	Middlewares []string `json:"middlewares,omitempty"`
	Service     string   `json:"service"`
	Rule        string   `json:"rule"`
	TLS         struct {
		CertResolver string `json:"certResolver"`
		Domains      []struct {
			Main string   `json:"main"`
			Sans []string `json:"sans"`
		} `json:"domains"`
	} `json:"tls,omitempty"`
	Status   string   `json:"status"`
	Using    []string `json:"using"`
	Name     string   `json:"name"`
	Provider string   `json:"provider"`
	Priority int64    `json:"priority,omitempty"`
}

func checkDuplicateRule(rule string, rules []cloudflare.DNSRecord) bool {
	for _, r := range rules {
		if r.Name == rule {
			return true
		}
	}
	return false
}

// GetTraefikRules gathers and formats all Traefik http routers
func GetTraefikRules(config Configuration, currentIP, zoneName string, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	for _, traefikInstance := range config.TraefikInstances {
		traefikURL, err := url.Parse(traefikInstance.URL)
		if err != nil {
			log.Fatalf("Unable to parse Traefik url %s: %s", traefikInstance.URL, err)
		}
		traefikURL.Path = path.Join(traefikURL.Path, "/api/http/routers")
		traefikHost := traefikURL.String()

		client := &http.Client{}
		req, err := http.NewRequest("GET", traefikHost, nil)
		if err != nil {
			log.Fatalf("Error creating http request for Traefik %s: %s", traefikInstance.Name, err)
		}
		if traefikInstance.Username != "" && traefikInstance.Password != "" {
			req.SetBasicAuth(traefikInstance.Username, traefikInstance.Password)
		}
		for k, v := range traefikInstance.CustomRequestHeaders {
			req.Header.Add(k, v)
			log.Debugf("Adding request header to Traefik %s: '%s: %s'", traefikInstance.Name, k, maskValue(v))
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Unable to get Traefik %s rules: %s", traefikInstance.Name, err)
		}
		if resp.StatusCode != 200 {
			log.Fatalf("Unable to get Traefik %s rules: http status code %d", traefikInstance.Name, resp.StatusCode)
		}

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Unable to read Traefik %s rules: %s", traefikInstance.Name, err)
		}

		var traefikRouters []TraefikRouter
		err = json.Unmarshal(respData, &traefikRouters)
		if err != nil {
			log.Fatalf("Unable to load Traefik %s rules: %s", traefikInstance.Name, err)
		}

		var content string
		switch config.Cloudflare.Defaults.Type {
		case "A":
			content = currentIP
		case "CNAME":
			content = zoneName
		}

		var ruleNames []string
	rules:
		for _, router := range traefikRouters {
			if !strings.Contains(router.Rule, zoneName) {
				continue
			}
			matches := traefikRegex.FindAllStringSubmatch(router.Rule, -1)
			for _, match := range matches {
				if match[0] == zoneName {
					continue
				}
				if !checkDuplicateRule(match[0], userRecords) {
					for _, ignoredRule := range traefikInstance.IgnoredRules {
						if strings.Contains(match[0], ignoredRule) {
							continue rules
						}
					}
					userRecords = append(userRecords, cloudflare.DNSRecord{
						Type:    config.Cloudflare.Defaults.Type,
						Name:    match[0],
						Content: content,
						Proxied: config.Cloudflare.Defaults.Proxied,
						TTL:     config.Cloudflare.Defaults.TTL,
					})
					ruleNames = append(ruleNames, match[0])
				}
			}
		}
		log.Debugf("Found rules in Traefik %s: %s", traefikInstance.Name, strings.Join(ruleNames, ", "))
	}

	return userRecords
}
