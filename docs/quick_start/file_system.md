# Adding your API - File System

By choosing a File System based configuration we have a static way of configure Janus (similar to nginx).

## 1. Boot it up

We highly recommend you to use one of our examples to start. Let's see the [front-proxy](/examples/front-proxy) example:

Make sure you have docker up and running on your platform and then run.

```sh
docker-compose up -d
```

This will spin up a janus server and will have a small proxy configuration that is going to a mock server that we spun up.

## 2. Verify that Janus is working

Issue the following cURL request to verify that Janus is properly forwarding
requests to your API. Note that [by default][proxy-port] Janus handles proxy
requests on port `:8080`:

If you access `http://localhost:8080/example` you should something like:

```json
{
    "message": "Hello World!"
}
```

A successful response means Janus is now forwarding requests made to
`http://localhost:8080` to the elected upstream target (chosen by the load balancer) we configured in step #1,
and is forwarding the response back to us.

## Understanding the directory structure

By default all apis configurations are splitted in separated files (both single and multiple api definitions per file are supported) and they are stored in `/etc/janus`. You can change that path by simply defining the configuration `database.dsn`, for instance, you can define the value to `file:///usr/local/janus`.

There are two required folder that needs to be there:

- `/etc/janus/apis` - Holds all API definitions
- `/etc/janus/auth` - Holds all your Auth servers configurations

## 4. Adding a new endpoint and authentication

To add a new endpoint or authentication you can see the [Add Endpoint tutorial](add_endpoint.md) but instead of using the admin API you'll add your configuration to a file and reload the docker instance:

```sh
docker-compose reload janus
```
