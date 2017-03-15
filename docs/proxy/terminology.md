### Terminology

`API`: This term refers to the API entity of Janus. You configure your APIs, that point to your own upstream services, through the Admin API.

`Middleware`: This refers to Janus "middleware", which are pieces of business logic that run in the proxying lifecycle. Middleware can be configured through the Admin API - either globally (all incoming traffic) or on a per-API basis.

`Client`: Refers to the downstream client making requests to Janus's proxy port.

`Upstream service`: Refers to your own API/service sitting behind Janus, to which client requests are forwarded.
