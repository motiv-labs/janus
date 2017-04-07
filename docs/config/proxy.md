# Proxy configuration

| Configuration         | Description                                                                            |
|-----------------------|----------------------------------------------------------------------------------------|
| preserve_hosts        | Enable the [preserve host](/docs/proxy/preserve_host_property.md) definition           |
| listen_path           | Defines the [endpoint](/docs/proxy/request_uri.md) that will be exposed in Janus       |
| upstream_url          | Defines the OAuth provider's [endpoint](/docs/proxy/upstream_url.md) for the request   |
| strip_path            | Enable the [strip URI](/docs/proxy/strip_uri_property.md) rule on this proxy           |
| enable_load_balancing | Enable load balancing for this proxy                                                   |
| methods               | Defines which [methods](/docs/proxy/request_http_method.md) are enabled for this proxy |
| hosts                 | Defines which [hosts](/docs/proxy/request_http_header.md) are enabled for this proxy   |
