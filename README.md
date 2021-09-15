# UnifiedPush over XMPP

[UnifiedPush](https://unifiedpush.org/) is an specification how push notifications delieveries between application server and application.
This is an implementation of the UnifiedPush specifications to delievere push notification over XMPP. 

In this project has following components:
- **Gateway** (also called an Push Provider or [Server](https://unifiedpush.org/spec/server/)) which could be registered as an XMPP Component on an Server
- **Distributor** for Linux using the [D-Bus Specification](https://unifiedpush.org/spec/dbus/) which implement an very small XMPP-Client to receive the push notifications

## XMPP Messages

### Register

Request for Register
```xml
<iq from="push-distributer@example.org/device" to="up.chat.sum7.eu" type="set" id="register-id">
  <register xmlns='unifiedpush.org'>
    <token>pubic-token</token>
  </register>
</iq>
```

**Responses**

on success:
```xml
<iq from="push-distributer@example.org/device" to="up.chat.sum7.eu" type="result" id="register-id">
  <register xmlns='unifiedpush.org'>
    <endpoint>https://an-endpoint-for-application-server.localhost/UP?token=endpoint-token</endpoint>
  </register>
</iq>
```

on failure:
```xml
<iq from="push-distributer@example.org/device" to="up.chat.sum7.eu" type="error" id="register-id">
  <register xmlns='unifiedpush.org'>
    <error>a reason of failure</error>
  </register>
</iq>
```

### Unregister

TODO

### Notification
For the push notification it-self the origin `<message/>` is used with following Position of Token and Content.

```xml
<message from="up.chat.sum7.eu" to="push-distributer@example.org/device" id="message-id">
  <subject>public-token</subject>
  <body>Here is the Notification content</body>
</message>
```

The message sender `from` should be validated from distributor, for not receiving invalid or manipulated push Messages.

## Wordings

We are using over the complete system three kind of **tokens**:
- **Endpoint Token** which is part of the *Endpoint* and is for using between Gateway and Application-Server
- **Public Token** which is used between Gateway and Distributor
- **App Token** which is used between Distributor and Application
