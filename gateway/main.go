package main

import (
	"flag"

	"dev.sum7.eu/genofire/golang-lib/file"
	"dev.sum7.eu/genofire/golang-lib/web"
	apiStatus "dev.sum7.eu/genofire/golang-lib/web/api/status"
	webM "dev.sum7.eu/genofire/golang-lib/web/metrics"
	"github.com/bdlm/log"
)

var VERSION = "development"

type configData struct {
	EndpointURL string `toml:"endpoint_url"`
	XMPP      XMPPService `toml:"xmpp"`
	Webserver web.Service `toml:"webserver"`
}

func main() {
	webM.VERSION = VERSION
	apiStatus.VERSION = VERSION
	webM.NAMESPACE = "unifiedpush"

	configPath := "config.toml"
	showVersion := false

	flag.StringVar(&configPath, "c", configPath, "path to configuration file")
	flag.BoolVar(&showVersion, "version", showVersion, "show current version")

	flag.Parse()

	if showVersion {
		log.WithField("version", VERSION).Info("Version")
		return
	}

	config := &configData{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panicf("open config file: %s", err)
	}
	// just for more beautiful config file - jere
	config.XMPP.EndpointURL = config.EndpointURL

	go func() {
		if err := config.XMPP.Run(); err != nil {
			log.Panicf("startup xmpp: %v", err)
		}
	}()

	config.Webserver.ModuleRegister(Bind(&config.XMPP))

	log.Info("startup")
	if err := config.Webserver.Run(); err != nil {
		log.Fatal(err)
	}
}
