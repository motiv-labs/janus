#### Using wildcard hostnames

To provide flexibility, Janus allows you to specify hostnames with wildcards in the hosts field. Wildcard hostnames allow any matching Host header to satisfy the condition, and thus match a given API.

Wildcard hostnames must contain only one asterisk at the leftmost or rightmost label of the domain. Examples:

`*.example.org` would allow Host values such as `a.example.com` and `x.y.example.com` to match.
`example.*` would allow Host values such as `example.com` and `example.org` to match.
A complete example would look like this:

```json
{
    "slug": "my-api",
    "hosts": ["*.example.com", "service.com"]
}
```

Which would allow the following requests to match this API:

```http
GET / HTTP/1.1
Host: an.example.com
```

```http
GET / HTTP/1.1
Host: service.com
```