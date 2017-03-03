# Proxy Authentication Methods

When configuring your API you can choose between different authentication methods, these are:

* OAuth 2.0
* Basic Authentication

Let's see how to configure 
We've tried to desing the API in a way that authentication provider are simple to setup 
and completely decoupled from the gateway. You can protect you API endpoints by simply 
following a few steps:

1. Configure a new authnetication provider

You can add a new authentication provider by sending a post request to `/oauth/servers`

```
http -v POST localhost:8080/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/oauth_server.json

HTTP/1.1 201 Created
Location: /oauth/servers/6fe999d3-63d5-4fd4-ba46-d9cd843e9133
```

2. Attaching the oauth server to an API

To attach this auth provider to one of our configured routes. You can do that by setting the
`oauth_server_id` propertry to use the returned ID of the authentication provider.

3. Restart the gateway to apply the configurations and done.
