# Add an endpoint

The main feature of the API Gateway is to proxy the requests to different services, so let's do this.

## Adding a new endpoint

Now that you are authenticated, you can send a request to `/apis` to create a proxy.

Create an `example.json` file containing your endpoint configuration:

```json
{
    "name" : "my-endpoint",
    "active" : true,
    "proxy" : {
        "listen_path" : "/example/*",
        "upstreams" : {
            "balancing": "roundrobin",
            "targets": [
                {"target": "http://www.mocky.io/v2/595625d22900008702cd71e8"}
            ]
        },
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

A successful response means Janus is now forwarding requests made to `http://localhost:8000` to the elected upstream target (chosen by the load balancer) configured in step #1, and is forwarding the response back to us.

[Next](add_plugins.md) we'll learn how to add plugins to our endpoint.
