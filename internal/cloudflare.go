package sf

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"strings"
)

var cf *cloudflare.API
var zoneID string

// SetupCloudflareClient creates a global var with a Cloudflare client instance
func SetupCloudflareClient() {
	var err error
	cf, err = cloudflare.New(config.Cloudflare.APIKey, config.Cloudflare.Email)
	if err != nil {
		log.Fatalf("Unable to setup Cloudflare client: %s", err)
	}
}

// GetCloudflareZoneID creates a global var containing the zone id
func GetCloudflareZoneID() {
	var err error
	zoneID, err = cf.ZoneIDByName(config.Cloudflare.ZoneName)
	if err != nil {
		log.Fatalf("Unable to get Cloudflare zone id: %s", err)
	}
	log.Debugf("Using Cloudflare DNS zone with id %s", zoneID)
}

// GetCloudflareDNSRecords gathers all DNS records in a given zone
func GetCloudflareDNSRecords() ([]cloudflare.DNSRecord, []cloudflare.DNSRecord) {
	dnsRecords, err := cf.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatalf("Unable to get Cloudflare DNS records: %s", err)
	}

	var cloudflareDNSRecords []cloudflare.DNSRecord
	var deleteGraceRecords []cloudflare.DNSRecord
	var cloudflareDNSRecordNames []string

	for _, dnsRecord := range dnsRecords {
		if dnsRecord.Type == "TXT" && strings.Contains(dnsRecord.Name, "_syncflaer._deletegrace") {
			deleteGraceRecords = append(deleteGraceRecords, dnsRecord)
			continue
		}
		if dnsRecord.Type != "CNAME" && dnsRecord.Type != "A" {
			continue
		}
		cloudflareDNSRecords = append(cloudflareDNSRecords, dnsRecord)
		cloudflareDNSRecordNames = append(cloudflareDNSRecordNames, dnsRecord.Name)
	}
	log.Debugf("Found Cloudflare DNS records: %s", strings.Join(cloudflareDNSRecordNames, ", "))

	return cloudflareDNSRecords, deleteGraceRecords
}

// CreateCloudflareDNSRecord is a wrapper function to create a DNS record
func CreateCloudflareDNSRecord(record cloudflare.DNSRecord) {
	_, err := cf.CreateDNSRecord(zoneID, record)
	if err != nil {
		var errMsg string
		if record.Type == "A" || record.Type == "CNAME" {
			errMsg = fmt.Sprintf("Unable to create DNS record %s: %s", record.Name, err)
		}
		if record.Type == "TXT" {
			errMsg = fmt.Sprintf("Unable to create delete grace DNS record %s: %s", record.Name, err)
		}
		addSlackMessage(errMsg, "danger")
		log.Error(errMsg)
		return
	}

	infoMsg := fmt.Sprintf("Created: name: %s, type: %s, content: %s, proxied: %t, ttl: %d", record.Name, record.Type, record.Content, record.Proxied, record.TTL)
	if record.Type == "A" || record.Type == "CNAME" {
		addSlackMessage(infoMsg, "good")
		log.Info(infoMsg)
	}
	if record.Type == "TXT" {
		log.Debug(infoMsg)
	}
}

// DeleteCloudflareDNSRecord is a wrapper function to delete a DNS record
func DeleteCloudflareDNSRecord(record cloudflare.DNSRecord) {
	err := cf.DeleteDNSRecord(zoneID, record.ID)
	if err != nil {
		var errMsg string
		if record.Type == "A" || record.Type == "CNAME" {
			errMsg = fmt.Sprintf("Unable to delete DNS record %s: %s", record.Name, err)
		}
		if record.Type == "TXT" {
			errMsg = fmt.Sprintf("Unable to delete delete grace DNS record %s: %s", record.Name, err)
		}
		addSlackMessage(errMsg, "danger")
		log.Error(errMsg)
		return
	}

	infoMsg := fmt.Sprintf("Deleted: %s", record.Name)
	if record.Type == "A" || record.Type == "CNAME" {
		addSlackMessage(infoMsg, "good")
		log.Info(infoMsg)
	}
	if record.Type == "TXT" {
		log.Debug(infoMsg)
	}
}

// UpdateCloudflareDNSRecord is a wrapper function to update a DNS record
func UpdateCloudflareDNSRecord(record cloudflare.DNSRecord) {
	err := cf.UpdateDNSRecord(zoneID, record.ID, record)
	if err != nil {
		var errMsg string
		if record.Type == "A" || record.Type == "CNAME" {
			errMsg = fmt.Sprintf("Unable to update DNS record %s: %s", record.Name, err)
		}
		if record.Type == "TXT" {
			errMsg = fmt.Sprintf("Unable to update delete grace DNS record %s: %s", record.Name, err)
		}
		addSlackMessage(errMsg, "danger")
		log.Error(errMsg)
		return
	}

	infoMsg := fmt.Sprintf("Updated: name: %s, type: %s, content: %s, proxied: %t, ttl: %d", record.Name, record.Type, record.Content, record.Proxied, record.TTL)
	if record.Type == "A" || record.Type == "CNAME" {
		addSlackMessage(infoMsg, "good")
		log.Info(infoMsg)
	}
	if record.Type == "TXT" {
		log.Debug(infoMsg)
	}
}

// GetMissingDNSRecords compares Cloudflare DNS records with Traefik rules and additionalRecords
func GetMissingDNSRecords(cloudflareDNSRecords []cloudflare.DNSRecord, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var missingRecords []cloudflare.DNSRecord

	for _, userRecord := range userRecords {
		recordFound := false
		for _, cloudflareDNSRecord := range cloudflareDNSRecords {
			if userRecord.Name == cloudflareDNSRecord.Name {
				recordFound = true
			}
		}
		if !recordFound {
			missingRecords = append(missingRecords, userRecord)
		}
	}

	return missingRecords
}

// GetOrphanedDNSRecords compares Cloudflare DNS records with Traefik rules and additionalRecords
func GetOrphanedDNSRecords(cloudflareDNSRecords []cloudflare.DNSRecord, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var orphanedRecords []cloudflare.DNSRecord

	for _, cloudflareDNSRecord := range cloudflareDNSRecords {
		recordFound := false
		for _, userRecord := range userRecords {
			if userRecord.Name == cloudflareDNSRecord.Name {
				recordFound = true
			}
		}
		if !recordFound {
			orphanedRecords = append(orphanedRecords, cloudflareDNSRecord)
		}
	}

	return orphanedRecords
}

// UpdateOutdatedDNSRecords checks whether the DNS records must be updated
func UpdateOutdatedDNSRecords(cloudflareDNSRecords []cloudflare.DNSRecord, userRecords []cloudflare.DNSRecord) {
	for _, dnsRecord := range cloudflareDNSRecords {
		for _, userRecord := range userRecords {
			if dnsRecord.Name != userRecord.Name {
				continue
			}
			if dnsRecord.Proxied == userRecord.Proxied && dnsRecord.TTL == userRecord.TTL && dnsRecord.Content == userRecord.Content {
				continue
			}
			dnsRecord.Type = userRecord.Type
			dnsRecord.Content = userRecord.Content
			dnsRecord.Proxied = userRecord.Proxied
			dnsRecord.TTL = userRecord.TTL

			UpdateCloudflareDNSRecord(dnsRecord)
		}
	}
}
