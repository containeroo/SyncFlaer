package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"os"

	internal "github.com/containeroo/syncflaer/internal"

	log "github.com/sirupsen/logrus"
)

const version string = "2.1.0"

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

	cloudflareDNSRecords := internal.GetCloudflareDNSRecords()
	deleteGraceRecords := internal.GetDeleteGraceRecords()

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
			if config.Cloudflare.DeleteGrace != 0 {
				deleteGraceRecord := internal.GetDeleteGraceRecord(orphanedRecord.Name, deleteGraceRecords)
				if deleteGraceRecord.Name == "" {
					internal.CreateDeleteGraceRecord(orphanedRecord.Name)
					continue
				}
				if deleteGraceRecord.Content != "1" {
					internal.UpdateDeleteGraceRecord(deleteGraceRecord, orphanedRecord.Name)
					continue
				}
				internal.DeleteCloudflareDNSRecord(deleteGraceRecord)
			}
			internal.DeleteCloudflareDNSRecord(orphanedRecord)
		}
	} else {
		log.Debug("No orphaned DNS records")
	}

	internal.CleanupDeleteGraceRecords(userRecords, cloudflareDNSRecords, deleteGraceRecords)

	internal.UpdateCloudflareDNSRecords(cloudflareDNSRecords, userRecords)

	internal.SendSlackMessage()
}
