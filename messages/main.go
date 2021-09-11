package messages

import (
	"encoding/xml"

	"mellium.im/xmpp/stanza"
)

// Namespace
const (
	Space = "unifiedpush.org"

	LocalRegister   = "register"
	LocalUnregister = "unregister"
)

// RegisterIQ with stanza
type RegisterIQ struct {
	stanza.IQ
	Register struct {
		XMLName  xml.Name      `xml:"unifiedpush.org register"`
		Token    *TokenData    `xml:"token"`
		Endpoint *EndpointData `xml:"endpoint"`
		Error    *ErrorData    `xml:"error"`
	} `xml:"register"`
}

type TokenData struct {
	XMLName xml.Name `xml:"token"`
	Body    string   `xml:",chardata"`
}

type EndpointData struct {
	XMLName xml.Name `xml:"endpoint"`
	Body    string   `xml:",chardata"`
}

type ErrorData struct {
	XMLName xml.Name `xml:"error"`
	Body    string   `xml:",chardata"`
}

// UnregisterIQ with stanza
type UnregisterIQ struct {
	stanza.IQ
	Unregister
}

// Unregister without stanza
type Unregister struct {
	XMLName xml.Name `xml:"unifiedpush.org unregister"`
	// set
	Token string `xml:"token,omitempty"`
	// result
	Success *string `xml:"success,omitempty"`
	// error
	Error string `xml:"error,omitempty"`
}

// Message of push notification - with stanza
type Message struct {
	stanza.Message
	Token string `xml:"subject"`
	Body  string `xml:"body"`
}

// MessageBody of push notification - without stanza
type MessageBody struct {
	Token string `xml:"subject"`
	Body  string `xml:"body"`
}
