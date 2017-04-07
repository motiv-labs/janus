# Adding your API - MongoDB

By choosing Mongo DB everything that we want configure on the gateway we can do it through a REST API, since all endpoints are protected, we need to [login first](auth.md).

## 1. Add your API

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

## 2. Verify that your API has been added

You can use the REST API to query all available APIs and Auth Providers. Simply make a request 
to `/apis`.

```bash
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
```

## 3. Forward your requests through Janus

Issue the following cURL request to verify that Janus is properly forwarding
requests to your API. Note that [by default][proxy-port] Janus handles proxy
requests on port `:8080`:

```bash
$ http -v GET http://localhost:8080/posts/1
```

<p align="center">
  <a href="http://g.recordit.co/vufeMjwEfg.gif">
    <img src="http://g.recordit.co/vufeMjwEfg.gif">
  </a>
</p>


A successful response means Janus is now forwarding requests made to
`http://localhost:8000` to the `upstream_url` we configured in step #1,
and is forwarding the response back to us.

Do you want to protect your API? Check it out [here](proxy_auth_methods.md) how to do it.
