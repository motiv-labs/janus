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
