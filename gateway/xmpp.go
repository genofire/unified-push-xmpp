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
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"

	"dev.sum7.eu/genofire/unified-push-xmpp/messages"
)

type XMPPService struct {
	Login    string
	Password string
	session *xmpp.Session
}

func (s *XMPPService) Run() error {
	var err error
	j := jid.MustParse(s.Login)
	if s.session, err = xmpp.DialClientSession(
		context.TODO(), j,
		xmpp.BindCustom(func(i jid.JID,r string) (jid.JID, error) {
			// Never run
			log.Infof("try to bind: %v with ressource %s", i, r)
			return j, nil
		}),
		xmpp.StartTLS(&tls.Config{
			ServerName: j.Domain().String(),
		}),
		// sasl.ScramSha1Plus <- problem with (my) ejabberd
		//xmpp.SASL("", xs.Password, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		xmpp.SASL("", s.Password, sasl.ScramSha1, sasl.Plain),
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
	// Send initial presence to let the server know we want to receive messages.
	err = s.session.Send(context.TODO(), stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return err
	}
	log.Infof("connected with %s", s.session.LocalAddr())
	s.session.Serve(mux.New(
		// register - get + set
		mux.IQFunc(stanza.SetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s.handleRegister),
		mux.IQFunc(stanza.GetIQ, xml.Name{Local: messages.LocalRegister, Space: messages.Space}, s.handleRegister),
		// unregister - get + set
		mux.IQFunc(stanza.SetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister),
		mux.IQFunc(stanza.GetIQ, xml.Name{Local: messages.LocalUnregister, Space: messages.Space}, s.handleUnregister),
		// auto accept
		mux.PresenceFunc(stanza.SubscribePresence, xml.Name{}, s.autoSubscribe),
	))
	return nil
}
// autoSubscribe to allow sending IQ
func (s *XMPPService) autoSubscribe(presHead stanza.Presence, t xmlstream.TokenReadEncoder) error {
	log.WithField("p", presHead).Info("autoSubscribe")
	// request eighter
	t.Encode(stanza.Presence{
		Type: stanza.SubscribePresence,
		To: presHead.From,
	})
	// agree
	t.Encode(stanza.Presence{
		Type: stanza.SubscribedPresence,
		To: presHead.From,
	})
	return nil
}

func (s *XMPPService) handleRegister(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	reply := messages.RegisterIQ{
		IQ: stanza.IQ{
			ID:   iq.ID,
			Type: stanza.ErrorIQ,
			To:   iq.From,
		},
	}
	defer func(){
		if err := t.Encode(reply); err != nil {
			log.Errorf("sending response: %v", err)
		}
	}()
	log.Infof("recieved iq: %v", iq)

	tokenData := messages.TokenData{}
	err := xml.NewTokenDecoder(t).Decode(&tokenData)
	if err != nil && err != io.EOF {
		log.Errorf("Error decoding message: %q", err)
		reply.Register.Error = &messages.ErrorData{ Body: "unable decode"}
		return nil
	}
	token := tokenData.Body
	if token == "" {
		log.Errorf("no token found: %v", token)
		reply.Register.Error = &messages.ErrorData{ Body: "no token"}
		return nil
	}
	endpoint :=  "https://localhost/UP?token=" + token + "&to=" +iq.From.String()
	reply.IQ.Type = stanza.ResultIQ
	reply.Register.Endpoint = &messages.EndpointData{ Body: endpoint}
	log.Infof("generate respone: %v", endpoint)
	return nil
}
func (s *XMPPService) handleUnregister(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	reply := messages.UnregisterIQ{
		IQ: stanza.IQ{
			ID:   iq.ID,
			Type: stanza.ErrorIQ,
			To:   iq.From,
		},
	}
	defer func(){
		if err := t.Encode(reply); err != nil {
			log.Errorf("sending response: %v", err)
		}
	}()
	log.Infof("unhandled: %v", start)

	reply.Unregister.Error = "not implemented"
	return nil
}

// SendMessage of an UP Notification
func (s *XMPPService) SendMessage(to, token, content string) error {
	return s.session.Encode(context.TODO(), messages.Message{
		Message: stanza.Message{
			To: jid.MustParse(to),
			// Type: stanza.ChatMessage,
			Type: stanza.NormalMessage,
		},
		Token: token,
		Body:  content,
	})
}
