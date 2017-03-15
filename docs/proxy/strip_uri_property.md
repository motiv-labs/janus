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
        "upstream_url": "http://my-api.com",
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
