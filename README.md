<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# Janus

> An Api Gateway written in Go

This is a lightweight, API Gateway and Management Platform enables you to control who accesses your API, 
when they access it and how they access it. API Gateway will also record detailed analytics on how your 
users are interacting with your API and when things go wrong. 

## Why Janus?

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Ianus, pronounced [ˈjaː.nus]) is the god of beginnings, 
gates, transitions, time, doorways, passages, and endings. He is usually depicted as having two faces, since he 
looks to the future and to the past. [Wikipedia](https://en.wikipedia.org/wiki/Janus)

We thought it would be nice to name the project after the God of the Gates :smile:

## What is an API Gateway?

An API Gateway sits in front of your application(s) and/or services and manages the heavy lifting of authorisation, 
access control and throughput limiting to your services. Ideally, it should mean that you can focus on creating 
services instead of implementing management infrastructure. For example if you have written a really awesome 
web service that provides geolocation data for all the cats in NYC, and you want to make it public, 
integrating an API gateway is a faster, more secure route that writing your own authorisation middleware.

## Key Features

This API Gateway offers powerful, yet lightweight features that allow fine gained control over your API ecosystem.

* **RESTFul API** - Full programatic access to the internals makes it easy to manage your API users, keys and Api Configuration from within your systems
* **Multiple access protocols** - Out of the box, we support JWT, OAtuh2, Basic Auth and Keyless access methods
* **Rate Limiting** - Easily rate limit your API users, rate limiting is granular and can be applied on a per-key basis
* **API Versioning** - API Versions can be easily set and deprecated at a specific time and date
* **Analytics logging** - Record detailed usage data on who is using your API's (raw data only)

## Installation

### Docker

The simplest way of installing janus is to run the docker image for it. Just check the [docker-compose.yml](ci/assets/docker-compose.yml)
example and then run.

```sh
docker-compose up -d
```

Now you should be able to get a response from the gateway, try the following command:

```sh
curl http://localhost:8080/
```

### Manual

You can get the binary and play with it in your own enviroment (or even deploy it whereever you like it).
Just go the [releases](https://github.com/hellofresh/janus/releases) and download the latest one for your platform.

Make sure you have the following dependencies installed in your enviroment:
 - Mongodb - For storing the proxies configurations
 - Redis - For caching and storing of expiration oauth tokens

And then just define where your dependencies are located

```sh
export DATABASE_DSN="mongodb://localhost:27017/janus"
export REDIS_DSN="redis://localhost:6379"
```

If you want you can have a stastd server so you can have some metrics about your gateway. For that just define:

```sh
export STATSD_DSN="statsd:8125"
```

## Getting Started

After you have *janus* up and running we need to setup our first proxy. Everything that we want to do on the gateway 
we do it through a REST API, since all endpoints are protected we need to login first.

```sh
http -v --json POST localhost:3000/login username=admin password=admin
```

The username and password are defined in a enviroment variable called `ADMIN_USERNAME` and `ADMIN_PASSWORD`, it defaults to *admin*/*admin*.

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>


### Creating a proxy

The main feature of the API Gateway is to proxy the requests to a different service, so lets do this.
Now that you are authenticate you can send a request to `/apis` to create a proxy.

```
http -v --json POST localhost:3000/apis "Authorization:Bearer yourToken" "Content-Type: application/json" < examples/api.json
```

<p align="center">
  <a href="http://g.recordit.co/Hi7SX8s5IA.gif">
    <img src="http://g.recordit.co/Hi7SX8s5IA.gif">
  </a>
</p>

This will create a proxy to `https://jsonplaceholder.typicode.com/posts` when you hit the api gateway on `GET /posts`.
Now just make a request to `GET /posts`

```
http -v --json GET http://localhost:3000/posts/1
```
<p align="center">
  <a href="http://g.recordit.co/vufeMjwEfg.gif">
    <img src="http://g.recordit.co/vufeMjwEfg.gif">
  </a>
</p>

Done! You just made your first request through the gateway.

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation
* Go lang: https://golang.org/
