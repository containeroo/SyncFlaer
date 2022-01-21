package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/containeroo/syncflaer/internal/kube"
	"os"
	"strconv"

	internal "github.com/containeroo/syncflaer/internal"

	log "github.com/sirupsen/logrus"
)

const version string = "5.3.0"

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
		log.Warn("Debug mode enabled! Sensitive data could be displayed in plain text!")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugf("SyncFlaer %s", version)

	slackHandler := internal.NewSlackHandler()

	config := internal.GetConfig(configFilePath)

	cf := internal.SetupCloudflareClient(&config.Cloudflare.APIToken)
	zoneIDs := internal.CreateCloudflareZoneMap(&config.Cloudflare.ZoneNames, cf)
	currentIP := internal.GetCurrentIP(&config.IPProviders)

	for zoneName, zoneID := range zoneIDs {
		cloudflareDNSRecords := internal.GetCloudflareDNSRecords(cf, zoneID)
		deleteGraceRecords := internal.GetDeleteGraceRecords(cf, zoneID)

		var userRecords []cloudflare.DNSRecord
		if config.TraefikInstances != nil {
			userRecords = internal.GetTraefikRules(config, currentIP, zoneName, userRecords)
		}
		if config.AdditionalRecords != nil {
			userRecords = internal.GetAdditionalRecords(config, currentIP, zoneName, userRecords)
		}
		if *config.Kubernetes.Enabled {
			kubeClient := kube.SetupKubernetesClient()
			userRecords = kube.GetIngresses(kubeClient, config, currentIP, zoneName, userRecords)
		}

		missingRecords := internal.GetMissingDNSRecords(cloudflareDNSRecords, userRecords)
		if missingRecords != nil {
			for _, missingRecord := range missingRecords {
				internal.CreateCloudflareDNSRecord(cf, zoneID, missingRecord, slackHandler)
			}
		} else {
			log.Debug("No missing DNS records")
		}

		orphanedRecords := internal.GetOrphanedDNSRecords(cloudflareDNSRecords, userRecords)
		if orphanedRecords != nil {
			for _, orphanedRecord := range orphanedRecords {
				if config.Cloudflare.DeleteGrace == 0 {
					internal.DeleteCloudflareDNSRecord(cf, zoneID, orphanedRecord, slackHandler)
					continue
				}

				existingDeleteGraceRecord := internal.GetDeleteGraceRecord(cf, orphanedRecord.Name, deleteGraceRecords)
				if existingDeleteGraceRecord.Name == "" {
					falseVar := false
					deleteGraceRecord := cloudflare.DNSRecord{
						Type:    "TXT",
						Name:    fmt.Sprintf("%s.%s", cf.DeleteGraceRecordPrefix(), orphanedRecord.Name),
						Content: strconv.Itoa(config.Cloudflare.DeleteGrace),
						Proxied: &falseVar,
					}
					internal.CreateCloudflareDNSRecord(cf, zoneID, deleteGraceRecord, slackHandler)
					continue
				}

				deleteGraceRecordContent, _ := strconv.Atoi(existingDeleteGraceRecord.Content)
				if deleteGraceRecordContent > 1 {
					internal.UpdateDeleteGraceRecord(cf, zoneID, existingDeleteGraceRecord, orphanedRecord.Name)
					continue
				}

				internal.DeleteCloudflareDNSRecord(cf, zoneID, orphanedRecord, slackHandler)
				internal.DeleteCloudflareDNSRecord(cf, zoneID, existingDeleteGraceRecord, slackHandler)
			}
		} else {
			log.Debug("No orphaned DNS records")
		}

		internal.CleanupDeleteGraceRecords(cf, zoneID, userRecords, cloudflareDNSRecords, deleteGraceRecords, slackHandler)

		internal.UpdateCloudflareDNSRecords(cf, zoneID, cloudflareDNSRecords, userRecords, slackHandler)
	}

	slackHandler.SendSlackMessage(config)
}
