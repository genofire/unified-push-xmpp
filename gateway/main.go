package main

import (
	"flag"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
)

func main() {
	configPath := "config.toml"
	flag.StringVar(&configPath, "c", configPath, "path to configuration file")
	flag.Parse()

	config := &XMPPService{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panicf("open config file: %s", err)
	}

	log.Info("startup")
	if err := config.Run(); err != nil {
		log.Errorf("startup xmpp: %v", err)
	}
}
