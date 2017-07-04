Feature: Proxy requests to upstream.

    Background:
        Given request JWT token is valid admin token

    Scenario: Proxy request to existing mock service
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://localhost:9089/hello-world","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com/status"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        Given request JSON payload '{"name":"posts","active":true,"proxy":{"preserve_host":false,"listen_path":"/posts/*","upstream_url":"http://localhost:9089/posts","strip_path":true,"append_path":false,"enable_load_balancing":false,"methods":["ALL"],"hosts":["hellofresh.*"]},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE"],"request_headers":["Origin","Authorization","Content-Type"],"exposed_headers":["X-Debug-Token","X-Debug-Token-Link"]}},{"name":"rate_limit","enabled":true,"config":{"limit":"10-S","policy":"local"}},{"name":"oauth2","enabled":true,"config":{"server_name":"local"}},{"name":"compression","enabled":true}]}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/posts"

        Given request JSON payload '{"name":"posts-public","active":true,"proxy":{"preserve_host":false,"listen_path":"/posts-public/*","upstream_url":"http://localhost:9089/posts","strip_path":true,"append_path":false,"enable_load_balancing":false,"methods":["ALL"]},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE"],"request_headers":["Origin","Authorization","Content-Type"],"exposed_headers":["X-Debug-Token","X-Debug-Token-Link"]}},{"name":"rate_limit","enabled":true,"config":{"limit":"10-S","policy":"local"}},{"name":"oauth2","enabled":true,"config":{"server_name":"local"}},{"name":"compression","enabled":true}]}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/posts-public"

        When I request "/example" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/example/nested/path" path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "hello" path with value 'world'

        When I request "/posts" path with "GET" method
        Then I should receive 404 response code
        And the response should contain "no API found with those values"

        When I request "/posts-public" path with "GET" method
        Then I should receive 200 response code
        And response JSON body is an array of length 100

        When I request "/posts-public/nested/path" path with "GET" method
        Then I should receive 404 response code
