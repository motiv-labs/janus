# Organization Auth

Create users with organizations and add an organization header to upstream requests.
The plugin works similarly to basic auth with the exception that it also tracks an organization for users.
It will also add the organization of the users to the header of upstream requests.

**Limitations**
1. This plugin only works as a Basic Authentication not Oauth.
2. This plugin only works with Cassandra DB repo.

## Configuration

The plain organization header config:

```json
{
  "name": "organization_header",
  "enabled":  true
}
```

Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: organization_header        |
| enabled                       | Is the plugin enabled?  |

## Usage

You need to create an user that will be used to authenticate. To create an user you can execute the following request:

{% codetabs name="HTTPie", type="bash" -%}
http -v POST http://localhost:8081/credentials/basic_auth "Authorization:Bearer yourToken" username=lanister password=pay-your-debt organization=motiv
{%- language name="CURL", type="bash" -%}
curl -X POST http://localhost:8081/credentials/basic_auth -H 'authorization: Bearer yourToken' -H 'content-type: application/json' -d '{"username": "lanister", "password": "pay-your-debt", "organization": "motiv"}'
{%- endcodetabs %}

| FORM PARAMETER | Description                                     |
|----------------|-------------------------------------------------|
| username       | The username to use in the Basic Authentication |
| password       | The password to use in the Basic Authentication |
| organization   | The organization of the user                    |

## Using the Credential

The authorization header must be base64 encoded. For example, if the credential uses `lanister` as the username and `pay-your-debt` as the password, then the field's value is the base64-encoding of lanister:pay-your-debt, or bGFuaXN0ZXI6cGF5LXlvdXItZGVidA==.

Then the `Authorization` header must appear as:

Authorization: Basic bGFuaXN0ZXI6cGF5LXlvdXItZGVidA==
Simply make a request with the header:

{% codetabs name="HTTPie", type="bash" -%}
http -v http://localhost:8080/example "Authorization:Basic bGFuaXN0ZXI6cGF5LXlvdXItZGVidA=="
{%- language name="CURL", type="bash" -%}
curl -v http://localhost:8080/example -H 'Authorization:Basic bGFuaXN0ZXI6cGF5LXlvdXItZGVidA=='
{%- endcodetabs %}

## Using the Header

Once the organization has been paired with a user any request that proxies through Janus will contain the `X-Organization` header with a value equal to the organization paired with the user.
