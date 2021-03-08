package sf

import (
	"fmt"
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
		Name         string   `yaml:"name"`
		URL          string   `yaml:"url"`
		Username     string   `yaml:"username"`
		Password     string   `yaml:"password"`
		IgnoredRules []string `yaml:"ignoredRules"`
	} `yaml:"traefikInstances"`
	AdditionalRecords []cloudflare.DNSRecord `yaml:"additionalRecords"`
	Cloudflare        struct {
		Email       string `yaml:"email"`
		APIKey      string `yaml:"apiKey"`
		DeleteGrace int    `yaml:"deleteGrace"`
		ZoneName    string `yaml:"zoneName"`
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

	// Check if env vars are set
	for i, traefikInstance := range config.TraefikInstances {
		envName := fmt.Sprintf("TRAEFIK_%s_PASSWORD", strings.ToUpper(traefikInstance.Name))
		if os.Getenv(envName) != "" {
			config.TraefikInstances[i].Password = os.Getenv(envName)
		}
	}
	if os.Getenv("SLACK_WEBHOOK") != "" {
		config.Notifications.Slack.WebhookURL = os.Getenv("SLACK_WEBHOOK")
	}
	if os.Getenv("CLOUDFLARE_APIKEY") != "" {
		config.Cloudflare.APIKey = os.Getenv("CLOUDFLARE_APIKEY")
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
		config.IPProviders = append(config.IPProviders, "https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip/")
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
	if config.Cloudflare.Email == "" {
		log.Fatal("Cloudflare email cannot be empty")
	}
	if config.Cloudflare.APIKey == "" {
		log.Fatal("Cloudflare api key cannot be empty")
	}
	if config.Cloudflare.ZoneName == "" {
		log.Fatal("Cloudflare zone name cannot be empty")
	}
	if config.Cloudflare.Defaults.Type != "A" && config.Cloudflare.Defaults.Type != "CNAME" {
		log.Fatalf("Supported Cloudflare default types are A or CNAME")
	}

	return config
}
