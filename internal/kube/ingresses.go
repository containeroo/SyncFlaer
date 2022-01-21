package kube

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	internal "github.com/containeroo/syncflaer/internal"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func GetIngresses(kubeClient kubernetes.Interface) *v1.IngressList {
	ingresses, err := kubeClient.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Errorf("Error getting ingresses: %s", err)
	}

	return ingresses
}

func BuildCloudflareDNSRecordsFromIngresses(config *internal.Configuration, currentIP string, ingresses *v1.IngressList, zoneName string, userRecords []cloudflare.DNSRecord) []cloudflare.DNSRecord {
	var ingressNames []string
	for _, ingress := range ingresses.Items {
		if ingress.Annotations["syncflaer.containeroo.ch/ignore"] == "true" {
			log.Debugf("Ignoring ingress %s/%s", ingress.Namespace, ingress.Name)
			continue
		}

	rules:
		for _, rule := range ingress.Spec.Rules {
			if !strings.HasSuffix(rule.Host, zoneName) {
				continue
			}
			for _, userRecord := range userRecords {
				if userRecord.Name == rule.Host {
					log.Warnf("DNS record %s already defined elsewhere (%s/%s). Skipping...", userRecord.Name, ingress.Namespace, ingress.Name)
					continue rules
				}
			}
			dnsRecord := cloudflare.DNSRecord{
				Name: rule.Host,
			}
			if ingress.Annotations["syncflaer.containeroo.ch/type"] == "" {
				dnsRecord.Type = config.Cloudflare.Defaults.Type
			}
			if ingress.Annotations["syncflaer.containeroo.ch/content"] == "" {
				switch dnsRecord.Type {
				case "A":
					dnsRecord.Content = currentIP
				case "CNAME":
					dnsRecord.Content = zoneName
				default:
					log.Errorf("%s is an unsupported type, only A and CNAME are supported", ingress.Annotations["syncflaer.containeroo.ch/type"])
					continue
				}
			}
			if ingress.Annotations["syncflaer.containeroo.ch/ttl"] == "" {
				dnsRecord.TTL = config.Cloudflare.Defaults.TTL
			}
			if ingress.Annotations["syncflaer.containeroo.ch/proxied"] == "" {
				dnsRecord.Proxied = config.Cloudflare.Defaults.Proxied
			}
			userRecords = append(userRecords, dnsRecord)
			ingressNames = append(ingressNames, rule.Host)
		}
	}

	log.Debugf("Found Kubernetes ingresses: %s", strings.Join(ingressNames, ", "))

	return userRecords
}
