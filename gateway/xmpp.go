package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"

	"github.com/bdlm/log"
	"mellium.im/sasl"
	"mellium.im/xmlstream"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
)

type XMPPService struct {
	Login    string
	Password string
}

// XMLElement is for Unmarshal undefined structs a fallback - any hasn't matched element
type XMLElement struct {
	XMLName  xml.Name
	InnerXML string `xml:",innerxml"`
}

func (xs *XMPPService) Run() error {
	j := jid.MustParse(xs.Login)
	s, err := xmpp.DialClientSession(
		context.TODO(), j,
		xmpp.BindResource(),
		xmpp.StartTLS(&tls.Config{
			ServerName: j.Domain().String(),
		}),
		// sasl.ScramSha1Plus <- problem with (my) ejabberd
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
	// TODO
	err = s.Encode(context.TODO(), struct{
		stanza.Message
		Token    string       `xml:"subject,omitempty"`
		Body    string       `xml:"body,omitempty"`
	}{
		Message: stanza.Message{
			To: jid.MustParse("up-test@chat.sum7.eu"),
			// Type: stanza.ChatMessage,
			Type: stanza.NormalMessage,
		},
		Token: "691499b4-adaf-4a92-b417-40e9a68f04a6",
		Body: "New Message ;) - Titel of UP-Developing",
	})
	s.Serve(mux.New())
	return nil
}
func (xs *XMPPService) HandleXMPP(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	log.Infof("unhandled: %v", start)
	return nil
}
