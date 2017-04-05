# Proxy Authentication Methods

When configuring your API you can choose between different authentication methods, these are:

* Basic/Digest Authentication
* OAuth 2.0
* JWT

We've tried to desing Janus in a way that the authentication provider are simple to setup 
and completely decoupled from the gateway. You can protect you API endpoints by simply 
following a few steps:

## Configure a new authentication provider

First we need to create an authentication provider. Use your choosen [configuration method](config.md)
to configure your auth provider.

## Attaching the oauth server to an API

To use the auth configuration we need to attach it to one of our configured APIs.
You can do that by simply adding the [oauth](/docs/plugins/oauth.md) plugin to your API.

## Restart the gateway to apply the configurations and done.

## Query your OAuth Servers

If you want to see the available auth providers that you've configured just do:

```
http -v GET localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json"
```
