
### Overview

From a high level perspective, Janus will listen for HTTP traffic on its configured proxy port (`8080` by default), recognize which upstream service is being requested, run the configured middlewares for that API, and forward the HTTP request upstream to your own API or service.

When a client makes a request to the proxy port, Janus will decide to which upstream service or API to route (or forward) the incoming request, depending on the API configuration in Janus, which is managed via the Admin API. You can configure APIs with various properties, but the three relevant ones for routing incoming traffic are hosts, uris, and methods.

If Janus cannot determine to which upstream API a given request should be routed, Janus will respond with:

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
Server: Janus/<x.x.x>

{
    "message": "no API found with those values"
}
```
