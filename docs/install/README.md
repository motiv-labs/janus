# Installation

You can install Janus using docker or manually deploying the binary.

### Docker

The simplest way of installing janus is to run the docker image for it. Just check the [docker-compose.yml](/examples/front-proxy/docker-compose.yml) example and then run it.

```sh
docker-compose up -d
```

Now you should be able to get a response from the gateway. 

Try the following command:

```sh
http http://localhost:8080/
```

You can find more details about how to use Janus docker image [here](docker.md).

### Manual

You can get the binary and play with it in your own environment (or even deploy it where ever you like).
Just go to the [releases](https://github.com/hellofresh/janus/releases) page and download the latest one for your platform.
