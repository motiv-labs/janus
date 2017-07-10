# Add an endpoint

> Note: If you are using the file based configuration you will not be able to use the write administration API to add new endpoints. Please check this [tutorial](file_system.md) if you'd like to add a new endpoint to a file beased configuration.

The main feature of the API Gateway is to proxy the requests to different services, so let's do this.

## Authenticating

To start using the Janus adminstration API you need to get a [JSON Web Token](https://jwt.io) and provide it in every single request
using the `Authorization` header.

To get a token you must execute:

```sh
http -v --json POST localhost:8081/login username=admin password=admin
```

The username and password are defined by the configuration called `web.credentials.username` and `web.credentials.password`. It defaults to *admin*/*admin*.

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>

With this token you can now request the administration endpoints of Janus


## Adding a new endpoint

Now that you are authenticated, you can send a request to `/apis` to create a proxy.

Just create an `example.json` file containing this:

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

And now lets add it to Janus:

```sh
http -v POST localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < example.json
```

This will create a proxy to `http://www.mocky.io/v2/595625d22900008702cd71e8` (which is a fake api) when you hit the api gateway on `GET /example`.

## 2. Verify that your API has been added


You can use the REST API to query all available APIs and Auth Providers. Simply make a request 
to `/apis`.

```bash
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
```

## 3. Forward your requests through Janus

Issue the following cURL request to verify that Janus is properly forwarding
requests to your API. Note that [by default][proxy-port] Janus handles proxy
requests on port `:8080`:

```bash
$ http -v GET http://localhost:8080/example
```

A successful response means Janus is now forwarding requests made to `http://localhost:8000` to the `upstream_url` we configured in step #1, and is forwarding the response back to us.

[Next](add_plugins.md) we'll learn how to add plugins to our endpoint.
