package sf

import (
	"io/ioutil"
	"os"

	"github.com/cloudflare/cloudflare-go"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var config *Configuration

// Configuration is a struct to store the script's configuration
type Configuration struct {
	RootDomain    string   `yaml:"rootDomain"`
	IPProviders   []string `yaml:"ipProviders"`
	Notifications struct {
		Slack struct {
			WebhookURL string `yaml:"webhookURL"`
			Username   string `yaml:"username"`
			Channel    string `yaml:"channel"`
			IconURL    string `yaml:"iconURL"`
		} `yaml:"slack"`
	} `yaml:"notifications"`
	Traefik struct {
		URL          string   `yaml:"url"`
		Username     string   `yaml:"username"`
		Password     string   `yaml:"password"`
		IgnoredRules []string `yaml:"ignoredRules"`
	} `yaml:"traefik"`
	AdditionalRecords []cloudflare.DNSRecord `yaml:"additionalRecords"`
	Cloudflare        struct {
		Email    string `yaml:"email"`
		APIKey   string `yaml:"apiKey"`
		Defaults struct {
			Type    string `yaml:"type"`
			Proxied bool   `yaml:"proxied"`
			TTL     int    `yaml:"ttl"`
		} `yaml:"defaults"`
	} `yaml:"cloudflare"`
}

// GetConfig creates a global var holding the configuration
func GetConfig(configFilePath string) {
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
	if os.Getenv("SLACK_WEBHOOK") != "" {
		config.Notifications.Slack.WebhookURL = os.Getenv("SLACK_WEBHOOK")
	}
	if os.Getenv("TRAEFIK_PASSWORD") != "" {
		config.Traefik.Password = os.Getenv("TRAEFIK_PASSWORD")
	}
	if os.Getenv("CLOUDFLARE_APIKEY") != "" {
		config.Cloudflare.APIKey = os.Getenv("CLOUDFLARE_APIKEY")
	}

	// Validate config
	if config.RootDomain == "" {
		log.Fatal("rootDomain cannot be empty")
	}
	if config.Traefik.URL == "" {
		log.Fatal("Traefik URL cannot be empty")
	}
	if config.Cloudflare.Email == "" {
		log.Fatal("Cloudflare email cannot be empty")
	}
	if config.Cloudflare.APIKey == "" {
		log.Fatal("Cloudflare api key cannot be empty")
	}

	// Set default values
	if config.Cloudflare.Defaults.Type == "" {
		config.Cloudflare.Defaults.Type = "CNAME"
	}
	if config.Cloudflare.Defaults.TTL == 0 || config.Cloudflare.Defaults.Proxied {
		config.Cloudflare.Defaults.TTL = 1
	}
	if config.IPProviders == nil {
		config.IPProviders = append(config.IPProviders, "https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip/")
	}
	if config.Notifications.Slack.Username == "" {
		config.Notifications.Slack.Username = "SyncFlaer"
	}
	if config.Notifications.Slack.IconURL == "" {
		// config.Notifications.Slack.IconURL = "https://placeholder.url.for.icon"
	}
}
