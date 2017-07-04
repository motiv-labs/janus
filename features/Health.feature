Feature: Run health-check against registered proxies.

    Background:
        Given request JWT token is valid admin token

    Scenario: Request health-checks statuses
        Given request JSON payload '{"name":"example-ok","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-ok/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-ok"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-ok"

        When I request "/status" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'OK'
        And response JSON body has "timestamp" path

        Given request JSON payload '{"name":"example-ok2","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-ok2/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-ok"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-ok2"

        When I request "/status" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'OK'
        And response JSON body has "timestamp" path

        Given request JSON payload '{"name":"example-partial","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-partial/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-partial"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-partial"

        When I request "/status" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'Partially Available'
        And response JSON body has "timestamp" path
        And response JSON body has "failures.example-partial" path

        Given request JSON payload '{"name":"example-broken","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-broken/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-broken"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-broken"

        When I request "/status" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'Partially Available'
        And response JSON body has "timestamp" path
        And response JSON body has "failures.example-partial" path
        And response JSON body has "failures.example-broken" path

    Scenario: Request health-check service status
        Given request JSON payload '{"name":"example-partial","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-partial/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-partial"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-partial"

        When I request "/status/example-partial" API path with "GET" method
        Then I should receive 400 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'Partially Available'
        And response JSON body has "timestamp" path
        And response JSON body has "failures.rabbitmq" path with value 'Failed during RabbitMQ health check'

        Given request JSON payload '{"name":"example-broken","active":true,"proxy":{"preserve_host":false,"listen_path":"/example-broken/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"http://localhost:9089/status-broken"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example-broken"

        When I request "/status/example-broken" API path with "GET" method
        Then I should receive 503 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'Unavailable'
        And response JSON body has "timestamp" path
        And response JSON body has "failures.mongodb" path with value 'Failed during MongoDB health check'

        When I request "/status/does-not-exist" API path with "GET" method
        Then I should receive 404 response code
        And the response should contain "Definition name is not found"
