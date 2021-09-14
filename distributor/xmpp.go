package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io"

	"github.com/bdlm/log"
	"github.com/google/uuid"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
	"unifiedpush.org/go/np2p_dbus/distributor"
	"unifiedpush.org/go/np2p_dbus/storage"

	"dev.sum7.eu/genofire/unified-push-xmpp/messages"
)

type XMPPService struct {
	Login    string
	Password string
	Gateway  string
	dbus     *distributor.DBus
	session  *xmpp.Session
	store    *storage.Storage
}

func (s *XMPPService) Run(dbus *distributor.DBus, store *storage.Storage) error {
	var err error
	s.dbus = dbus
	s.store = store
	j := jid.MustParse(s.Login)
	if s.session, err = xmpp.DialClientSession(
		context.TODO(), j,
		xmpp.BindResource(),
		xmpp.StartTLS(&tls.Config{
			ServerName: j.Domain().String(),
		}),
		//TODO sasl.ScramSha1Plus <- problem with (my) ejabberd
		//xmpp.SASL("", s.Password, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		xmpp.SASL("", s.Password, sasl.ScramSha1, sasl.Plain),
	); err != nil {
		return err
	}
	defer func() {
		log.Info("closing session")
		if err := s.session.Close(); err != nil {
			log.Errorf("error closing session: %q", err)
		}
		log.Println("closing connection")
		if err := s.session.Conn().Close(); err != nil {
			log.Errorf("error closing connection: %q", err)
		}
	}()
	// Send initial presence to let the server know we want to receive messages.
	err = s.session.Send(context.TODO(), stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return err
	}
	s.session.Serve(mux.New(
		mux.MessageFunc("", xml.Name{Local: "subject"}, s.message),
	))
	return nil
}

// handler of incoming message - forward to DBUS
func (s *XMPPService) message(msgHead stanza.Message, t xmlstream.TokenReadEncoder) error {
	logger := log.WithFields(map[string]interface{}{
		"to":   msgHead.To.String(),
		"from": msgHead.From.String(),
		"id":   msgHead.ID,
	})
	d := xml.NewTokenDecoder(t)
	msg := messages.MessageBody{}
	err := d.Decode(&msg)
	if err != nil && err != io.EOF {
		log.WithField("msg", msg).Errorf("error decoding message: %q", err)
		return nil
	}
	from := msgHead.From.Bare().String()
	if s.Gateway == "" || from != s.Gateway {
		log.WithField("from", from).Info("message not from gateway, that is no notification")
		return nil
	}

	if msg.Body == "" || msg.PublicToken == "" {
		log.Infof("empty: %v", msgHead)
		return nil
	}
	logger = logger.WithFields(map[string]interface{}{
		"publicToken": msg.PublicToken,
		"content":     msg.Body,
	})

	conn := s.store.GetConnectionbyPublic(msg.PublicToken)
	if conn == nil {
		logger.Warnf("no appID and appToken found for publicToken")
	}
	logger = logger.WithFields(map[string]interface{}{
		"appID":    conn.AppID,
		"appToken": conn.AppToken,
	})
	if s.dbus.
		NewConnector(conn.AppID).
		Message(conn.AppToken, msg.Body, msgHead.ID) != nil {
		logger.Errorf("Error send unified push: %q", err)
		return nil
	}
	logger.Infof("recieve unified push")

	return nil
}

// Register handler of DBUS Distribution
func (s *XMPPService) Register(appID, appToken string) (string, string, error) {
	publicToken := uuid.New().String()
	logger := log.WithFields(map[string]interface{}{
		"appID":       appID,
		"appToken":    appToken,
		"publicToken": publicToken,
	})
	iq := messages.RegisterIQ{
		IQ: stanza.IQ{
			Type: stanza.SetIQ,
			To:   jid.MustParse(s.Gateway),
		},
	}
	iq.Register.Token = &messages.TokenData{Body: publicToken}
	t, err := s.session.EncodeIQ(context.TODO(), iq)
	if err != nil {
		logger.Errorf("xmpp send IQ for register: %v", err)
		return "", "xmpp unable send iq to gateway", err
	}
	defer func() {
		if err := t.Close(); err != nil {
			logger.Errorf("unable to close registration response %v", err)
		}
	}()
	d := xml.NewTokenDecoder(t)
	iqRegister := messages.RegisterIQ{}
	if err := d.Decode(&iqRegister); err != nil {
		logger.Errorf("xmpp recv IQ for register: %v", err)
		return "", "xmpp unable recv iq to gateway", err
	}
	if endpoint := iqRegister.Register.Endpoint; endpoint != nil {
		logger.WithField("endpoint", endpoint.Body).Info("success")
		// update Endpoint
		conn := s.store.NewConnectionWithToken(appID, appToken, publicToken, endpoint.Body)
		return conn.Endpoint, "", nil
	}
	errStr := "Unknown Error"
	if errr := iqRegister.Register.Error; errr != nil {
		errStr = errr.Body
	}
	err = errors.New(errStr)
	logger.WithField("error", err).Error("unable to register")
	return "", errStr, err
}

// Unregister handler of DBUS Distribution
func (xs *XMPPService) Unregister(token string) {
	conn, _ := xs.store.DeleteConnection(token)
	log.WithFields(map[string]interface{}{
		"appID":       conn.AppID,
		"appToken":    conn.AppToken,
		"publicToken": conn.PublicToken,
		"endpoint":    conn.Endpoint,
	}).Info("distributor-unregister")
	_ = xs.dbus.NewConnector(conn.AppID).Unregistered(conn.AppToken)
}
