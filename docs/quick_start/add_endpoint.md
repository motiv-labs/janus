# Add an endpoint

> Note: If you are using the file-based configuration you will not be able to use the write administration API to add new endpoints. Please check this [tutorial](./file_system.md) if you'd like to add a new endpoint using the file-based configuration.

The main feature of the API Gateway is to proxy the requests to different services, so let's do this.

## Authenticating

To start using the Janus adminstration API you need to get a [JSON Web Token](https://jwt.io) and provide it in every single request using the `Authorization` header.

You can choose to log in with either `github` or `basic` providers.

### Github

To login with Github, you need to send a valid Github access token in the Authorization header. This token will be exchanged for a JWT that you can use to make requests to the admin gateway API.

You can choose to either go through the [oAuth2](https://developer.github.com/v3/guides/basics-of-authentication/) flows to authorize an user on github, or generate a [Personal Access Token](https://github.com/settings/tokens) and provide that instead.

Authentication is then performed with the following request:

{% codetabs name="HTTPie", type="bash" -%}
http -v --json POST localhost:8081/login?provider=github "Authorization:Bearer githubToken"
{%- language name="CURL", type="bash" -%}
curl -X "POST" localhost:8081/login?provider=github -H 'Authorization:Bearer githubToken'
{%- endcodetabs %}

You can also configure which organizations/teams will be allowed to log into the Admin API. This can be done with the following [configuration](../install/configuration.md):

```toml
  [web.credentials]
    secret = "secret"

    [web.credentials.github]
    organizations = ["hellofresh"]
    teams = [
      {organizationName = "hellofresh", TeamName = "Devs"}
    ]
```

### Basic

Alternatively, you can authenticate against the admin API using HTTP `Basic` Authentication.

{% codetabs name="HTTPie", type="bash" -%}
http -v --json POST localhost:8081/login username=admin password=admin
{%- language name="CURL", type="bash" -%}
curl -X "POST" localhost:8081/login -d '{"username": "admin", "password": "admin"}'
{%- endcodetabs %}

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>

The username and password default to *admin*/*admin*, and **should be changed** using the following [configuration](../install/configuration.md):

```toml
  [web.credentials]
    secret = "secret"

    [web.credentials.basic]
    users = [
      {username = "admin", password = "admin"}
    ]
```


## Adding a new endpoint

Now that you are authenticated, you can send a request to `/apis` to create a proxy.

Create an `example.json` file containing your endpoint configuration:

```json
{
    "name" : "my-endpoint",
    "active" : true,
    "proxy" : {
        "listen_path" : "/example/*",
        "upstream_url" : "http://www.mocky.io/v2/595625d22900008702cd71e8",
        "methods" : ["GET"]
    }
}
```

And now let's add it to Janus:

{% codetabs name="HTTPie", type="bash" -%}
http -v POST localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < example.json
{%- language name="CURL", type="bash" -%}
curl -X "POST" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json" -d @example.json
{%- endcodetabs %}

This will create a proxy to `http://www.mocky.io/v2/595625d22900008702cd71e8` (which is a fake api) when you hit Janus on `GET /example`.

## 2. Verify that your API has been added


You can use the REST API to query all available APIs and Auth Providers. Simply make a request 
to `/apis`:

{% codetabs name="HTTPie", type="bash" -%}
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
{%- language name="CURL", type="bash" -%}
curl -X "GET" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json"
{%- endcodetabs %}

## 3. Forward your requests through Janus

Issue the following request to verify that Janus is properly forwarding
requests to your API. Note that [by default][proxy-port] Janus handles proxy
requests on port `:8080`:

{% codetabs name="HTTPie", type="bash" -%}
http -v GET http://localhost:8080/example
{%- language name="CURL", type="bash" -%}
curl -vX "GET" http://localhost:8080/example
{%- endcodetabs %}

A successful response means Janus is now forwarding requests made to `http://localhost:8000` to the `upstream_url` we configured in step #1, and is forwarding the response back to us.

[Next](add_plugins.md) we'll learn how to add plugins to our endpoint.
