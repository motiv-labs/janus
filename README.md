# hellofresh/api-gateway

# owner
https://github.com/italolelis

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
* **Multiple access protocols** - Out of the box, Tyk supports Token-based, HMAC Signed, Basic Auth and Keyless access methods
* **Rate Limiting** - Easily rate limit your API users, rate limiting is granular and can be applied on a per-key basis
* **Quotas** - Enforce usage quotas on users to manage capacity or charge for tiered access
* **API Versioning** - API Versions can be easily set and deprecated at a specific time and date
Blacklist/Whitelist/Ignored endpoint access - Enforce strict security models on a version-by-version basis to your access points
* **Analytics logging** - Record detailed usage data on who is using your API's (raw data only)

The API Gateway is written in Go, which makes it fast and easy to set up. Its only dependencies are a Mongo database (for analytics) and Redis, though it can be deployed without either (not recommended).

## Contributing

To start contributing, please check https://github.com/hellofresh/api-gateway/blob/master/.github/CONTRIBUTING.md.

## Documentation
* Go lang: https://golang.org/
