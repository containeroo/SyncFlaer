package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"os"
	"strconv"
	"strings"

	internal "github.com/containeroo/syncflaer/internal"

	log "github.com/sirupsen/logrus"
)

const version string = "2.0.0-rc4"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	configFilePath, printVersion, debug := internal.ParseFlags()

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugf("SyncFlaer %s", version)

	config := internal.GetConfig(configFilePath)

	internal.SetupCloudflareClient()
	internal.GetCloudflareZoneID()

	internal.GetCurrentIP()

	cloudflareDNSRecords, deleteGraceRecords := internal.GetCloudflareDNSRecords()

	var userRecords []cloudflare.DNSRecord
	userRecords = internal.GetTraefikRules(userRecords)
	userRecords = internal.GetAdditionalRecords(userRecords)

	missingRecords := internal.GetMissingDNSRecords(cloudflareDNSRecords, userRecords)
	if missingRecords != nil {
		for _, missingRecord := range missingRecords {
			internal.CreateCloudflareDNSRecord(missingRecord)
		}
	} else {
		log.Debug("No missing DNS records")
	}

	orphanedRecords := internal.GetOrphanedDNSRecords(cloudflareDNSRecords, userRecords)
	if orphanedRecords != nil {
		for _, orphanedRecord := range orphanedRecords {
			if config.Cloudflare.DeleteGrace == 0 {
				internal.DeleteCloudflareDNSRecord(orphanedRecord)
				for _, deleteGraceRecord := range deleteGraceRecords {
					log.Infof("Cleaning up delete grace DNS record for %s", strings.TrimPrefix(deleteGraceRecord.Name, "_syncflaer._deletegrace."))
					internal.DeleteCloudflareDNSRecord(deleteGraceRecord)
				}
				continue
			}
			deleteGraceRecordFound := false
			for _, deleteGraceRecord := range deleteGraceRecords {
				if !strings.Contains(deleteGraceRecord.Name, orphanedRecord.Name) {
					continue
				}
				deleteGraceRecordFound = true
				deleteGrace, _ := strconv.Atoi(deleteGraceRecord.Content)
				deleteGrace -= 1
				if deleteGrace > 0 {
					deleteGraceRecord.Content = strconv.Itoa(deleteGrace)
					log.Infof("Waiting %d more runs until deleting DNS record %s", deleteGrace-1, orphanedRecord.Name)
					internal.UpdateCloudflareDNSRecord(deleteGraceRecord)
					continue
				}
				internal.DeleteCloudflareDNSRecord(orphanedRecord)
				internal.DeleteCloudflareDNSRecord(deleteGraceRecord)
			}
			if !deleteGraceRecordFound {
				deleteGraceRecord := cloudflare.DNSRecord{
					Type:    "TXT",
					Name:    fmt.Sprintf("_syncflaer._deletegrace.%s", orphanedRecord.Name),
					Content: strconv.Itoa(config.Cloudflare.DeleteGrace),
				}
				log.Infof("Waiting %d more runs until deleting DNS record %s", config.Cloudflare.DeleteGrace-1, orphanedRecord.Name)
				internal.CreateCloudflareDNSRecord(deleteGraceRecord)
			}
		}
	} else {
		for _, deleteGraceRecord := range deleteGraceRecords {
			log.Infof("DNS record %s is not orphaned anymore", strings.TrimPrefix(deleteGraceRecord.Name, "_syncflaer._deletegrace."))
			internal.DeleteCloudflareDNSRecord(deleteGraceRecord)
		}
		log.Debug("No orphaned DNS records")
	}

	internal.UpdateOutdatedDNSRecords(cloudflareDNSRecords, userRecords)

	internal.SendSlackMessage()
}
