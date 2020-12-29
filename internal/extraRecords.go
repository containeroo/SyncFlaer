package sf

import (
	"strings"

	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
)

// GetAdditionalRecords gathers and checks configured additionalRecords
func GetAdditionalRecords(userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var additionalRecordNames []string

	for _, additionalRecord := range config.AdditionalRecords {
		if additionalRecord.Name == "" {
			log.Error("Additional DNS record name cannot be empty")
			continue
		}
		if additionalRecord.Type == "" {
			additionalRecord.Type = config.Cloudflare.Defaults.Type
		}
		if additionalRecord.Content == "" {
			if additionalRecord.Type == "A" {
				additionalRecord.Content = currentIP
			}
			if additionalRecord.Type == "CNAME" {
				additionalRecord.Content = config.RootDomain
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
		Name:    config.RootDomain,
		Content: currentIP,
		Proxied: config.Cloudflare.Defaults.Proxied,
		TTL:     config.Cloudflare.Defaults.TTL,
	}
	userRecords = append(userRecords, rootDNSRecord)
	log.Debugf("Found additional DNS records: %s", strings.Join(additionalRecordNames, ", "))
	return userRecords
}
