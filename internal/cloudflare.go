package sf

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// CloudflareClient is a wrapper for the Cloudflare API
type CloudflareClient struct {
	client *cloudflare.API
}

func (c *CloudflareClient) DeleteGraceRecrodPrefix() string {
	return "_syncflaer._deletegrace"
}

// SetupCloudflareClient creates a Cloudflare client instance
func SetupCloudflareClient(config Configuration) *CloudflareClient {
	cf, err := cloudflare.NewWithAPIToken(config.Cloudflare.APIToken)
	if err != nil {
		log.Fatalf("Unable to setup Cloudflare client: %s", err)
	}
	cfc := CloudflareClient{client: cf}

	return &cfc
}

// CreateCloudflareZoneMap creates a map containing the zone ids
func CreateCloudflareZoneMap(config Configuration, cf *CloudflareClient) map[string]string {
	zoneIDs := make(map[string]string)
	for _, zoneName := range config.Cloudflare.ZoneNames {
		zoneID, err := cf.client.ZoneIDByName(zoneName)
		if err != nil {
			log.Fatalf("Unable to get Cloudflare zone %s id: %s", zoneName, err)
		}
		log.Debugf("Using Cloudflare DNS zone %s with id %s", zoneName, zoneID)
		zoneIDs[zoneName] = zoneID
	}

	return zoneIDs
}

// GetCloudflareDNSRecords gathers all DNS records in a given zone
func GetCloudflareDNSRecords(cf *CloudflareClient, zoneID string) []cloudflare.DNSRecord {
	dnsRecords, err := cf.client.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{})
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

// GetDeleteGraceRecords gathers all delete grace DNS records in a given zone
func GetDeleteGraceRecords(cf *CloudflareClient, zoneID string) []cloudflare.DNSRecord {
	dnsRecords, err := cf.client.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{
		Type: "TXT",
	})
	if err != nil {
		log.Fatalf("Unable to get delete grace DNS records: %s", err)
	}

	var deleteGraceRecords []cloudflare.DNSRecord
	var deleteGraceRecordNames []string

	for _, dnsRecord := range dnsRecords {
		if !strings.Contains(dnsRecord.Name, cf.DeleteGraceRecrodPrefix()) {
			continue
		}
		deleteGraceRecordNames = append(deleteGraceRecordNames, dnsRecord.Name)
		deleteGraceRecords = append(deleteGraceRecords, dnsRecord)
	}
	if deleteGraceRecordNames != nil {
		log.Debugf("Found delete grace DNS records: %s", strings.Join(deleteGraceRecordNames, " ,"))
	}
	return deleteGraceRecords
}

// CreateCloudflareDNSRecord is a wrapper function to create a DNS record
func CreateCloudflareDNSRecord(cf *CloudflareClient, zoneID string, record cloudflare.DNSRecord, slackHandler *SlackHandler) {
	_, err := cf.client.CreateDNSRecord(context.Background(), zoneID, record)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to create DNS record %s: %s", record.Name, err)
		slackHandler.AddSlackMessage(errMsg, "danger")
		log.Error(errMsg)
		return
	}

	infoMsg := fmt.Sprintf("Created: name: %s, type: %s, content: %s, proxied: %t, ttl: %d", record.Name, record.Type, record.Content, *record.Proxied, record.TTL)
	slackHandler.AddSlackMessage(infoMsg, "good")
	log.Info(infoMsg)
}

// DeleteCloudflareDNSRecord is a wrapper function to delete a DNS record
func DeleteCloudflareDNSRecord(cf *CloudflareClient, zoneID string, record cloudflare.DNSRecord, slackHandler *SlackHandler) {
	err := cf.client.DeleteDNSRecord(context.Background(), zoneID, record.ID)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to delete DNS record %s: %s", record.Name, err)
		slackHandler.AddSlackMessage(errMsg, "danger")
		log.Error(errMsg)
		return
	}

	infoMsg := fmt.Sprintf("Deleted: %s", record.Name)
	if record.Type != "TXT" {
		slackHandler.AddSlackMessage(infoMsg, "good")
		log.Info(infoMsg)
		return
	}
	log.Debug(infoMsg)
}

