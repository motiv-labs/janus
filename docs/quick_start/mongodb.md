# Adding your API - File System

By choosing Mongo DB everything that we want configure on the gateway we can do it through a REST API, since all endpoints are protected, we need to login first.

```sh
http -v --json POST localhost:8081/login username=admin password=admin
```

The username and password are defined in an environmental variable called `ADMIN_USERNAME` and `ADMIN_PASSWORD`. It defaults to *admin*/*admin*.

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>

**Important Note**: We have two main servers running: 

* By default on port **8081** we have teh Admin REST API running
* By default on port **8080** we have our configured APIs

### Creating a proxy

The main feature of the API Gateway is to proxy the requests to a different service, so let's do this.
Now that you are authenticated, you can send a request to `/apis` to create a proxy.

```
http -v POST localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/apis/posts.json
```

<p align="center">
  <a href="http://g.recordit.co/Hi7SX8s5IA.gif">
    <img src="http://g.recordit.co/Hi7SX8s5IA.gif">
  </a>
</p>

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

Do you want to protect your API? Check it out [here](proxy_auth_methods.md) how to do it.