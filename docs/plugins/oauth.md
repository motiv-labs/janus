# OAuth 2.0 Authentication

Add an OAuth 2.0 authentication layer with the Authorization Code Grant, Client Credentials, Implicit Grant or Resource Owner Password Credentials Grant flow.

This plugin allows you to set the enpoints of an authentication provider. This means that Janus is not attached in any way
to the oauth flow and it simply delegate that to the oauth server.

> Note: As per the OAuth2 specs, this plugin requires the underlying API to be served over HTTPS. To avoid any confusion, we recommend that you configure your underlying API to be only served through HTTPS. 

## Configuration

Configuring the plugin is straightforward, you can add it on top of an API by executing the following request on your Janus server:

```
http -v POST localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/apis/auth0.json
```

Here is a simple definition of the available configurations

| Configuration                 | Description                                                                               |
|-------------------------------|-------------------------------------------------------------------------------------------|
| name                          | The unique name of your OAuth Server                                                      |
| oauth_endpoints.authorize     | Defines the [proxy configuration](/docs/config/proxy.md) for the `authorize` endpoint     |
| oauth_endpoints.token         | Defines the [proxy configuration](/docs/config/proxy.md) for the `token` endpoint         |
| oauth_endpoints.info          | Defines the [proxy configuration](/docs/config/proxy.md) for the `info` endpoint          |
| oauth_endpoints.revoke        | Defines the [proxy configuration](/docs/config/proxy.md) for the `revoke` endpoint        |
| oauth_client_endpoints.create | Defines the [proxy configuration](/docs/config/proxy.md) for the `create` client endpoint |
| oauth_client_endpoints.remove | Defines the [proxy configuration](/docs/config/proxy.md) for the `remove` client endpoint |
| allowed_access_types          | The allowed access types for this oauth server                                            |
| allowed_authorize_types       | The allowed authorize types for this oauth server                                         |
| auth_login_redirect           | The auth login redirect URL                                                               |
