package sf

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var currentIP string

// GetCurrentIP checks the current public IP
func GetCurrentIP() {
	rand.Seed(time.Now().UnixNano())
	providers := config.IPProviders
	rand.Shuffle(len(config.IPProviders), func(i, j int) { providers[i], providers[j] = providers[j], providers[i] })

	for _, provider := range providers {
		resp, err := http.Get(provider)
		if err != nil {
			log.Errorf("Unable to get public ip: %s", err)
			continue
		}
		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Unable to get public ip: %s", err)
			continue
		}
		currentIP = strings.TrimSpace(string(ip))
		log.Debugf("Got public ip %s from provider %s", currentIP, provider)
		break
	}
	if currentIP == "" {
		log.Fatal("Unable to get public ip from any provider")
	}
}
