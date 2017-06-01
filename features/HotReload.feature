Feature: Apply proxy changes to all instances, when changing via API on single instance.

    Background:
        Given request JWT token is valid admin token

    Scenario: Proxy registered on the primary instance is available on the secondary instance
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://www.mocky.io/v2/58c6c60710000040151b7cad","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I wait for a while
        And I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I wait for a while
        And I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

    Scenario: Proxy registered on the secondary instance is available on the primary instance
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://www.mocky.io/v2/58c6c60710000040151b7cad","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" secondary API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I wait for a while
        And I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I wait for a while
        And I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

    Scenario: Proxy removed on the primary instance is not available on the secondary instance
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://www.mocky.io/v2/58c6c60710000040151b7cad","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I wait for a while
        And I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/apis/example" API path with "DELETE" method
        Then I should receive 204 response code

        When I wait for a while
        And I request "/example" path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I wait for a while
        And I request "/example" secondary path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I request "/example" path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I request "/example" secondary path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

    Scenario: Proxy removed on the secondary instance is not available on the primary instance
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://www.mocky.io/v2/58c6c60710000040151b7cad","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" secondary API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        When I request "/example" secondary path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I wait for a while
        And I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/apis/example" secondary API path with "DELETE" method
        Then I should receive 204 response code

        When I wait for a while
        And I request "/example" secondary path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I wait for a while
        And I request "/example" path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I request "/example" secondary path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'

        When I request "/example" path with "GET" method
        Then I should receive 404 response code
        And response JSON body has "error" path with value 'no API found with those values'
