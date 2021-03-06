package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/bdlm/log"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/disco"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
	"unifiedpush.org/go/np2p_dbus/distributor"
	"unifiedpush.org/go/np2p_dbus/storage"

	"dev.sum7.eu/genofire/unified-push-xmpp/messages"
)

// Demo Server as fallback
var (
	XMPPUPDemoJID = jid.MustParse("up.chat.sum7.eu")
)

type XMPPService struct {
	Login       string `toml:"login"`
	Password    string `toml:"password"`
	Gateway     string `toml:"gateway"`
	KeepGateway bool   `toml:"keep_gateway"`
	gateway     jid.JID
	dbus        *distributor.DBus
	session     *xmpp.Session
	store       *storage.Storage
}

func (s *XMPPService) Run(dbus *distributor.DBus, store *storage.Storage) error {
	var err error
	s.dbus = dbus
	s.store = store
	j := jid.MustParse(s.Login)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	if s.session, err = xmpp.DialClientSession(
		ctx, j,
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
	err = s.session.Send(ctx, stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return err
	}
	go s.checkServer()
	go s.selectGateway()
	log.Debug("xmpp client is running")
	s.session.Serve(mux.New(
		// disco.Handle(),
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
		return nil
	}
	from := msgHead.From.String()
	if settings := strings.Split(conn.Settings, ":"); len(settings) > 1 && settings[0] == from {
		log.WithField("from", from).Info("message not from gateway, that is no notification")
		return nil
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
	logger.Infof("receive unified push")

	return nil
}

// checkServer - background job
func (s *XMPPService) checkServer() {
	domain := s.session.LocalAddr().Domain()
	logger := log.WithField("instance", domain.String())
	logger.Debug("check running")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*500))
	defer cancel()
	info, err := disco.GetInfo(ctx, "", domain, s.session)
	if err != nil {
		log.Errorf("check server: %v", err)
		return
	}

	// check if server support msgoffline
	supportMSGOffline := false
	for _, f := range info.Features {
		if f.Var == "msgoffline" {
			supportMSGOffline = true
			break
		}
	}
	if !supportMSGOffline {
		log.Warn("your server does not support offline messages (XEP-0160) - it is need to deliever messages later, if this distributer has current no connection")
	}
	logger.Info("your instance checked")
	return
}

// selectGateway - background job
func (s *XMPPService) selectGateway() {
	if gateway, err := jid.Parse(s.Gateway); err != nil {
		if err := s.findGateway(); err != nil {
			log.Panicf("no gateway found: %v", err)
		} else {
			log.WithField("gateway", s.gateway.String()).Info("using found UnifiedPush")
		}
	} else {
		if err := s.testAndUseGateway(gateway); err != nil {
			log.Panic(err)
		} else {
			log.WithField("gateway", s.gateway.String()).Info("using configured UnifiedPush")
		}
	}
	// function to renew endpoint if new gateway was detected
	if s.KeepGateway {
		return
	}
	conns := s.store.GetUnequalSettings(s.gateway.String() + ":" + s.session.LocalAddr().Bare().String())
	if len(conns) <= 0 {
		return
	}
	log.WithField("count", len(conns)).Info("register apps for new gateway")
	for _, i := range conns {
		s.Register(i.AppID, i.AppToken)
	}
}

// findGateway
func (s *XMPPService) findGateway() error {
	domain := s.session.LocalAddr().Domain()
	log.WithField("instance", domain.String()).Infof("no gateway configured, try to find one on your instance")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*500))
	defer cancel()
	iter := disco.FetchItemsIQ(ctx, "", stanza.IQ{To: domain}, s.session)
	if err := iter.Err(); err != nil {
		iter.Close()
		return err
	}
	addresses := []jid.JID{iter.Item().JID}
	for iter.Next() {
		if err := iter.Err(); err != nil {
			iter.Close()
			return err
		}
		addresses = append(addresses, iter.Item().JID)
	}
	iter.Close()
	for _, j := range addresses {
		log.Debugf("check for UnifiedPush gateway: %s", j)
		if err := s.testAndUseGateway(j); err == nil {
			return nil
		}
	}
	log.WithField("gateway", XMPPUPDemoJID.String()).Infof("no UnifiedPush gateway on your instance - try demo server")
	return s.testAndUseGateway(XMPPUPDemoJID)
}

// testAndUseGateway
func (s *XMPPService) testAndUseGateway(address jid.JID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*500))
	defer cancel()
	info, err := disco.GetInfo(ctx, "", address, s.session)
	if err != nil {
		return err
	}
	for _, f := range info.Features {
		if f.Var == messages.Space {
			s.gateway = address
			log.WithField("gateway", s.gateway.String()).Debug("tested UnifiedPush XMPP gateway should work")
			return nil
		}
	}
	return errors.New("this is no UnifiedPush gateway")
}

// Register handler of DBUS Distribution
func (s *XMPPService) Register(appID, appToken string) (string, string, error) {
	logger := log.WithFields(map[string]interface{}{
		"appID":    appID,
		"appToken": appToken,
	})
	conn := s.store.NewConnection(appID, appToken, s.gateway.String()+":"+s.session.LocalAddr().Bare().String())
	if conn == nil {
		errStr := "error to store public token"
		err := errors.New(errStr)
		logger.WithField("error", err).Error("unable to register")
		return "", errStr, err
	}
	logger = logger.WithField("publicToken", conn.PublicToken)
	iq := messages.RegisterIQ{
		IQ: stanza.IQ{
			Type: stanza.SetIQ,
			To:   s.gateway,
		},
	}
	iq.Register.Token = &messages.TokenData{Body: conn.PublicToken}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*500))
	defer cancel()
	t, err := s.session.EncodeIQ(ctx, iq)
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
		return endpoint.Body, "", nil
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
func (s *XMPPService) Unregister(appToken string) {
	logger := log.WithFields(map[string]interface{}{
		"appToken": appToken,
	})
	conn, err := s.store.DeleteConnection(appToken)
	if err != nil {
		log.WithField("error", err).Error("delete connection on storage")
		return
	}
	logger = logger.WithFields(map[string]interface{}{
		"appID":       conn.AppID,
		"publicToken": conn.PublicToken,
		"gateway":     conn.Settings,
	})
	if err = s.dbus.NewConnector(conn.AppID).Unregistered(conn.AppToken); err != nil {
		logger.WithField("error", err).Error("send unregister per dbus ")
		return
	}
	logger.Info("distributor-unregister")
}
