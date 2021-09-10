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
	Register
}

// Register without stanza
type Register struct {
	XMLName xml.Name `xml:"unifiedpush.org register"`
	// set
	Token string `xml:"token,omitempty"`
	// result
	Endpoint string `xml:"endpoint,omitempty"`
	// error
	Error string `xml:"error,omitempty"`
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