// UpdateCloudflareDNSRecords updates the public IP and additionalRecords
func UpdateCloudflareDNSRecords(cf *CloudflareClient, zoneID string, cloudflareDNSRecords, userRecords []cloudflare.DNSRecord, slackHandler *SlackHandler) {
	for _, dnsRecord := range cloudflareDNSRecords {
		for _, userRecord := range userRecords {
			if dnsRecord.Name != userRecord.Name {
				continue
			}
			if *dnsRecord.Proxied == *userRecord.Proxied && dnsRecord.TTL == userRecord.TTL && dnsRecord.Content == userRecord.Content {
				continue
			}
			updatedDNSRecord := cloudflare.DNSRecord{
				Type:    userRecord.Type,
				Content: userRecord.Content,
				Proxied: userRecord.Proxied,
				TTL:     userRecord.TTL,
			}
			err := cf.client.UpdateDNSRecord(context.Background(), zoneID, dnsRecord.ID, updatedDNSRecord)
			if err != nil {
				errMsg := fmt.Sprintf("Unable to update DNS record %s: %s", dnsRecord.Name, err)
				slackHandler.AddSlackMessage(errMsg, "danger")
				log.Error(errMsg)
				continue
			}

			infoMsg := fmt.Sprintf("Updated: name: %s, type: %s, content: %s, proxied: %t, ttl: %d", dnsRecord.Name, updatedDNSRecord.Type, updatedDNSRecord.Content, *updatedDNSRecord.Proxied, updatedDNSRecord.TTL)
			slackHandler.AddSlackMessage(infoMsg, "good")
			log.Info(infoMsg)
		}
	}
}

// GetMissingDNSRecords compares Cloudflare DNS records with Traefik rules and additionalRecords
func GetMissingDNSRecords(cloudflareDNSRecords, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var missingRecords []cloudflare.DNSRecord

	for _, userRecord := range userRecords {
		recordFound := false
		for _, cloudflareDNSRecord := range cloudflareDNSRecords {
			if userRecord.Name == cloudflareDNSRecord.Name {
				recordFound = true
				break
			}
		}
		if !recordFound {
			missingRecords = append(missingRecords, userRecord)
		}
	}

	return missingRecords
}

// GetOrphanedDNSRecords compares Cloudflare DNS records with Traefik rules and additionalRecords
func GetOrphanedDNSRecords(cloudflareDNSRecords, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var orphanedRecords []cloudflare.DNSRecord

	for _, cloudflareDNSRecord := range cloudflareDNSRecords {
		recordFound := false
		for _, userRecord := range userRecords {
			if userRecord.Name == cloudflareDNSRecord.Name {
				recordFound = true
				break
			}
		}
		if !recordFound {
			orphanedRecords = append(orphanedRecords, cloudflareDNSRecord)
		}
	}

	return orphanedRecords
}

func GetDeleteGraceRecord(cf *CloudflareClient, orphanedRecordName string, deleteGraceRecords []cloudflare.DNSRecord) cloudflare.DNSRecord {
	var deleteGraceRecordFound cloudflare.DNSRecord
	for _, deleteGraceRecord := range deleteGraceRecords {
		if deleteGraceRecord.Name == fmt.Sprintf("%s.%s", cf.DeleteGraceRecrodPrefix(), orphanedRecordName) {
			deleteGraceRecordFound = deleteGraceRecord
			break
		}
	}
	return deleteGraceRecordFound
}

func UpdateDeleteGraceRecord(cf *CloudflareClient, zoneID string, deleteGraceRecord cloudflare.DNSRecord, orphanedRecordName string) {
	newDeleteGrace, _ := strconv.Atoi(deleteGraceRecord.Content)
	newDeleteGrace--
	deleteGraceRecord.Content = strconv.Itoa(newDeleteGrace)
	err := cf.client.UpdateDNSRecord(context.Background(), zoneID, deleteGraceRecord.ID, deleteGraceRecord)
	if err != nil {
		log.Errorf("Unable to update delete grace DNS record %s: %s", deleteGraceRecord.Name, err)
		return
	}

	log.Infof("Waiting %s more runs until DNS record %s gets deleted", deleteGraceRecord.Content, orphanedRecordName)
}

func CleanupDeleteGraceRecords(cf *CloudflareClient, zoneID string, userRecords, cloudflareDNSRecords, deleteGraceRecords []cloudflare.DNSRecord, slackHandler *SlackHandler) {
	for _, deleteGraceRecord := range deleteGraceRecords {
		dnsRecordFound := false
		var dnsRecordName string
		for _, userRecord := range userRecords {
			if deleteGraceRecord.Name == fmt.Sprintf("%s.%s", cf.DeleteGraceRecrodPrefix(), userRecord.Name) {
				dnsRecordFound = true
				dnsRecordName = userRecord.Name
				break
			}
		}
		if dnsRecordFound {
			DeleteCloudflareDNSRecord(cf, zoneID, deleteGraceRecord, slackHandler)
			log.Infof("DNS record %s is not orphaned anymore", dnsRecordName)
			continue
		}
		for _, cloudflareDNSRecord := range cloudflareDNSRecords {
			if deleteGraceRecord.Name == fmt.Sprintf("%s.%s", cf.DeleteGraceRecrodPrefix(), cloudflareDNSRecord.Name) {
				dnsRecordFound = true
				break
			}
		}
		if !dnsRecordFound {
			DeleteCloudflareDNSRecord(cf, zoneID, deleteGraceRecord, slackHandler)
			log.Debugf("Cleaned up delete grace DNS record %s", deleteGraceRecord.Name)
		}
	}
}
