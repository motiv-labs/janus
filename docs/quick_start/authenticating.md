## Authenticating

To start using the Janus administration API you need to get a [JSON Web Token](https://jwt.io) and provide it in every single request using the `Authorization` header.

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
  # The algorithm that you want to use to create your JWT
  algorithm = "HS256"
  # This is the secret that you will use to encrypt your JWT
  secret = "secret key"

  [web.credentials.github]
  # The github owner/organizations that will be allowed to login on the private API
  organizations = ["hellofresh"]
  # A map of the owner/organization and the team name that will have access to the private API
  teams = {hellofresh = "devs"}
```

### Basic

Alternatively, you can authenticate against the admin API using HTTP `Basic` Authentication.

{% codetabs name="HTTPie", type="bash" -%}
http -v --json POST localhost:8081/login username=admin password=admin
{%- language name="CURL", type="bash" -%}
curl -X "POST" localhost:8081/login -d '{"username": "admin", "password": "admin"}' -H "Content-Type: application/json"
{%- endcodetabs %}

The username and password default to *admin*/*admin*, and **should be changed** using the following [configuration](../install/configuration.md):

```toml
[web.credentials]
  # The algorithm that you want to use to create your JWT
  algorithm = "HS256"
  # This is the secret that you will use to encrypt your JWT
  secret = "secret key"

  [web.credentials.basic]
  # A dictionary with the user and password
  users = [
    {admin = "admin"}
  ]
```
