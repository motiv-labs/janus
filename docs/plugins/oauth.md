# OAuth 2.0 Authentication

Add an OAuth 2.0 authentication layer with the Authorization Code Grant, Client Credentials, Implicit Grant or Resource Owner Password Credentials Grant flow.

This plugin allows you to set the enpoints of an authentication provider. This means that Janus is not attached in any way
to the oauth flow and it simply delegate that to the oauth server.

> Note: As per the OAuth2 specs, this plugin requires the underlying API to be served over HTTPS. To avoid any confusion, we recommend that you configure your underlying API to be only served through HTTPS. 

## Configuration

> To enable this plugin to your API you should configure an [OAuth Server](auth/oauth.md) first!

Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| server_name                   | Defines the `oauth server name` to be used as your oauth provider |
