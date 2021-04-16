package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"os"

	internal "github.com/containeroo/syncflaer/internal"

	log "github.com/sirupsen/logrus"
)

const version string = "4.0.0"

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
	zoneIDs := internal.CreateCloudflareZoneMap()

	internal.GetCurrentIP()

	for zoneName, zoneID := range zoneIDs {
		cloudflareDNSRecords := internal.GetCloudflareDNSRecords(zoneID)
		deleteGraceRecords := internal.GetDeleteGraceRecords(zoneID)

		var userRecords []cloudflare.DNSRecord
		userRecords = internal.GetTraefikRules(zoneName, userRecords)
		userRecords = internal.GetAdditionalRecords(zoneName, userRecords)

		missingRecords := internal.GetMissingDNSRecords(cloudflareDNSRecords, userRecords)
		if missingRecords != nil {
			for _, missingRecord := range missingRecords {
				internal.CreateCloudflareDNSRecord(zoneID, missingRecord)
			}
		} else {
			log.Debug("No missing DNS records")
		}

		orphanedRecords := internal.GetOrphanedDNSRecords(cloudflareDNSRecords, userRecords)
		if orphanedRecords != nil {
			for _, orphanedRecord := range orphanedRecords {
				if config.Cloudflare.DeleteGrace != 0 {
					deleteGraceRecord := internal.GetDeleteGraceRecord(orphanedRecord.Name, deleteGraceRecords)
					if deleteGraceRecord.Name == "" {
						internal.CreateDeleteGraceRecord(zoneID, orphanedRecord.Name)
						continue
					}
					if deleteGraceRecord.Content != "1" {
						internal.UpdateDeleteGraceRecord(zoneID, deleteGraceRecord, orphanedRecord.Name)
						continue
					}
					internal.DeleteCloudflareDNSRecord(zoneID, deleteGraceRecord)
				}
				internal.DeleteCloudflareDNSRecord(zoneID, orphanedRecord)
			}
		} else {
			log.Debug("No orphaned DNS records")
		}

		internal.CleanupDeleteGraceRecords(zoneName, zoneID, userRecords, cloudflareDNSRecords, deleteGraceRecords)

		internal.UpdateCloudflareDNSRecords(zoneID, cloudflareDNSRecords, userRecords)
	}

	internal.SendSlackMessage()
}
