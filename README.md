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

## What is an API Gateway?

An API Gateway sits in front of your application(s) and/or services and manages the heavy lifting of authorisation, 
access control and throughput limiting to your services. Ideally, it should mean that you can focus on creating 
services instead of implementing management infrastructure. For example if you have written a really awesome 
web service that provides geolocation data for all the cats in NYC, and you want to make it public, 
integrating an API gateway is a faster, more secure route that writing your own authorisation middleware.

## Key Features

This API Gateway offers powerful, yet lightweight features that allow fine gained control over your API ecosystem.

* **RESTFul API** - Full programatic access to the internals makes it easy to manage your API users, keys and Api Configuration from within your systems
* **Multiple access protocols** - Out of the box, we support Token-based, Basic Auth and Keyless access methods
* **Rate Limiting** - Easily rate limit your API users, rate limiting is granular and can be applied on a per-key basis
* **API Versioning** - API Versions can be easily set and deprecated at a specific time and date
* **Analytics logging** - Record detailed usage data on who is using your API's (raw data only)

The API Gateway is written in Go, which makes it fast and easy to set up. Its only dependencies are a Mongo database (for analytics) and Redis, though it can be deployed without either (not recommended).

## Installation

### Docker

The simplest way of installing janus is to run the docker image for it. Just check the [docker-compose.yml](ci/assets/docker-compose.yml)
example and then run.

```sh
docker-compose up -d
```

Now you should be able to get a response from the gateway, try the following command:

```sh
curl http://localhost:8080/apis/
```

### Manual

You can get the binary and play with it in your own enviroment (or even deploy it whereever you like it).

```sh
go get github.com/hellofresh/janus
```

Make sure you have the following dependencies installed in your enviroment:
 - Mongodb - For storing the proxies configurations
 - Redis - For caching and storing of expiration oauth tokens

And then just define where your dependencies are located

```sh
export DATABASE_DSN: 'mongodb://localhost:27017/janus'
export REDIS_DSN: 'redis://localhost:6379'
```

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation
* Go lang: https://golang.org/
