package sf

import (
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetCurrentIP(ipProviders *[]string) string {
	rand.Seed(time.Now().UnixNano())
	providers := *ipProviders
	rand.Shuffle(len(providers), func(i, j int) { providers[i], providers[j] = providers[j], providers[i] })

	var success bool
	var currentIP string
	for _, provider := range providers {
		success = false
		resp, err := http.Get(provider)
		if err != nil {
			log.Errorf("Unable to get public ip from %s: %s", provider, err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Errorf("Unable to get public ip from %s: http status code %d", provider, resp.StatusCode)
			continue
		}
		ip, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Unable to get public ip from %s: %s", provider, err)
			continue
		}
		currentIP = strings.TrimSpace(string(ip))
		if net.ParseIP(currentIP) == nil {
			log.Errorf("Public ip %s from %s is invalid", currentIP, provider)
			continue
		}
		log.Debugf("Got public ip %s from provider %s", currentIP, provider)
		success = true
		break
	}
	if !success {
		log.Fatal("Unable to get public ip from any configured provider")
	}

	return currentIP
}
