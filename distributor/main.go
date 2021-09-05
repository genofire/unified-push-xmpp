package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bdlm/log"
	"github.com/godbus/dbus/v5"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Panicf("on dbus connection: %v", err)
	}
	defer conn.Close()
	log.Debug("connect to dbus")

	rp, err := conn.RequestName(DBUSName, dbus.NameFlagReplaceExisting)
	if err != nil {
		log.Panicf("register name on dbus: %v", err)
	}
	if rp != dbus.RequestNameReplyPrimaryOwner {
		log.Panicf("a other dbus service is running with %s: %v", DBUSName, rp)
	}
	log.Debug("register name to dbus")

	d := NewDistributor(conn)
	if err := conn.ExportAll(d, DBUSDistributorPath, DBUSDistributorInterface); err != nil {
		log.Panicf("export distributor on %s: %v", DBUSDistributorPath, err)
	}

	log.Info("startup")

	// Wait for INT/TERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Infof("received %s", sig)
}
