package main

import (
	"flag"
	"os"
	"path/filepath"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
	"unifiedpush.org/go/np2p_dbus/distributor"
	"unifiedpush.org/go/np2p_dbus/storage"
)

var dbus *distributor.DBus

type configData struct {
	StoragePath string      `toml"storage_path"`
	XMPP        XMPPService `toml:"xmpp"`
}

func defaultPath(given, filename string) string {
	if given != "" {
		return given
	}
	basedir := os.Getenv("XDG_CONFIG_HOME")
	if len(basedir) == 0 {
		basedir = os.Getenv("HOME")
		if len(basedir) == 0 {
			basedir = "./" // FIXME: set to cwd if dunno wth is going on
		}
		basedir = filepath.Join(basedir, ".config")
	}
	basedir = filepath.Join(basedir, "unifiedpushxmpp")
	os.MkdirAll(basedir, 0o700)
	return filepath.Join(basedir, filename)
}

func main() {
	configPath := ""
	flag.StringVar(&configPath, "c", configPath, "path to configuration file")
	flag.Parse()

	config := &configData{}
	if err := file.ReadTOML(defaultPath(configPath, "config.toml"), config); err != nil {
		log.Panicf("open config file: %s", err)
	}

	store, err := storage.InitStorage(defaultPath(config.StoragePath, "database.db"))
	if err != nil {
		log.Panicf("open storage: %s", err)
	}

	dbus = distributor.NewDBus("org.unifiedpush.Distributor.xmpp")
	dbus.StartHandling(&config.XMPP)

	log.Info("startup")
	if err := config.XMPP.Run(dbus, store); err != nil {
		log.Errorf("startup xmpp: %v", err)
	}
}
