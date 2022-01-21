package sf

import (
	"strings"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
)

func GetAdditionalRecords(config *Configuration, currentIP, zoneName string, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var additionalRecordNames []string

additionalRecords:
	for _, additionalRecord := range config.AdditionalRecords {
		if additionalRecord.Name == "" {
			log.Error("Additional DNS record name cannot be empty")
			continue
		}
		if !strings.HasSuffix(additionalRecord.Name, zoneName) {
			continue
		}
		for _, userRecord := range userRecords {
			if userRecord.Name == additionalRecord.Name {
				log.Warnf("DNS record %s is already defined in a Traefik route. Skipping...", userRecord.Name)
				continue additionalRecords
			}
		}
		if additionalRecord.Type == "" {
			additionalRecord.Type = config.Cloudflare.Defaults.Type
		}
		if additionalRecord.Content == "" {
			switch additionalRecord.Type {
			case "A":
				additionalRecord.Content = currentIP
			case "CNAME":
				additionalRecord.Content = zoneName
			default:
				log.Errorf("%s is an unsupported type, only A or CNAME are supported", additionalRecord.Type)
				continue
			}
		}
		if additionalRecord.TTL == 0 {
			additionalRecord.TTL = 1
		}
		if additionalRecord.Proxied == nil {
			additionalRecord.Proxied = config.Cloudflare.Defaults.Proxied
		}
		userRecords = append(userRecords, additionalRecord)
		additionalRecordNames = append(additionalRecordNames, additionalRecord.Name)
	}
	if *config.ManagedRootRecord {
		rootDNSRecord := cloudflare.DNSRecord{
			Type:    "A",
			Name:    zoneName,
			Content: currentIP,
			Proxied: config.Cloudflare.Defaults.Proxied,
			TTL:     config.Cloudflare.Defaults.TTL,
		}
		userRecords = append(userRecords, rootDNSRecord)
	}

	log.Debugf("Found additional DNS records: %s", strings.Join(additionalRecordNames, ", "))

	return userRecords
}
