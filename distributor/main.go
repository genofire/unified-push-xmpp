package main

import (
	"os"
	"os/signal"
	"syscall"
	"errors"

	"github.com/bdlm/log"
	"unifiedpush.org/go/np2p_dbus/distributor"
)

var dbus *distributor.DBus

func main() {
	dbus = distributor.NewDBus("org.unifiedpush.Distributor.xmpp")
	dbus.StartHandling(handler{})

	log.Info("startup")

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Infof("received %s", sig)
}

type handler struct {
}

func (h handler) Register(appName, token string) (string,string,error) {
	log.WithFields(map[string]interface{}{
		"name": appName,
		"token": token,
	}).Info("distributor-register")
	endpoint := "https://up.chat.sum7.eu/UP?appid="+appName+"&token="+token
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
