package sf

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	SkipUpdateCheck *bool    `yaml:"skipUpdateCheck"`
	IPProviders     []string `yaml:"ipProviders"`
	Notifications   struct {
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
		DefaultOverrides     []struct {
			Rule    string `yaml:"rule"`
			Type    string `yaml:"type"`
			Content string `yaml:"content"`
			Proxied *bool  `yaml:"proxied"`
			TTL     int    `yaml:"ttl"`
		} `yaml:"defaultOverrides"`
	} `yaml:"traefikInstances"`
	Kubernetes struct {
		Enabled *bool `yaml:"enabled"`
	} `yaml:"kubernetes"`
	ManagedRootRecord *bool                  `yaml:"managedRootRecord"`
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

func maskValue(value string) string {
	rs := []rune(value)
	for i := 0; i < len(rs)-3; i++ {
		rs[i] = '*'
	}
	return string(rs)
}

func GetConfig(configFilePath string) *Configuration {
	log.Debugf("Loading config file %s", configFilePath)
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Unable to load config file %s from disk: %s", configFilePath, err)
	}

	var config Configuration
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}

	// Check if environment variables are used
	var envVarName string

	for i, traefikInstance := range config.TraefikInstances {
		if strings.HasPrefix(traefikInstance.Password, "env:") {
			envVarName = strings.Replace(traefikInstance.Password, "env:", "", 1)
			config.TraefikInstances[i].Password = os.Getenv(envVarName)
			log.Debugf("Got Traefik %s password '%s' from environment variable '%s'", traefikInstance.Name, maskValue(config.TraefikInstances[i].Password), envVarName)
		}

		for k, v := range traefikInstance.CustomRequestHeaders {
			if strings.HasPrefix(v, "env:") {
				envVarName = strings.Replace(v, "env:", "", 1)
				config.TraefikInstances[i].CustomRequestHeaders[k] = os.Getenv(envVarName)
				log.Debugf("Got Traefik %s customRequestHeader '%s' from environment variable '%s' with value '%s'", traefikInstance.Name, k, envVarName, maskValue(config.TraefikInstances[i].CustomRequestHeaders[k]))
			}
		}
	}

	if strings.HasPrefix(config.Notifications.Slack.WebhookURL, "env:") {
		envVarName = strings.Replace(config.Notifications.Slack.WebhookURL, "env:", "", 1)
		config.Notifications.Slack.WebhookURL = os.Getenv(envVarName)
		log.Debugf("Got Slack webhook URL '%s' from environment variable '%s'", maskValue(config.Notifications.Slack.WebhookURL), envVarName)
	}

	if strings.HasPrefix(config.Cloudflare.APIToken, "env:") {
		envVarName = strings.Replace(config.Cloudflare.APIToken, "env:", "", 1)
		config.Cloudflare.APIToken = os.Getenv(envVarName)
		log.Debugf("Got Cloudflare API token '%s' from environment variable '%s'", maskValue(config.Cloudflare.APIToken), envVarName)
	}

	// Set default values
	trueVar := true
	falseVar := false
	if config.ManagedRootRecord == nil {
		config.ManagedRootRecord = &trueVar
		log.Debugf("ManagedRootRecord is not set, defaulting to %t", *config.ManagedRootRecord)
	}

	if config.Kubernetes.Enabled == nil {
		falseVar := false
		config.Kubernetes.Enabled = &falseVar
		log.Debugf("Kubernetes enabled is not set, defaulting to %t", *config.Kubernetes.Enabled)
	}

	if config.Cloudflare.Defaults.Type == "" {
		config.Cloudflare.Defaults.Type = "CNAME"
		log.Debugf("Cloudflare default type is empty, defaulting to %s", config.Cloudflare.Defaults.Type)
	}

	if config.Cloudflare.Defaults.Proxied == nil {
		config.Cloudflare.Defaults.Proxied = &trueVar
		log.Debugf("Cloudflare default proxied is empty, defaulting to %t", *config.Cloudflare.Defaults.Proxied)
	}

	if config.Cloudflare.Defaults.TTL == 0 || *config.Cloudflare.Defaults.Proxied {
		config.Cloudflare.Defaults.TTL = 1
		log.Debugf("Cloudflare default TTL is empty, defaulting to %d", config.Cloudflare.Defaults.TTL)
	}

	if config.IPProviders == nil {
		config.IPProviders = append(config.IPProviders, "https://ifconfig.me/ip", "https://ipecho.net/plain", "https://myip.is/ip/", "https://checkip.amazonaws.com", "https://api.ipify.org")
		log.Debugf("IP providers is empty, defaulting to %s", strings.Join(config.IPProviders, ", "))
	}

	if config.SkipUpdateCheck == nil {
		config.SkipUpdateCheck = &falseVar
		log.Debugf("SkipUpdateCheck is not set, defaulting to %t", *config.SkipUpdateCheck)
	}

	if config.Notifications.Slack.Username == "" {
		config.Notifications.Slack.Username = "SyncFlaer"
		log.Debugf("Slack username is empty, defaulting to %s", config.Notifications.Slack.Username)
	}

	if config.Notifications.Slack.IconURL == "" {
		config.Notifications.Slack.IconURL = "https://www.cloudflare.com/img/cf-facebook-card.png"
		log.Debugf("Slack icon URL is empty, defaulting to %s", config.Notifications.Slack.IconURL)
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
		log.Fatal("Supported Cloudflare default types are A or CNAME")
	}

	return &config
}
