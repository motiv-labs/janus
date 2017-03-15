### Request Host header

Routing a request based on its Host header is the most straightforward way to proxy traffic through Janus, as this is the intended usage of the HTTP Host header. Janus makes it easy to do so via the hosts field of the API entity.

`hosts` accepts multiple values, which must be in an array format when specifying them via the Admin API:

```json
{
    "hosts": ["my-api.com", "example.com", "service.com"]
}
```

To satisfy the hosts condition of this API, any incoming request from a client must now have its Host header set to one of:

```http
Host: my-api.com
```

or:

```http
Host: example.com
```

or:

```http
Host: service.com
```
