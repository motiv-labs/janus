##### The `strip_path` property

It may be desirable to specify a URI prefix to match an API, but not
include it in the upstream request. To do so, use the `strip_path` boolean
property by configuring an API like this:

```json
{
    "name": "My API",
    "proxy": {
        "strip_path" : true,
        "listen_path": "/service/*",
        "upstreams" : {
            "balancing": "roundrobin",
            "targets": [
                {"target": "http://my-api.com"}
            ]
        },
        "methods": ["GET"]
    }
}
```

Enabling this flag instructs Janus that when proxying this API, it should **not**
include the matching URI prefix in the upstream request's URI. For example, the
following client's request to the API configured as above:

```http
GET /service/path/to/resource HTTP/1.1
Host: my-api.com
```

Will cause Janus to send the following request to your upstream service:

```http
GET /path/to/resource HTTP/1.1
Host: my-api.com
```

The `strip_path` property can be used in tandem with `Named URL Parameters`.
i.e:
```json
{
    "name": "My API",
    "proxy": {
        "strip_path" : true,
        "listen_path": "/prepath/{service}/*",
        "upstreams" : {
            "balancing": "roundrobin",
            "targets": [
                {"target": "http://{service}.com"}
            ]
        },
        "methods": ["GET"]
    }
}
```

Akin to the previous example a request to janus like this:
```http
GET /prepath/my-service/path/to/resource HTTP/1.1
Host: janus
```

Will cause Janus to send the following request to your upstream service:

```http
GET /path/to/resource HTTP/1.1
Host: my-service.com
```

This is because when the `strip_path` property is set to **true** and a `Named URL parameter` is used, the first instance of each section of the `listen_path` (delineated by `/`) will be removed from the upstream request. This includes the parameter name and regex.
<br>
In addition to that, the first instance of each `Named URL parameter` will be removed from the upstream request.
