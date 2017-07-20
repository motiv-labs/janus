# OAuth 2.0


## 1. Add your OAuth Server

The main feature of the API Gateway is to proxy the requests to a different service, so let's do this.
Now that you are authenticated, you can send a request to `/oauth/servers` to create a proxy.

```
http -v POST localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/front-proxy-auth/auth/auth.json
```

## 2. Verify that your API has been added

You can use the REST API to query all available APIs and Auth Providers. Simply make a request to `/oauth/servers`.

```bash
http -v GET localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json"
```

## 3. Forward your requests through Janus

Issue the following cURL request to verify that Janus is properly forwarding
requests to your OAuth Server.

This request is an example of a simple `client_credentials` flow of [OAuth 2.0](), you can try any flow that you like.

```bash
$ http -v GET http://localhost:8080/auth/token?grant_type=client_credentials "Authorization: Basic YourBasicToken"
```

# Reference

| Configuration                 | Description                                                                               |
|-------------------------------|-------------------------------------------------------------------------------------------|
| name                          | The unique name of your OAuth Server                                                      |
| oauth_endpoints.authorize     | Defines the [proxy configuration](/docs/config/proxy.md) for the `authorize` endpoint     |
| oauth_endpoints.token         | Defines the [proxy configuration](/docs/config/proxy.md) for the `token` endpoint         |
| oauth_endpoints.introspection | Defines the [proxy configuration](/docs/config/proxy.md) for the `introspection` endpoint |
| oauth_endpoints.revoke        | Defines the [proxy configuration](/docs/config/proxy.md) for the `revoke` endpoint        |
| oauth_client_endpoints.create | Defines the [proxy configuration](/docs/config/proxy.md) for the `create` client endpoint |
| oauth_client_endpoints.remove | Defines the [proxy configuration](/docs/config/proxy.md) for the `remove` client endpoint |
| allowed_access_types          | The allowed access types for this oauth server                                            |
| allowed_authorize_types       | The allowed authorize types for this oauth server                                         |
| auth_login_redirect           | The auth login redirect URL                                                               |
| secrets                       | A map of client_id: client_secret that allows you to authenticate only with the client_id |
| token_strategy.name           | The token strategy for this server. Could be `introspection` or `jwt`                           |
| token_strategy.settings.secret| If you use JWT you should set your secret or private certificate string here              |

## Token Strategy Settings

### `jwt`

JWT token validation strategy performs token validation against signature and expiration date. Currently the following
signature methods are supported:
 
* `HS256` - HMAC with SHA256 hash (symmetric key)
* `HS384` - HMAC with SHA384 hash (symmetric key)
* `HS512` - HMAC with SHA512 hash (symmetric key)
* `RS256` - RSA with SHA256 hash (asymmetric key)
* `RS384` - RSA with SHA384 hash (asymmetric key)
* `RS512` - RSA with SHA512 hash (asymmetric key)

Settings structure has the following format:

```json
[
    {"alg": "<alg1>", "key": "<key1>"},
    {"alg": "<alg2>", "key": "<key2>"},
    ...
]
```

List of signing methods allows signing method and keys rotation w/out immediate invalidation of the old one, so the
tokens signed with ald and new methods will be valid.

For backward compatibility the following settings format is also valid: `{"secret": "<key>"}` that is equal to the
new format `[{"alg": "HS256", "key", "<key>"}]`.
