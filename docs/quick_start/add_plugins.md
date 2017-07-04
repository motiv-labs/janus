# Add Plugins

Any API Definition can make use of plugins. A plugin is a way to add functionalities to your definition, like rate limit, CORS, authentication, etc...

In this tutorial we will add a [rate limit](/docs/plugins/rate_limit.md) plugin to our definition.

## 1. Add the plugin

Let's create a file called `rate_limit.json` with our new additions to the API definition. We will create a rate limit rule that we can only send `5 req/m`
 just for us to play around with it.

```json
{
  "plugins": [
    {
      "name": "rate_limit",
      "enabled": true,
      "config": {
        "limit": "5-M",
        "policy": "local"
      }
    }
  ]
}
```

Now lets update our API definition:

```sh
http -v PUT localhost:8081/apis/my-endpoint "Authorization:Bearer yourToken" "Content-Type: application/json" < rate_limit.json
```

Done! Now Janus already reloaded the configuration and the rate limit is enabled on your endpoint.

## 2. Check if the plugin is working

Lets make a few requests to our endpoint and see if the rate limit is working.

```bash
$ http -v GET http://localhost:8080/example
```

You should see a few extra headers on the response of this request:

```
X-Ratelimit-Limit →5
X-Ratelimit-Remaining →4
X-Ratelimit-Reset →1498773715
```

This means the plugin is working properly. Now lets make the same request 4 more times... On the fifth time you should get:

```
Status Code: 429 Too Many Requests

Limit exceeded
```

After 1 minute you should be able to make 5 more requests :)

In the [next part](add_auth.md) we'll learn how to protect our endpoint.
