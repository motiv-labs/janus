Feature: Login to API management section of the system.

    Scenario: Try to login with valid credentials
        Given request JSON payload:
            """
            {
                "username":"admin",
                "password":"admin"
            }
            """
        And request header "Content-Type" is set to "application/json"
        When I request "/login" API path with "POST" method
        Then I should receive 200 response code
        And the response should contain "token"
        And the response should contain "expire"

    Scenario: Try to login with invalid credentials
        Given request JSON payload:
            """
            {
                "username":"admin111",
                "password":"admin111"
            }
            """
        And request header "Content-Type" is set to "application/json"
        When I request "/login" API path with "POST" method
        Then I should receive 401 response code

    Scenario: Try to login with missed credentials
        Given request header "Content-Type" is set to "application/json"
        And request JSON payload:
            """
            {
                "username":"admin111"
            }
            """
        When I request "/login" API path with "POST" method
        Then I should receive 401 response code

        Given request header "Content-Type" is set to "application/json"
        And request JSON payload:
            """
            {
                "password":"admin111"
            }
            """
        When I request "/login" API path with "POST" method
        Then I should receive 401 response code

        Given request header "Content-Type" is set to "application/json"
        And request JSON payload:
            """
            {}
            """
        When I request "/login" API path with "POST" method
        Then I should receive 401 response code
