CREATE KEYSPACE IF NOT EXISTS janus with replication  = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE janus;

CREATE TABLE IF NOT EXISTS janus.user (
    username text,
    password text,
    PRIMARY KEY (username));

CREATE TABLE IF NOT EXISTS janus.api_definition (
    name text,
    definition text,
    PRIMARY KEY (name));

CREATE TABLE IF NOT EXISTS janus.oauth (
    name text,
    oauth text,
    PRIMARY KEY (name));

CREATE TABLE IF NOT EXISTS janus.organization (
    username text,
    password text,
    organization text,
    PRIMARY KEY (username));

CREATE TABLE IF NOT EXISTS janus.organization_config (
    organization text,
    priority int,
    content_per_day int,
    config text,
    PRIMARY KEY (organization));
