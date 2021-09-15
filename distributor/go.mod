module dev.sum7.eu/genofire/unified-push-xmpp/distributor

go 1.17

require (
	dev.sum7.eu/genofire/golang-lib v0.0.0-20210912204316-9b2fe62df536
	dev.sum7.eu/genofire/unified-push-xmpp/messages v0.0.0-20210914093612-4a88e1d4a772
	github.com/bdlm/log v0.1.20
	github.com/google/uuid v1.3.0
	mellium.im/sasl v0.2.1
	mellium.im/xmlstream v0.15.3-0.20210221202126-7cc1407dad4c
	mellium.im/xmpp v0.19.0
	unifiedpush.org/go/np2p_dbus v0.0.0-20210914192133-e2c19b86a23f
)

require (
	github.com/bdlm/std v1.0.1 // indirect
	github.com/godbus/dbus/v5 v5.0.5 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gorm.io/driver/sqlite v1.1.5 // indirect
	gorm.io/gorm v1.21.15 // indirect
	mellium.im/reader v0.1.0 // indirect
)

replace dev.sum7.eu/genofire/unified-push-xmpp/messages => ../messages

replace mellium.im/xmpp => ../../../../mellium.im/xmpp
