# Adding your API - File System

By choosing a File System based configuration we have a static way of configure the gateway (similar to nginx).

## 1. Preparing the folder structure

By default all apis configurations are splited in separeted files (one for each API) and they live in
`/etc/janus`. Of course you can configure that route by simply defining the configuration `database.dsn`, for instance,
you can define the value to `file:///usr/local/janus`.

Let's use the default directory and create it `/etc/janus`. We need two main folders `apis` and `auth` each one of them
holds a set of configurations for our proxies. Just run:

```bash
mkdir -p /etc/janus/apis
mkdir -p /etc/janus/auth
```

## 2. Add your API

The main feature of the API Gateway is to proxy the requests to a different service, so let's do this.

Just place [this example](../../examples/apis/posts.json) in your `apis` directory.
This will create a proxy to `https://jsonplaceholder.typicode.com/posts` when you hit Janus on `GET /posts`.

Now restart Janus to apply the changes.

```bash
sudo sv restart janus
```

## 3. Verify that your API has been added

You can use the REST API to query all available APIs and Auth Providers. Simply make a request 
to `/apis`.

```bash
http -v GET localhost:8081/apis "Authorization:Bearer yourToken" "Content-Type: application/json"
```

## 4. Forward your requests through Janus

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
