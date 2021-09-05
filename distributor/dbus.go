package main

import (
	"github.com/godbus/dbus/v5"
	"github.com/bdlm/log"
)


type Distributor struct {
	dbus *dbus.Conn
}

func NewDistributor(dbus *dbus.Conn) *Distributor {
	return &Distributor{
		dbus: dbus,
	}
}

func (d Distributor) Register(name, token string) (thing, reason string, err *dbus.Error) {
	logger:= log.WithFields(map[string]interface{}{
		"name": name,
		"token": token,
	})

	endpoint := "https://up.chat.sum7.eu/UP?appid="+name+"&token="+token

	c := NewConector(d.dbus, name)
	if err := c.NewEndpoint(token, endpoint); err != nil {
		logger.Errorf("distributor-register error on NewEndpoint: %v", err)
		return "REGISTRATION_FAILED", err.Error(), nil
	}

	logger.Info("distributor-register")
	return "NEW_ENDPOINT", "", nil
}

func (d Distributor) Unregister(token string) *dbus.Error {
	log.WithFields(map[string]interface{}{
		"token": token,
	}).Info("distributor-unregister")
	return nil
}

type Connector struct {
	obj dbus.BusObject
}

func NewConector(dbus *dbus.Conn,appid string) *Connector {
	obj := dbus.Object(appid, DBUSConnectorPath)
	return &Connector{
		obj: obj,
	}
}

func (c Connector) Message(token, contents, id string) error {
	return c.obj.Call(DBUSConnectorInterface+".Message", dbus.FlagNoReplyExpected, token, contents, id).Err
}

func (c Connector) NewEndpoint(token, endpoint string) error {
	return c.obj.Call(DBUSConnectorInterface+".NewEndpoint", dbus.FlagNoReplyExpected, token, endpoint).Err
}

func (c Connector) Unregistered(token string) error {
	return c.obj.Call(DBUSConnectorInterface+".Unregistered", dbus.FlagNoReplyExpected, token).Err
}
