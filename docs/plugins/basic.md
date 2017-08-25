# HTTP Basic Auth

Add Basic Authentication to your APIs, with username and password protection. The plugin will check for valid credentials in the `Authorization` header.

## Configuration

The plain basic auth config:

```json
"basic_auth": {
    "enabled": true
}
```

Here is a simple definition of the available configurations.

| Configuration                 | Description                                                         |
|-------------------------------|---------------------------------------------------------------------|
| name                          | Name of the plugin to use, in this case: basic_auth        |
| enabled                       | Is the plugin enabled?  |

## Usage

In order to use the plugin, you first need to create some users first. By enabling this plugins in any endpoint There is a simple API that you can use to create new users.

## Create an User

You need to create an user that will be used to authenticate. To create an user you can execute the following request:

{% codetabs name="HTTPie", type="bash" -%}
http -v POST http://localhost:8081/credentials/basic_auth "Authorization:Bearer yourToken" username=lanister password=pay-your-debt
{%- language name="CURL", type="bash" -%}
curl -X POST http://localhost:8081/credentials/basic_auth -H 'authorization: Bearer yourToken' -H 'content-type: application/json' -d '{"username": "lanister", "password": "pay-your-debt"}'
{%- endcodetabs %}

| FORM PARAMETER | Description                                     |
|----------------|-------------------------------------------------|
| username       | The username to use in the Basic Authentication |
| password       | The password to use in the Basic Authentication |

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
