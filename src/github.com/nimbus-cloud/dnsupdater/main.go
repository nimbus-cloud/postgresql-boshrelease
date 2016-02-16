package main

import (
	"flag"
	"log"
	"github.com/nimbus-cloud/dnsupdater/config"

	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/nimbus-cloud/dnsupdater/dns"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.Parse()
}

func main() {
	var cfg *config.Config

	if configFile != "" {
		log.Printf("Loading file: %s\n", configFile)
		cfg = config.InitConfigFromFile(configFile)
	} else {
		log.Fatal("Config file not specified.")
	}

	logger := logger.NewLogger(logger.LevelInfo)
	dnsUpdater := dns.NewDnsUpdater(*cfg, logger)

	go dnsUpdater.Run()

	// error chan not being used by anything, just to prevent main from exiting
	errCh := make(chan error, 1)
	select {
	case <-errCh:
		return
	}
}