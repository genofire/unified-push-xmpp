package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"io"

	"github.com/bdlm/log"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
	"unifiedpush.org/go/np2p_dbus/distributor"
)

type XMPPService struct {
	Login    string
	Password string
	Gateway  string
	dbus     *distributor.DBus
}

func (xs *XMPPService) Run(dbus *distributor.DBus) error {
	xs.dbus = dbus
	j := jid.MustParse(xs.Login)
	s, err := xmpp.DialClientSession(
		context.TODO(), j,
		xmpp.BindResource(),
		xmpp.StartTLS(&tls.Config{
			ServerName: j.Domain().String(),
		}),
		//TODO sasl.ScramSha1Plus <- problem with (my) ejabberd
		//xmpp.SASL("", xs.Password, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		xmpp.SASL("", xs.Password, sasl.ScramSha1, sasl.Plain),
	)
	if err != nil {
		return err
	}
	defer func() {
		log.Info("Closing session…")
		if err := s.Close(); err != nil {
			log.Errorf("Error closing session: %q", err)
		}
		log.Println("Closing conn…")
		if err := s.Conn().Close(); err != nil {
			log.Errorf("Error closing connection: %q", err)
		}
	}()
	// Send initial presence to let the server know we want to receive messages.
	err = s.Send(context.TODO(), stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return err
	}
	s.Serve(mux.New(
		mux.MessageFunc("",xml.Name{Local: "subject"}, xs.message),
	))
	return nil
}

func (xs *XMPPService) message(msgHead stanza.Message, t xmlstream.TokenReadEncoder) error {
	d := xml.NewTokenDecoder(t)
	msg := struct {
		Token string `xml:"subject"`
		Body string `xml:"body"`
	}{}
	err := d.Decode(&msg)
	if err != nil && err != io.EOF {
		log.WithField("msg", msg).Errorf("Error decoding message: %q", err)
		return nil
	}
	from := msgHead.From.Bare().String()
	if xs.Gateway == "" || from != xs.Gateway {
		log.WithField("from", from).Info("message not from gateway, that is no notification")
		return nil
	}

	if msg.Body == "" || msg.Token == "" {
		log.Infof("empty: %v", msgHead)
		return nil
	}

	//TODO Lockup for appid by token in storage
	if xs.dbus.
		NewConnector("cc.malhotra.karmanyaah.testapp.golibrary").
		Message(msg.Token, msg.Body, msgHead.ID) != nil {
		log.Errorf("Error send unified push: %q", err)
		return nil
	}
	log.Infof("recieve unified push: %v", msg)

	return nil
}
