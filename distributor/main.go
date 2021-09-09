package main

import (
	"errors"
	"flag"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
	"unifiedpush.org/go/np2p_dbus/distributor"
)

var dbus *distributor.DBus

func main() {
	configPath := "config.toml"
	flag.StringVar(&configPath, "c", configPath, "path to configuration file")
	flag.Parse()

	config := &XMPPService{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panicf("open config file: %s", err)
	}

	dbus = distributor.NewDBus("org.unifiedpush.Distributor.xmpp")
	dbus.StartHandling(handler{})

	log.Info("startup")
	if err := config.Run(dbus); err != nil {
		log.Errorf("startup xmpp: %v", err)
	}
}

type handler struct {
}

func (h handler) Register(appName, token string) (string, string, error) {
	log.WithFields(map[string]interface{}{
		"name":  appName,
		"token": token,
	}).Info("distributor-register")
	endpoint := "https://up.chat.sum7.eu/UP?appid=" + appName + "&token=" + token
	if endpoint != "" {
		return endpoint, "", nil
	}
	return "", "reason to app", errors.New("Unknown error")
}
func (h handler) Unregister(token string) {
	log.WithFields(map[string]interface{}{
		"token": token,
	}).Info("distributor-unregister")
	appID := ""
	_ = dbus.NewConnector(appID).Unregistered(token)
}
