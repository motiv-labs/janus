<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# Janus

[![Build Status](https://travis-ci.org/hellofresh/janus.svg?branch=master)](https://travis-ci.org/hellofresh/janus)
[![Coverage Status](https://coveralls.io/repos/github/hellofresh/janus/badge.svg?branch=feature%2Fopentracing)](https://coveralls.io/github/hellofresh/janus?branch=feature%2Fopentracing)
[![GoDoc](https://godoc.org/github.com/hellofresh/janus?status.svg)](https://godoc.org/github.com/hellofresh/janus)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/janus)](https://goreportcard.com/report/github.com/hellofresh/janus)

> An API Gateway written in Go

This is a lightweight API Gateway and Management Platform that enables you to control who accesses your API,
when they access it and how they access it. API Gateway will also record detailed analytics on how your
users are interacting with your API and when things go wrong.

## Why Janus?

> In ancient Roman religion and myth, Janus (/ˈdʒeɪnəs/; Latin: Ianus, pronounced [ˈjaː.nus]) is the god of beginnings,
gates, transitions, time, doorways, passages, and endings. He is usually depicted as having two faces since he
looks to the future and to the past. [Wikipedia](https://en.wikipedia.org/wiki/Janus)

We thought it would be nice to name the project after the God of the Gates :smile:

## What is an API Gateway?

An API Gateway sits in front of your application(s) and/or services and manages the heavy lifting of authorisation,
access control and throughput limiting to your services. Ideally, it should mean that you can focus on creating
services instead of implementing management infrastructure. For example, if you have written a really awesome
web service that provides geolocation data for all the cats in NYC, and you want to make it public,
integrating an API gateway is a faster, more secure route than writing your own authorisation middleware.

## Key Features

This API Gateway offers powerful, yet lightweight features that allows fine gained control over your API ecosystem.

* No dependency hell, single binary made with go
* [OpenTracing](http://opentracing.io/) support for Distributed tracing (Supports Google Cloud Platform, Zipkin and Appdash)
* HTTP/2 support
* REST API, full programatic access to the internals makes it easy to manage your API users, keys and API Configuration from within your systems
* Rate Limiting, easily rate limit your API users, rate limiting is granular and can be applied on a per-key basis
* CORS Filter, enable cors for your API, or even for specific endpoints
* Multiple auth protocols, out of the box, we support JWT, OAuth 2.0 and Basic Auth access methods
* Graceful shutdown of http connections
* Small [official](https://quay.io/repository/hellofresh/janus) docker image included

## Installation

### Docker

The simplest way of installing janus is to run the docker image for it. Just check the [docker-compose.yml](ci/assets/docker-compose.yml)
example and then run it.

```sh
docker-compose up -d
```

Now you should be able to get a response from the gateway. 

Try the following command:

```sh
http http://localhost:8080/
```

### Manual

You can get the binary and play with it in your own environment (or even deploy it where ever you like).
Just go to the [releases](https://github.com/hellofresh/janus/releases) page and download the latest one for your platform.

## Getting Started

After you have *janus* up and running we need to setup our first proxy. You can choose between:

* [File System](https://hellofresh.gitbooks.io/janus/quick_start/file_system.html)
* [MongoDB](https://hellofresh.gitbooks.io/janus/quick_start/mongodb.html)

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation

* Janus Docs: https://hellofresh.gitbooks.io/janus
* Janus Go Docs: https://godoc.org/github.com/hellofresh/janus
* Go lang: https://golang.org/
