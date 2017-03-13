# Getting Started - File System

By choosing a File System based configuration we have a static way of configure the gateway (similar to nginx).

### Preparing the folder structure

By default all apis configurations are splited in separeted files (one for each API) and they live in
`/etc/janus`. Of course you can configure that route by simply defining the `DATABASE_DSN`, for instance,
you can define the value to `file:///usr/local/janus`.

Let's use the default directory and create it `/etc/janus`. We need two main folders `apis` and `auth` each one of them
holds a set of configurations for our proxies. Just run:

```
mkdir -p /etc/janus/apis
mkdir -p /etc/janus/auth
```

### Creating a proxy

The main feature of the API Gateway is to proxy the requests to a different service, so let's do this.

Just place [this example](examples/apis/posts.json) in your `apis` directory.
This will create a proxy to `https://jsonplaceholder.typicode.com/posts` when you hit the api gateway on `GET /posts`.
Now just make a request to `GET /posts`

```
http -v --json GET http://localhost:8080/posts/1
```
<p align="center">
  <a href="http://g.recordit.co/vufeMjwEfg.gif">
    <img src="http://g.recordit.co/vufeMjwEfg.gif">
  </a>
</p>

Done! You just made your first request through the gateway.

### Checking available proxies

You can use the REST API to query (Read Only) all available APIs and Auth Providers. Simply make a request 
to `/apis`.

```
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
```

Do you want to protect your API? Check it out [here](proxy_auth_methods.md) how to do it.
