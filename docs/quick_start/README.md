# Quick Start

1. [Download](#download)
2. [Install](#install)
3. [Run](#run)
3. [Configure](#configure)

## Download

You can get Janus for nearly any OS and architecture. You can get the latest Janus release on [Github](https://github.com/hellofresh/janus/releases).

## Install and run

We highly recommend you to use one of our examples to start. Let's see the [front-proxy-mongo](https://github.com/hellofresh/janus/blob/master/examples/front-proxy-mongo) example:

Make sure you have docker up and running on your platform and then run.

```sh
docker-compose up -d
```

This will spin up a janus server and will have a small proxy configuration that is going to a mock server that we spun up.

## Configure

If you access `http://localhost:8080/example` you should see something like:

```json
{
    "message": "Hello World!"
}
```

That means that Janus already proxied your request to an upstream. But of course you don't just want to do that. For this reason
now is the perfect time for you to learn about all the available configurations that you can play with.

> Note: If you are using the file-based configuration you will not be able to use the write administration API to add/modify/remove new endpoints, plugins, etc. Please check this [tutorial](file_system.md) if you'd like to add a new endpoint using the file-based configuration.

Next, let's learn about how to [configure a new endpoint](authenticating.md).
