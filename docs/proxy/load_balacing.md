### Load Balancing

Janus provides multiple ways of load balancing requests to multiple backend services: a `roundrobin` (or just `rr`) method,
 and a `weight` method.

#### Round Robin

```json
{
    "name": "My API",
    "proxy": {
        "listen_path": "/foo/*",
        "upstreams" : {
            "balancing": "rr",
            "targets": [
                {"target": "http://my-api1.com"},
                {"target": "http://my-api2.com"},
                {"target": "http://my-api3.com"}
            ]
        },
        "methods": ["GET"]
    }
}
```

This configuration will apply the `roundrobin` algorithm and balance the requests to your upstreams.

#### Weight

```json
{
    "name": "My API",
    "proxy": {
        "listen_path": "/foo/*",
        "upstreams" : {
            "balancing": "weight",
            "targets": [
                {"target": "http://my-api1.com", "weight": 30},
                {"target": "http://my-api2.com", "weight": 10},
                {"target": "http://my-api3.com", "weight": 60}
            ]
        },
        "methods": ["GET"]
    }
}
```

This configuration will apply the `weight` algorithm and balance the requests to your upstreams.
