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
	zoneID, err = cf.ZoneIDByName(config.RootDomain)
	if err != nil {
		log.Fatalf("Unable to get Cloudflare zone id: %s", err)
	}
	log.Debugf("Using Cloudflare DNS zone with id %s", zoneID)
}

// GetCloudflareDNSRecords gathers all DNS records in a given zone
func GetCloudflareDNSRecords() []cloudflare.DNSRecord {
	dnsRecords, err := cf.DNSRecords(zoneID, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatalf("Unable to get Cloudflare DNS records: %s", err)
	}

	var cloudflareDNSRecords []cloudflare.DNSRecord
	var cloudflareDNSRecordNames []string

	for _, dnsRecord := range dnsRecords {
		if dnsRecord.Type != "CNAME" && dnsRecord.Type != "A" {
			continue
		}
		cloudflareDNSRecords = append(cloudflareDNSRecords, dnsRecord)
		cloudflareDNSRecordNames = append(cloudflareDNSRecordNames, dnsRecord.Name)
	}
	log.Debugf("Found Cloudflare DNS records: %s", strings.Join(cloudflareDNSRecordNames, ", "))
	return cloudflareDNSRecords
}

// CreateCloudflareDNSRecord is a wrapper function to create a DNS record
func CreateCloudflareDNSRecord(record cloudflare.DNSRecord) {
	newDNSRecord := cloudflare.DNSRecord{
		Type:    record.Type,
		Name:    record.Name,
		Content: record.Content,
		Proxied: record.Proxied,
		TTL:     record.TTL,
	}

	_, err := cf.CreateDNSRecord(zoneID, newDNSRecord)
	if err != nil {
		addSlackMessage(fmt.Sprintf("Unable to create DNS record %s: %s", newDNSRecord.Name, err), "danger")
		log.Errorf("Unable to create DNS record %s: %s", newDNSRecord.Name, err)
		return
	}

	addSlackMessage(fmt.Sprintf("Created: %s IN %s %s, proxied %t, ttl %d", newDNSRecord.Name, newDNSRecord.Type, newDNSRecord.Content, newDNSRecord.Proxied, newDNSRecord.TTL), "good")
	log.Infof("Created: %s IN %s %s, proxied %t, ttl %d", newDNSRecord.Name, newDNSRecord.Type, newDNSRecord.Content, newDNSRecord.Proxied, newDNSRecord.TTL)
}

// DeleteCloudflareDNSRecord is a wrapper function to delete a DNS record
func DeleteCloudflareDNSRecord(record cloudflare.DNSRecord) {
	err := cf.DeleteDNSRecord(zoneID, record.ID)
	if err != nil {
		addSlackMessage(fmt.Sprintf("Unable to delete DNS record %s: %s", record.Name, err), "danger")
		log.Errorf("Unable to delete DNS record %s: %s", record.Name, err)
		return
	}

	addSlackMessage(fmt.Sprintf("Deleted: %s", record.Name), "good")
	log.Infof("Deleted: %s", record.Name)
}

// UpdateCloudflareDNSRecords updates the public IP and additionalRecords
func UpdateCloudflareDNSRecords(cloudflareDNSRecords []cloudflare.DNSRecord, userRecords []cloudflare.DNSRecord) {
	for _, dnsRecord := range cloudflareDNSRecords {
		for _, userRecord := range userRecords {
			if dnsRecord.Name != userRecord.Name {
				continue
			}
			if dnsRecord.Proxied == userRecord.Proxied && dnsRecord.TTL == userRecord.TTL && dnsRecord.Content == userRecord.Content {
				continue
			}
			updatedDNSRecord := cloudflare.DNSRecord{
				Type:    userRecord.Type,
				Content: userRecord.Content,
				Proxied: userRecord.Proxied,
				TTL:     userRecord.TTL,
			}
			err := cf.UpdateDNSRecord(zoneID, dnsRecord.ID, updatedDNSRecord)
			if err != nil {
				addSlackMessage(fmt.Sprintf("Unable to update DNS record %s: %s", dnsRecord.Name, err), "danger")
				log.Errorf("Unable to update DNS record %s: %s", dnsRecord.Name, err)
				continue
			}
			addSlackMessage(fmt.Sprintf("Updated: %s IN %s %s, proxied %t, ttl %d", dnsRecord.Name, updatedDNSRecord.Type, updatedDNSRecord.Content, updatedDNSRecord.Proxied, updatedDNSRecord.TTL), "good")
			log.Infof("Updated: %s IN %s %s, proxied %t, ttl %d", dnsRecord.Name, updatedDNSRecord.Type, updatedDNSRecord.Content, updatedDNSRecord.Proxied, updatedDNSRecord.TTL)
		}
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
