package sf

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var config Configuration

// Configuration struct holds SyncFlaer configuration
type Configuration struct {
	IPProviders   []string `yaml:"ipProviders"`
	Notifications struct {
		Slack struct {
			WebhookURL string `yaml:"webhookURL"`
			Username   string `yaml:"username"`
			Channel    string `yaml:"channel"`
			IconURL    string `yaml:"iconURL"`
		} `yaml:"slack"`
	} `yaml:"notifications"`
	TraefikInstances []struct {
		Name                 string            `yaml:"name"`
		URL                  string            `yaml:"url"`
		Username             string            `yaml:"username"`
		Password             string            `yaml:"password"`
		CustomRequestHeaders map[string]string `yaml:"customRequestHeaders"`
		IgnoredRules         []string          `yaml:"ignoredRules"`
	} `yaml:"traefikInstances"`
	AdditionalRecords []cloudflare.DNSRecord `yaml:"additionalRecords"`
	Cloudflare        struct {
		APIToken    string   `yaml:"apiToken"`
		DeleteGrace int      `yaml:"deleteGrace"`
		ZoneNames   []string `yaml:"zoneNames"`
		Defaults    struct {
			Type    string `yaml:"type"`
			Proxied *bool  `yaml:"proxied"`
			TTL     int    `yaml:"ttl"`
		} `yaml:"defaults"`
	} `yaml:"cloudflare"`
}

// GetConfig creates a global var holding the configuration
func GetConfig(configFilePath string) Configuration {
	log.Debugf("Loading config file %s", configFilePath)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Unable to load config file %s from disk: %s", configFilePath, err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	// Check if env vars are used
	for i, traefikInstance := range config.TraefikInstances {
		if strings.HasPrefix(traefikInstance.Password, "env:") {
			config.TraefikInstances[i].Password = os.Getenv(strings.Replace(traefikInstance.Password, "env:", "", 1))
		}

		for k, v := range traefikInstance.CustomRequestHeaders {
			if strings.HasPrefix(v, "env:") {
				config.TraefikInstances[i].CustomRequestHeaders[k] = os.Getenv(strings.Replace(v, "env:", "", 1))
			}
		}
	}

	if strings.HasPrefix(config.Notifications.Slack.WebhookURL, "env:") {
		config.Notifications.Slack.WebhookURL = os.Getenv(strings.Replace(config.Notifications.Slack.WebhookURL, "env:", "", 1))
	}

	if strings.HasPrefix(config.Cloudflare.APIToken, "env:") {
		config.Cloudflare.APIToken = os.Getenv(strings.Replace(config.Cloudflare.APIToken, "env:", "", 1))
	}

	// Set default values
	if config.Cloudflare.Defaults.Type == "" {
		config.Cloudflare.Defaults.Type = "CNAME"
		log.Debug("Cloudflare default type is empty, defaulting to ", config.Cloudflare.Defaults.Type)
	}

	if config.Cloudflare.Defaults.TTL == 0 || *config.Cloudflare.Defaults.Proxied {
		config.Cloudflare.Defaults.TTL = 1
		log.Debug("Cloudflare default TTL is empty, defaulting to ", config.Cloudflare.Defaults.TTL)
	}

	if config.IPProviders == nil {
		config.IPProviders = append(config.IPProviders, "https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip/", "https://checkip.amazonaws.com")
		log.Debug("IP providers is empty, defaulting to ", strings.Join(config.IPProviders, ", "))
	}

	if config.Notifications.Slack.Username == "" {
		config.Notifications.Slack.Username = "SyncFlaer"
		log.Debug("Slack username is empty, defaulting to ", config.Notifications.Slack.Username)
	}

	if config.Notifications.Slack.IconURL == "" {
		config.Notifications.Slack.IconURL = "https://www.cloudflare.com/img/cf-facebook-card.png"
		log.Debug("Slack icon URL is empty, defaulting to ", config.Notifications.Slack.IconURL)
	}

	// Validate config
	for _, traefikInstance := range config.TraefikInstances {
		if traefikInstance.Name == "" {
			log.Fatal("Traefik instance name cannot be empty")
		}

		if traefikInstance.URL == "" {
			log.Fatalf("Traefik URL for instance %s cannot be empty", traefikInstance.Name)
		}
	}

	if config.Cloudflare.APIToken == "" {
		log.Fatal("Cloudflare API token cannot be empty")
	}

	if config.Cloudflare.ZoneNames == nil {
		log.Fatal("Cloudflare zone name cannot be empty")
	}

	if config.Cloudflare.Defaults.Type != "A" && config.Cloudflare.Defaults.Type != "CNAME" {
		log.Fatalf("Supported Cloudflare default types are A or CNAME")
	}

	return config
}
