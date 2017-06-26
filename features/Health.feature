Feature: Run health-check against registered proxies.

    Background:
        Given request JWT token is valid admin token

    Scenario: Request health-checks statuses
        Given request JSON payload '{"name":"example","active":true,"proxy":{"preserve_host":false,"listen_path":"/example/*","upstream_url":"http://www.mocky.io/v2/58c6c60710000040151b7cad","strip_path":false,"append_path":false,"enable_load_balancing":false,"methods":["GET"]},"health_check":{"url":"https://example.com"}}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/example"

        Given request JSON payload '{"name":"posts","active":true,"proxy":{"preserve_host":false,"listen_path":"/posts/*","upstream_url":"https://jsonplaceholder.typicode.com/posts","strip_path":true,"append_path":false,"enable_load_balancing":false,"methods":["ALL"],"hosts":["hellofresh.*"]},"plugins":[{"name":"cors","enabled":true,"config":{"domains":["*"],"methods":["GET","POST","PUT","PATCH","DELETE"],"request_headers":["Origin","Authorization","Content-Type"],"exposed_headers":["X-Debug-Token","X-Debug-Token-Link"]}},{"name":"rate_limit","enabled":true,"config":{"limit":"10-S","policy":"local"}},{"name":"oauth2","enabled":true,"config":{"server_name":"local"}},{"name":"compression","enabled":true}]}'
        When I request "/apis" API path with "POST" method
        Then I should receive 201 response code
        And header "Location" should be "/apis/posts"

        When I request "/status" API path with "GET" method
        Then I should receive 200 response code
        And response JSON body has "system" path
        And response JSON body has "status" path with value 'Partially Available' 
        And response JSON body has "timestamp" path
