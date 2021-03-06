# UnifiedPush over XMPP - gateway

The gateway is an XMPP Component and Webserver.

The Webserver receive notifications from application server by using an JWT-Token as Endpoint Token.
The Signature JWT-Token is validated, so that only the Token generated by this gateway could be used (and no spaming over this Gateway is possible).
The JWT-Token contains the XMPP-Address of the Distributor, so that this Gateway does not require to store anything.
Current the JWT-Token is not encrypted like in [RFC 7516(https://tools.ietf.org/html/rfc7516) as an JWE [we are working on it](https://dev.sum7.eu/genofire/unified-push-xmpp/issues/5).
So the XMPP-Address could be readed by Application-Server, Application and Distributor (and could be leaked there - we hope you could trust them till we implement JWE)


The XMPP Component implements [XEP-0225](https://xmpp.org/extensions/xep-0225.html) it could be plugged in at every common server (like [ejabberd](https://docs.ejabberd.im/admin/configuration/listen/#ejabberd-service) or [Prosody](https://prosody.im/doc/components)) with an Secret and domain name.

## Install

How to configure this gateway take a look into the [config_example.toml](config_example.toml), we prefer it place it under `/etc/up-gateway.conf`

The Builded `gateway` binary, we place it under `/usr/local/bin/up-gateway`.

To get it running, we prefer following `systemd.service` file under `/etc/systemd/system/up-gateway.service`:
```toml
[Unit]
Description=UnifiedPush gateway
After=network.target
After=ejabberd.service

[Service]
Type=simple
ExecStart=/usr/local/bin/up-gateway -c /etc/up-gateway.conf
Restart=always
RestartSec=5sec

[Install]
WantedBy=multi-user.target
```

And run `systemctl enable --now up-gateway.service` to startup and start at boot.

## Demo Gateway
An Demo gateway is running under Endpoint-Address [https://up.chat.sum7.eu/UP](https://up.chat.sum7.eu/UP) and under the XMPP-Address (jid of component) [xmpp:up.chat.sum7.eu](xmpp:up.chat.sum7.eu).
Have fun by using it ;)
