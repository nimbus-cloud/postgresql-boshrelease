package main

import (
	"flag"
	"log"
	"github.com/nimbus-cloud/dnsupdater/config"
	"github.com/nimbus-cloud/dnsupdater/dns"

	"github.com/cloudfoundry/bosh-utils/logger"

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

	level, err := logger.Levelify(cfg.Logging.Level)
	if err != nil {
		panic(err.Error())
	}

	logger := logger.NewLogger(level)
	dnsUpdater := dns.NewDnsUpdater(*cfg, logger)

	go dnsUpdater.Run()

	// error chan not being used by anything, just to prevent main from exiting
	errCh := make(chan error, 1)
	select {
	case <-errCh:
		return
	}
}