package sf

import (
	"strings"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
)

// GetAdditionalRecords gathers and checks configured additionalRecords
func GetAdditionalRecords(zoneName string, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var additionalRecordNames []string

additionalRecords:
	for _, additionalRecord := range config.AdditionalRecords {
		if additionalRecord.Name == "" {
			log.Error("Additional DNS record name cannot be empty")
			continue
		}
		if !strings.Contains(additionalRecord.Name, zoneName) {
			continue
		}
		for _, userRecord := range userRecords {
			if userRecord.Name == additionalRecord.Name {
				log.Warnf("DNS record %s is already defined in a Traefik route. Skipping...", additionalRecord.Name)
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
		userRecords = append(userRecords, additionalRecord)
		additionalRecordNames = append(additionalRecordNames, additionalRecord.Name)
	}
	rootDNSRecord := cloudflare.DNSRecord{
		Type:    "A",
		Name:    zoneName,
		Content: currentIP,
		Proxied: config.Cloudflare.Defaults.Proxied,
		TTL:     config.Cloudflare.Defaults.TTL,
	}
	userRecords = append(userRecords, rootDNSRecord)
	log.Debugf("Found additional DNS records: %s", strings.Join(additionalRecordNames, ", "))

	return userRecords
}
