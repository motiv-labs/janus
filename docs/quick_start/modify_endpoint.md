# Modify (Update/Delete) an endpoint

As your infrastructure grows and changes you need to modify existing endpoints configurations and delete unused one. In this section we'll add one more endpoint, like we did in the last section, will make some changes to it and then remove it.

## 1. Adding a new endpoint

Now that you are authenticated, you can send a request to `/apis` to create a proxy.

Create an `example-modify.json` file containing your endpoint configuration:

```json
{
    "name" : "my-endpoint-to-modify",
    "active" : true,
    "proxy" : {
        "listen_path" : "/temporary/*",
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
http -v POST localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < example-modify.json
{%- language name="CURL", type="bash" -%}
curl -X "POST" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json" -d @example-modify.json
{%- endcodetabs %}

This will create a proxy to `http://www.mocky.io/v2/595625d22900008702cd71e8` (which is a fake api) when you hit Janus on `GET /temporary`.

By this point this is exactly what we did in the previous section, when we added first endpoint.

You can make sure that the configuration is created using an API, same as we did before

{% codetabs name="HTTPie", type="bash" -%}
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
{%- language name="CURL", type="bash" -%}
curl -X "GET" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json"
{%- endcodetabs %}

And that it proxies requests to a desired destination, again, same as we did previously

{% codetabs name="HTTPie", type="bash" -%}
http -v GET http://localhost:8080/temporary
{%- language name="CURL", type="bash" -%}
curl -vX "GET" http://localhost:8080/temporary
{%- endcodetabs %}

## 2. Update existing endpoint

Updating existing endpoint is almost the same as creating new one except for the method and api endpoint used.

First, let's modify a file containing endpoint configuration and change `active` to `false` to deactivate existing endpoint, so the resulting file will look like this:

```json
{
    "name" : "my-endpoint-to-modify",
    "active" : false,
    "proxy" : {
        "listen_path" : "/temporary/*",
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

To modify an endpoint we'll use the same `/apis` REST API endpoint, but with the `PUT` method and endpoint name, that is set in the `name` field, as resource:

{% codetabs name="HTTPie", type="bash" -%}
http -v PUT localhost:8081/apis/my-endpoint-to-modify "Authorization:Bearer yourToken" "Content-Type: application/json" < example-modify.json
{%- language name="CURL", type="bash" -%}
curl -X "PUT" localhost:8081/apis/my-endpoint-to-modify -H "Authorization:Bearer yourToken" -H "Content-Type: application/json" -d @example-modify.json
{%- endcodetabs %}

Make sure that the configuration is modified

{% codetabs name="HTTPie", type="bash" -%}
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
{%- language name="CURL", type="bash" -%}
curl -X "GET" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json"
{%- endcodetabs %}

And that it is not proxying requests anymore

{% codetabs name="HTTPie", type="bash" -%}
http -v GET http://localhost:8080/temporary
{%- language name="CURL", type="bash" -%}
curl -vX "GET" http://localhost:8080/temporary
{%- endcodetabs %}

## 3. Delete existing endpoint

Some endpoints become outdated and it does not make sense to store all of them in non-active state, so we can simply remove them.

To modify an endpoint we'll use the same `/apis` REST API endpoint, but with the `DELETE` method and endpoint name, that is set in the `name` field, as resource:

{% codetabs name="HTTPie", type="bash" -%}
http -v DELETE localhost:8081/apis/my-endpoint-to-modify "Authorization:Bearer yourToken"
{%- language name="CURL", type="bash" -%}
curl -X "DELETE" localhost:8081/apis/my-endpoint-to-modify -H "Authorization:Bearer yourToken"
{%- endcodetabs %}

Make sure that the configuration is removed

{% codetabs name="HTTPie", type="bash" -%}
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
{%- language name="CURL", type="bash" -%}
curl -X "GET" localhost:8081/apis -H "Authorization:Bearer yourToken" -H "Content-Type: application/json"
{%- endcodetabs %}

[Next](add_plugins.md) we'll learn how to add plugins to our endpoint.
