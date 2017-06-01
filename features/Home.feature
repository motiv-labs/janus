Feature: Retrieve welcome line of the service.
    As an anonymous user, I need to be able to see welcome line of the system.
    In order to know that the service is up and running.

    Scenario: Check welcome line of the service
        When I request "/" API path with "GET" method
        Then I should receive 200 response code
        And header "Content-Type" should be "application/json"
        And the response should contain "Welcome to Janus v"
