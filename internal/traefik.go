package sf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/traefik/v2/pkg/rules"
)

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

func GetTraefikRules(config *Configuration, currentIP, zoneName string, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
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
			log.Fatalf("Error creating http client for Traefik %s: %s", traefikInstance.Name, err)
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
		for _, router := range traefikRouters {
			parsedDomains, err := rules.ParseDomains(router.Rule)
			if err != nil {
				log.Fatalf("Unable to parse rule %s for Traefik %s: %s", router.Rule, traefikInstance.Name, err)
			}

		parsedDomain:
			for _, parsedDomain := range parsedDomains {
				if !strings.HasSuffix(parsedDomain, zoneName) {
					continue
				}
				if parsedDomain == zoneName {
					continue
				}
				if checkDuplicateRule(parsedDomain, userRecords) {
					continue
				}
				for _, ignoredRule := range traefikInstance.IgnoredRules {
					if strings.HasSuffix(parsedDomain, ignoredRule) {
						continue parsedDomain
					}
				}
				traefikRecord := cloudflare.DNSRecord{}
				traefikRecord.Name = parsedDomain
				traefikRecord.Type = config.Cloudflare.Defaults.Type
				traefikRecord.Content = content
				traefikRecord.Proxied = config.Cloudflare.Defaults.Proxied
				traefikRecord.TTL = config.Cloudflare.Defaults.TTL

				for _, defaultOverride := range traefikInstance.DefaultOverrides {
					if defaultOverride.Rule != parsedDomain {
						continue
					}
					if defaultOverride.Type != "" {
						traefikRecord.Type = defaultOverride.Type
					}
					if defaultOverride.Content != "" {
						traefikRecord.Content = defaultOverride.Content
					}
					if defaultOverride.Proxied != nil {
						traefikRecord.Proxied = defaultOverride.Proxied
					}
					if defaultOverride.TTL != 0 {
						traefikRecord.TTL = defaultOverride.TTL
					}
				}

				userRecords = append(userRecords, traefikRecord)
				ruleNames = append(ruleNames, parsedDomain)
			}
		}
		log.Debugf("Found rules in Traefik %s: %s", traefikInstance.Name, strings.Join(ruleNames, ", "))
	}

	return userRecords
}
