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
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"

	"dev.sum7.eu/genofire/unified-push-xmpp/messages"
)

type XMPPService struct {
	Addr    string `toml:"address"`
	JID     string `toml:"jid"`
	Secret  string `toml:"secret"`
	// hidden here for beautiful config file
	EndpointURL  string `toml:"-"`
	session *xmpp.Session
}

func (s *XMPPService) Run() error {
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
		log.Info("Closing session…")
		if err := s.session.Close(); err != nil {
			log.Errorf("Error closing session: %q", err)
		}
		log.Println("Closing conn…")
		if err := s.session.Conn().Close(); err != nil {
			log.Errorf("Error closing connection: %q", err)
		}
	}()
	/* Send initial presence to let the server know we want to receive messages.
	err = s.session.Send(context.TODO(), stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return err
	}*/
	log.Infof("connected with %s", s.session.LocalAddr())
	return s.session.Serve(mux.New(
		// register - get + set
		mux.IQFunc(stanza.SetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s.handleRegister),
		mux.IQFunc(stanza.GetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s.handleRegister),
		// unregister - get + set
		mux.IQFunc(stanza.SetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister),
		mux.IQFunc(stanza.GetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister),
		// mux.IQFunc("", xml.Name{}, s.handleDisco),
	))
	/* -
	return s.session.Serve(xmpp.HandlerFunc(func(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
		log.Info(start)
		return nil
	}))
	*/
}

func (s *XMPPService) handleRegister(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
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
		log.Errorf("decoding message: %q", err)
		reply.Register.Error = &messages.ErrorData{Body: "unable decode"}
		return nil
	}
	token := tokenData.Body
	if token == "" {
		log.Errorf("no token found: %v", token)
		reply.Register.Error = &messages.ErrorData{Body: "no token"}
		return nil
	}
	endpoint := s.EndpointURL+"/UP?token=" + token + "&to=" + iq.From.String()
	reply.IQ.Type = stanza.ResultIQ
	reply.Register.Endpoint = &messages.EndpointData{Body: endpoint}
	log.Infof("generate respone: %v", endpoint)
	return nil
}
func (s *XMPPService) handleUnregister(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
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
	log.Infof("unregistered unhandled: %v", start)

	reply.Unregister.Error = "not implemented"
	return nil
}

func (s *XMPPService) handleDisco(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	reply := stanza.IQ{
		ID:   iq.ID,
		Type: stanza.ErrorIQ,
		From: iq.To,
		To:   iq.From,
	}
	defer func() {
		if err := t.Encode(reply); err != nil {
			log.Errorf("sending response: %v", err)
		}
	}()
	log.Infof("recieved iq: %v", iq)
	return nil
}

// SendMessage of an UP Notification
func (s *XMPPService) SendMessage(to, token, content string) error {
	log.WithFields(map[string]interface{}{
		"to": to,
		"token": token,
	}).Info("forward message to xmpp")
	return s.session.Encode(context.TODO(), messages.Message{
		Message: stanza.Message{
			To:   jid.MustParse(to),
			From: jid.MustParse(s.JID),
			// Type: stanza.ChatMessage,
			Type: stanza.NormalMessage,
		},
		Token: token,
		Body:  content,
	})
}
