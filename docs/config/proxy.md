# Proxy configuration

| Configuration         | Description                                                                            |
|-----------------------|----------------------------------------------------------------------------------------|
| preserve_hosts        | Enable the [preserve host](/docs/proxy/preserve_host_property.md) definition           |
| listen_path           | Defines the [endpoint](/docs/proxy/request_uri.md) that will be exposed in Janus       |
| upstreams             | Defines the [endpoints](/docs/proxy/upstreams.md) that the request will be forwarded to|
| strip_path            | Enable the [strip URI](/docs/proxy/strip_uri_property.md) rule on this proxy           |
| methods               | Defines which [methods](/docs/proxy/request_http_method.md) are enabled for this proxy |
| hosts                 | Defines which [hosts](/docs/proxy/request_http_header.md) are enabled for this proxy   |
| forwarding_timeouts.dial_timeout | The amount of time to wait until a connection to a backend server can be established. Defaults to 30 seconds. If zero, no timeout exists. You must use any format that is compatible with [time.Duration](https://golang.org/pkg/time/#Duration) |
| forwarding_timeouts.response_header_timeout | The amount of time to wait for a server's response headers after fully writing the request (including its body, if any). If zero, no timeout exists. You must use any format that is compatible with [time.Duration](https://golang.org/pkg/time/#Duration) |
