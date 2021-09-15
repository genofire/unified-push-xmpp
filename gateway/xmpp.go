package main

import (
	"context"
	"encoding/xml"
	"io"
	"net"

	"github.com/bdlm/log"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/component"
	"mellium.im/xmpp/disco"
	"mellium.im/xmpp/disco/info"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/ping"
	"mellium.im/xmpp/stanza"

	"dev.sum7.eu/genofire/unified-push-xmpp/messages"
)

type XMPPService struct {
	Addr    string `toml:"address"`
	JID     string `toml:"jid"`
	Secret  string `toml:"secret"`
	session *xmpp.Session
}

func (s *XMPPService) Run(jwt JWTSecret, endpoint string) error {
	var err error
	j := jid.MustParse(s.JID)
	ctx := context.TODO()
	conn, err := net.Dial("tcp", s.Addr)
	if err != nil {
		return err
	}
	if s.session, err = component.NewSession(
		ctx, j.Domain(),
		[]byte(s.Secret), conn,
	); err != nil {
		return err
	}
	defer func() {
		log.Info("closing xmpp connection")
		if err := s.session.Close(); err != nil {
			log.Errorf("Error closing session: %q", err)
		}
		log.Info("closing xmpp connection")
		if err := s.session.Conn().Close(); err != nil {
			log.Errorf("Error closing connection: %q", err)
		}
	}()
	log.Infof("connected with %s", s.session.LocalAddr())
	return s.session.Serve(mux.New(
		disco.Handle(),
		ping.Handle(),
		XMPPUpHandle(jwt, endpoint),
	))
}

// SendMessage of an UP Notification
func (s *XMPPService) SendMessage(to jid.JID, publicToken, content string) error {
	log.WithFields(map[string]interface{}{
		"to":          to.String(),
		"publicToken": publicToken,
	}).Debug("forward message to xmpp")
	return s.session.Encode(context.TODO(), messages.Message{
		Message: stanza.Message{
			To:   to,
			From: jid.MustParse(s.JID),
			// Type: stanza.ChatMessage,
			Type: stanza.NormalMessage,
		},
		PublicToken: publicToken,
		Body:        content,
	})
}

// XMPPUpHandler struct
// for handling UnifiedPush specifical requests
type XMPPUpHandler struct {
	jwtSecret JWTSecret
	endpoint  string
}

// XMPPUpHandle - setup UnfiedPush handler to mux
func XMPPUpHandle(jwt JWTSecret, endpoint string) mux.Option {
	return func(m *mux.ServeMux) {
		s := &XMPPUpHandler{jwtSecret: jwt, endpoint: endpoint}
		// register - get + set - need direct handler (not IQFunc) to bind ForIdentities and ForFeatures to disco
		mux.IQ(stanza.SetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s)(m)
		mux.IQ(stanza.GetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s)(m)
		// unregister - get + set
		mux.IQFunc(stanza.SetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister)(m)
		mux.IQFunc(stanza.GetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister)(m)
	}
}

var (
	upIdentity = info.Identity{
		Category: "pubsub",
		Type:     "push",
		Name:     "Unified Push over XMPP",
	}
	upFeature = info.Feature{Var: messages.Space}
)

// ForIdentities disco handler
func (h *XMPPUpHandler) ForIdentities(node string, f func(info.Identity) error) error {
	if node != "" {
		log.Debugf("response disco feature for %s", node)
		return nil
	}
	var err error
	err = f(upIdentity)
	if err != nil {
		return err
	}
	log.Debug("response disco identity")
	return nil
}

// ForFeatures disco handler
func (h *XMPPUpHandler) ForFeatures(node string, f func(info.Feature) error) error {
	if node != "" {
		log.Debugf("response disco feature for %s", node)
		return nil
	}
	var err error
	err = f(upFeature)
	if err != nil {
		return err
	}
	log.Debug("response disco feature")
	return nil
}

// HandleIQ - handleRegister for UnifiedPush request
func (h *XMPPUpHandler) HandleIQ(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	if start.Name.Local != "register" || start.Name.Space != messages.Space {
		return nil
	}
	reply := messages.RegisterIQ{
		IQ: stanza.IQ{
			ID:   iq.ID,
			Type: stanza.ErrorIQ,
			From: iq.To,
			To:   iq.From,
		},
	}
	defer func() {
		if err := t.Encode(reply); err != nil {
			log.Errorf("sending register response: %v", err)
		}
	}()
	log.Infof("recieved iq: %v", iq)

	tokenData := messages.TokenData{}
	err := xml.NewTokenDecoder(t).Decode(&tokenData)
	if err != nil && err != io.EOF {
		log.Warnf("decoding message: %q", err)
		reply.Register.Error = &messages.ErrorData{Body: "unable decode"}
		return nil
	}
	publicToken := tokenData.Body
	if publicToken == "" {
		log.Warnf("no token found: %v", publicToken)
		reply.Register.Error = &messages.ErrorData{Body: "no token"}
		return nil
	}
	endpointToken, err := h.jwtSecret.Generate(iq.From, publicToken)
	if err != nil {
		log.Errorf("unable entpointToken generation: %v", err)
		reply.Register.Error = &messages.ErrorData{Body: "endpointToken error on gateway"}
		return nil
	}
	endpoint := h.endpoint + "/UP?token=" + endpointToken
	reply.IQ.Type = stanza.ResultIQ
	reply.Register.Endpoint = &messages.EndpointData{Body: endpoint}
	log.Debugf("generate respone: %v", endpoint)
	return nil
}

// handleUnregister for UnifiedPush request
func (h *XMPPUpHandler) handleUnregister(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	reply := messages.UnregisterIQ{
		IQ: stanza.IQ{
			ID:   iq.ID,
			Type: stanza.ErrorIQ,
			From: iq.To,
			To:   iq.From,
		},
	}
	defer func() {
		if err := t.Encode(reply); err != nil {
			log.Errorf("sending unregister response: %v", err)
		}
	}()
	log.Debugf("unregistered unhandled: %v", start)

	reply.Unregister.Error = "not implemented"
	return nil
}
