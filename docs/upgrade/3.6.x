# 3.5.x to 3.6.x Upgrade Notes

If you're using `MongoDB` as configuration database run the following script against `api_specs` collection:

```javascript
db.getCollection('api_specs').find({}).forEach(function(doc) {
    if (!doc.proxy.upstreams || !doc.proxy.upstreams.targets) && !!doc.proxy.upstream_url {
        doc.plugins.push({
          "upstreams": {
            "balancing": "roundrobin",
            "targets": [{"target": doc.proxy.upstream_url}]
           }
        });
        
        delete doc.proxy.upstream_url;
        db.api_specs.update({"_id": doc._id}, doc);
    }
});
```

For the `oauth_servers` collection, run:

```javascript
db.getCollection('oauth_servers').find({}).forEach(function(doc) {
    
    fn = function(p) {
      if !!p && (!p.upstreams || !p.upstreams.targets) && !!p.upstream_url {
          p.push({
            "upstreams": {
              "balancing": "roundrobin",
              "targets": [{"target": p.upstream_url}]
             }
          });

          delete p.upstream_url;
      }
    }
    
    
    fn(doc.oauth_endpoints.authorize);
    fn(doc.oauth_endpoints.token);
    fn(doc.oauth_endpoints.info);
    fn(doc.oauth_endpoints.revoke);
    fn(doc.oauth_endpoints.introspect);
    fn(doc.oauth_client_endpoints.create);
    fn(doc.oauth_client_endpoints.remove);
    
    db.oauth_servers.update({"_id": doc._id}, doc);
});
```
