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


INSERT INTO janus.api_definition (name, definition) VALUES ('private', '{"name":"private","active":true,"proxy":{"preserve_host":false,"listen_path":"/private/{service}/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://{service}-app:8080/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","DELETE","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE","OPTIONS"],"request_headers":["X-Requested-With","Authorization","Content-Type"]}},{"name":"basic_auth","enabled":true,"config":null}],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.api_definition (name, definition) VALUES ('public-github-webhook-listener', '{"name":"public-github-webhook-listener","active":true,"proxy":{"preserve_host":false,"listen_path":"/public/github-webhook-listener/webhook/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://github-webhook-listener-app:8080/webhook/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.api_definition (name, definition) VALUES ('public-taxi-dropoff-strapi', '{"name":"public-taxi-dropoff-strapi","active":true,"proxy":{"preserve_host":false,"listen_path":"/public/taxi-dropoff-strapi/binary/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://taxi-dropoff-strapi-app:8080/binary/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.api_definition (name, definition) VALUES ('public-taxi-dropoff-dotcms', '{"name":"public-taxi-dropoff-dotcms","active":true,"proxy":{"preserve_host":false,"listen_path":"/public/taxi-dropoff-dotcms/binary/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://taxi-dropoff-dotcms-app:8080/binary/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.api_definition (name, definition) VALUES ('public-taxi-pickup-dotcms', '{"name":"public-taxi-pickup-dotcms","active":true,"proxy":{"preserve_host":false,"listen_path":"/public/taxi-pickup-dotcms/diffuseContentPayload/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://taxi-pickup-dotcms-app:8080/diffuseContentPayload/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.api_definition (name, definition) VALUES ('public-taxi-pickup-strapi', '{"name":"public-taxi-pickup-strapi","active":true,"proxy":{"preserve_host":false,"listen_path":"/public/taxi-pickup-strapi/diffuseContentPayload/*","upstreams":{"balancing":"roundrobin","targets":[{"target":"http://taxi-pickup-strapi-app:8080/diffuseContentPayload/","weight":0}]},"insecure_skip_verify":false,"strip_path":true,"append_path":false,"methods":["GET","POST","PUT","OPTIONS"],"hosts":[],"forwarding_timeouts":{"dial_timeout":"0s","response_header_timeout":"0s"}},"plugins":[],"health_check":{"url":"","timeout":0}}');
INSERT INTO janus.user (username, password) VALUES ('user', '$2a$10$SWpROghFjIrxUfESOaJ/pu/SwLqPZP.8hTuMpT3Iq4XWIk/7mWJV6');
