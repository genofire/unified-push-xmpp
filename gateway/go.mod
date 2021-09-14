module dev.sum7.eu/genofire/unified-push-xmpp/gateway

go 1.17

require (
	dev.sum7.eu/genofire/golang-lib v0.0.0-20210912204316-9b2fe62df536
	dev.sum7.eu/genofire/unified-push-xmpp/messages v0.0.0-20210914093612-4a88e1d4a772
	github.com/bdlm/log v0.1.20
	github.com/bdlm/std v1.0.1
	github.com/gin-gonic/gin v1.7.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/stretchr/testify v1.7.0
	mellium.im/xmlstream v0.15.3-0.20210221202126-7cc1407dad4c
	mellium.im/xmpp v0.19.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chenjiandongx/ginprom v0.0.0-20210617023641-6c809602c38a // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sessions v0.0.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-contrib/static v0.0.1 // indirect
	github.com/gin-gonic/autotls v0.0.3 // indirect
	github.com/go-mail/mail v2.3.1+incompatible // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/gorm v1.21.15 // indirect
	gorm.io/plugin/prometheus v0.0.0-20210820101226-2a49866f83ee // indirect
	mellium.im/reader v0.1.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace dev.sum7.eu/genofire/unified-push-xmpp/messages => ../messages

replace mellium.im/xmpp => ../../../../mellium.im/xmpp
