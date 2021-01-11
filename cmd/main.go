package main

import (
	"github.com/cloudflare/cloudflare-go"
	"os"

	internal "github.com/containeroo/syncflaer/internal"

	log "github.com/sirupsen/logrus"
)

const version string = "1.0.3"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	configFilePath, debug := internal.ParseFlags()

	if !debug {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("SyncFlaer %s", version)

	internal.GetConfig(configFilePath)

	internal.SetupCloudflareClient()
	internal.GetCloudflareZoneID()

	internal.GetCurrentIP()

	cloudflareDNSRecords := internal.GetCloudflareDNSRecords()

	var userRecords []cloudflare.DNSRecord
	userRecords = internal.GetTraefikRules(userRecords)
	userRecords = internal.GetAdditionalRecords(userRecords)

	missingRecords := internal.CheckMissingDNSRecords(cloudflareDNSRecords, userRecords)
	if missingRecords != nil {
		for _, missingRecord := range missingRecords {
			internal.CreateCloudflareDNSRecord(missingRecord)
		}
	} else {
		log.Debug("No missing DNS records")
	}

	orphanedRecords := internal.CheckOrphanedDNSRecords(cloudflareDNSRecords, userRecords)
	if orphanedRecords != nil {
		for _, orphanedRecord := range orphanedRecords {
			internal.DeleteCloudflareDNSRecord(orphanedRecord)
		}
	} else {
		log.Debug("No orphaned DNS records")
	}

	internal.UpdateCloudflareDNSRecords(cloudflareDNSRecords, userRecords)

	internal.SendSlackMessage()

	log.Debug("Done")
}
