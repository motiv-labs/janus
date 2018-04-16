# Clustering / High Availability

Multiple Janus nodes pointing to the same datastore must belong to the same "Janus Cluster".

A Janus cluster allows you to scale the system horizontally by adding more machines to handle a bigger load of incoming requests, and they all share the same data since they point to the same datastore.

A Janus cluster can be created in one datacenter, or in multiple datacenters, in both cloud or on-premise environments. Janus will take care of joining and leaving a node automatically in a cluster, as long as the node is configured properly.

## Configuration update

Some backends requires that you define an update interval, which is used to check for changes 
on that storage. You can do that by setting the cluster configuration like this:

```toml
[cluster]
  UpdateFrequency = "5s"
```

