##### The `append_path` property

You might also want to always append the `listen_path` to the upstream `upstream_url`. 
To do so, use the `append_path` boolean property by configuring an API like this:

```json
{
    "name": "My API",
    "proxy": {
        "append_path" : true,
        "listen_path": "/service/*",
        "upstream_url": "http://my-api.com/example",
    }
}
```

Enabling this flag instructs Janus that when proxying this API, it should **always**
include the matching URI prefix in the upstream request's URI. For example, the
following client's request to the API configured as above:

```http
GET /service/path/to/resource HTTP/1.1
Host: my-api.com
```

Will cause Janus to send the following request to your upstream service:

```http
GET /example/service/path/to/resource HTTP/1.1
Host: my-api.com
```
