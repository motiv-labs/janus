# Add Authentication

When configuring your API you can choose between different authentication methods, these are:

* Basic/Digest Authentication
* OAuth 2.0
* JWT

We tried to design Janus in a way that the authentication provider are simple to setup 
and completely decoupled from the gateway. 

Let's add an OAuth2 authentication to our endpoint.

## 1. Configure OAuth2

First of all let's configure our oAuth2 provider. This could be any OAuth2 provider: Google, Facebook, etc...
Let's bring a container up with a mocked OAuth2 server:

```sh
cd examples/front-proxy-auth
docker-compose up -d auth-service
```

Lets create a file with the oAuth2 configuration called `auth.json`:

```json
{
    "name" : "local",
    "oauth_endpoints" : {
        "token" : {
            "listen_path" : "/auth/token",
            "upstream_url" : "http://auth-service:8080/token",
            "strip_path" : true,
            "append_path" : false,
            "methods" : ["POST"]
        }
    },
    "token_strategy" : {
        "name" : "jwt",
        "settings" : {
            "secret" : "secret"
        }
    }
}
```

So, what we've done here? 

1. The first thing is to give a `name` for the oAuth2 server. 
2. Within `oauth_endpoints` we setup only one endpoint for this example, which is the `token`. If you look closely you will see that the `oauth_endpoints.token` is just a proxy configuration, exactly the same that we used to configure our first endpoint.
3. We've defined a `token_strategy`. Here you can choose between `jwt` or `storage`, storage means that Janus will be in charge of storing and managing (expiring, refreshing, etc) the tokens once they are returned from your oauth provider. JWT means that Janus will only check for expiration and secret of the tokens, but it wont store them.
This allows Janus to not go on the auth service on every single request to check the validity of the token.

If you want to check all available configurations for OAuth2 please visit [this](/docs/auth/oauth.md).

Now lets add this configuration to Janus:

```sh
http -v POST localhost:8081/oauth/servers "Authorization:Bearer yourToken" "Content-Type: application/json" < auth.json
```

## 2. Add a plugin for your endpoint

Now that we have an oauth2 available to use, lets add it to our endpoint, just create a file called `auth_plugin.json`:

```json
{
  "plugins": [
    {
      "name": "oauth",
      "enabled": true,
      "config": {
        "server_name": "local"
      }
    }
  ]
}
```

```sh
http -v PUT localhost:8081/apis/my-endpoint "Authorization:Bearer yourToken" "Content-Type: application/json" < auth_plugin.json
```

## Testing the endpoint

If we make a request to our endpoint, it should fail:

```bash
http -v GET http://localhost:8080/example

HTTP/1.1 400 Bad Request
Content-Length: 30
Content-Type: application/json
Date: Mon, 03 Jul 2017 10:33:17 GMT

"authorization field missing"
```

Adding an Authorization field with a wrong token, gives us this:

```bash
http -v GET http://localhost:8080/example "Authorization:Bearer wrongToken"

HTTP/1.1 401 Unauthorized
Content-Length: 30
Content-Type: application/json
Date: Mon, 03 Jul 2017 10:33:31 GMT

"access token not authorized"
```

So, lets get a valid token:
```bash
http -v POST http://localhost:8080/auth/token

HTTP/1.1 200 OK
Content-Length: 188
Content-Type: application/json
Date: Mon, 03 Jul 2017 10:47:09 GMT
Server: Jetty(9.2.z-SNAPSHOT)
Vary: Origin

{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphbnVzIn0.PvBI5yIdPVtR8RVJWWZEEEVv9Bk83Q_rS7vYcKNX1wM",
    "expires_in": 21600,
    "token_type": "Bearer"
}
```

Now if we request with the right token you should be able to go through.

```bash
http -v GET http://localhost:8080/example "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkphbnVzIn0.PvBI5yIdPVtR8RVJWWZEEEVv9Bk83Q_rS7vYcKNX1wM"

HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 46
Content-Type: application/json
Date: Mon, 03 Jul 2017 12:44:07 GMT
Server: Jetty(9.2.z-SNAPSHOT)
Vary: Accept-Encoding, User-Agent

{
    "message": "Hello World!"
}
```

Of course in a real world scenario your auth service would have to check for a client ID and Secret, set an expiration on the token, etc...

