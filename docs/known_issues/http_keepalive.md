# Stale HTTP Keep-Alive

Janus proxies requests with HTTP keep-alive enabled to reduce resource usage and latency.
Unless explicitly configured, the proxy uses default idle connection timeout of 90 seconds.
This means keep-alive connections that are not used for 90 seconds are closed, and the next request will use a fresh connection.    

This creates a stale connection if an endpoint is under constant load – therefore never reaching 90 seconds idle timeout.
If a DNS update occurs and the target host of the endpoint points to a different address, this connection will still stay alive indefinitely, proxying requests to the wrong address until either Janus is restarted, or the endpoint gets idle long enough to exceed the 90 seconds of the idle connection timeout so it's closed and a new connection is created.

To help with this problem, Janus can be started with an optional environment variable `CONN_PURGE_INTERVAL` that flushes the idle HTTP keep-alive connections periodically. 
This allows Janus to retain the benefits of HTTP keep-alive connections while limiting the maximum duration of stale connection kept alive.   
