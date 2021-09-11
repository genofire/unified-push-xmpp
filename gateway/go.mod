module dev.sum7.eu/genofire/unified-push-xmpp/gateway

go 1.17

require (
	dev.sum7.eu/genofire/golang-lib v0.0.0-20210907112925-492a959d8452
	github.com/bdlm/log v0.1.20
	mellium.im/sasl v0.2.1
	mellium.im/xmlstream v0.15.3-0.20210221202126-7cc1407dad4c
	mellium.im/xmpp v0.19.0
)

require (
	dev.sum7.eu/genofire/unified-push-xmpp/messages v0.0.0-00010101000000-000000000000 // indirect
	github.com/bdlm/std v1.0.1 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200103221440-774c71fcf114 // indirect
	mellium.im/reader v0.1.0 // indirect
)

replace dev.sum7.eu/genofire/unified-push-xmpp/messages => ../messages
