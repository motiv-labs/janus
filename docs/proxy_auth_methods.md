# Proxy Authentication Methods

When configuring your API you can choose between different authentication methods, these are:

* Basic Authentication
* OAuth 2.0
* JWT

We've tried to desing Janus in a way that authentication provider are simple to setup 
and completely decoupled from the gateway. You can protect you API endpoints by simply 
following a few steps:

1. Configure a new authentication provider

First we need to create an authentication provider. Use your choosen [configuration method](docs/config.md)
to configure your auth provider.

2. Attaching the oauth server to an API

To use the auth configuration we need attach it to one of our configured APIs.
You can do that by setting the `oauth_server_slug` propertry to use the configured
`slug` of the authentication provider.

3. Enable the protection on your API

After attaching it to the API now you need to make sure that you enable it. To do this
you just need to set the property `use_oauth` to `true`.

You can choose between `use_oauth`, `use_basic` or `use_jwt`.

4. Restart the gateway to apply the configurations and done.

5. Query your Auth Servers

If you want to see the available auth providers that you've configured just do:

```
http -v GET localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json"
```
