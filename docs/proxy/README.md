# Proxy Reference

Janus listens for traffic on four ports, which by default are:

`:8080` on which Janus listens for incoming HTTP traffic from your clients, and forwards it to your upstream services.

`:8443` on which Janus listens for incoming HTTPS traffic. This port has a similar behavior as the `:8080` port, except that it expects HTTPS traffic only. This port can be disabled via the configuration file.

`:8081` on which the [Admin API](admin_api.md) used to configure Janus listens.

`:8444` on which the [Admin API](admin_api.md) listens for HTTPS traffic.
