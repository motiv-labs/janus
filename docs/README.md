> An API Gateway written in Go EDITED

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
