package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/tidwall/gjson"
)

const (
	headerAuthorization = "Authorization"
)

// RegisterRequestContext registers godog suite context for handling HTTP-requests related steps
func RegisterRequestContext(s *godog.Suite, port, apiPort, portSecondary, apiPortSecondary int, adminCred config.Credentials) {
	ctx := &requestContext{
		port:             port,
		apiPort:          apiPort,
		portSecondary:    portSecondary,
		apiPortSecondary: apiPortSecondary,
		adminCred:        adminCred,
	}

	ctx.requestHeaders = make(http.Header)

	s.Step(`^I request "([^"]*)" path with "([^"]*)" method$`, ctx.iRequestPathWithMethod)
	s.Step(`^I request "([^"]*)" API path with "([^"]*)" method$`, ctx.iRequestAPIPathWithMethod)
	s.Step(`^I request "([^"]*)" secondary path with "([^"]*)" method$`, ctx.iRequestSecondaryPathWithMethod)
	s.Step(`^I request "([^"]*)" secondary API path with "([^"]*)" method$`, ctx.iRequestSecondaryAPIPathWithMethod)
	s.Step(`^I should receive (\d+) response code$`, ctx.iShouldReceiveResponseCode)
	s.Step(`^header "([^"]*)" should be "([^"]*)"$`, ctx.headerShouldBe)
	s.Step(`^header "([^"]*)" should start with "([^"]*)"$`, ctx.headerShouldStartWith)
	s.Step(`^the response should contain "([^"]*)"$`, ctx.responseShouldContain)
	s.Step(`^response JSON body has "([^"]*)" path with value \'([^']*)\'$`, ctx.responseJSONBodyHasPathWithValue)
	s.Step(`^response JSON body has "([^"]*)" path`, ctx.responseJSONBodyHasPath)
	s.Step(`^response JSON body is an array of length (\d+)$`, ctx.responseJSONBodyIsAnArrayOfLength)
	s.Step(`^request JSON payload \'([^']*)\'$`, ctx.requestJSONPayload)
	s.Step(`^request header "([^"]*)" is set to "([^"]*)"$`, ctx.requestHeaderIsSetTo)
	s.Step(`^request JWT token is not set$`, ctx.requestJWTTokenIsNotSet)
	s.Step(`^request JWT token is valid admin token$`, ctx.requestJWTTokenIsValidAdminToken)
}

type requestContext struct {
	port    int
	apiPort int

	portSecondary    int
	apiPortSecondary int

	adminCred config.Credentials

	requestBody    *bytes.Buffer
	requestHeaders http.Header
	response       *http.Response
	responseBody   []byte
}

func (c *requestContext) iRequestAPIPathWithMethod(path, method string) error {
	url := fmt.Sprintf("http://localhost:%d%s", c.apiPort, path)
	return c.doRequest(url, method)
}

func (c *requestContext) iRequestPathWithMethod(path, method string) error {
	url := fmt.Sprintf("http://localhost:%d%s", c.port, path)
	return c.doRequest(url, method)
}

func (c *requestContext) iRequestSecondaryAPIPathWithMethod(path, method string) error {
	url := fmt.Sprintf("http://localhost:%d%s", c.apiPortSecondary, path)
	return c.doRequest(url, method)
}

func (c *requestContext) iRequestSecondaryPathWithMethod(path, method string) error {
	url := fmt.Sprintf("http://localhost:%d%s", c.portSecondary, path)
	return c.doRequest(url, method)
}

func (c *requestContext) doRequest(url, method string) error {
	var req *http.Request
	var err error
	if method == http.MethodGet || method == http.MethodDelete {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, c.requestBody)
	}
	if nil != err {
		return fmt.Errorf("Failed to instantiate request instance: %v", err)
	}

	req.Header = c.requestHeaders

	// Inform to close the connection after the transaction is complete
	req.Header.Set("Connection", "close")

	c.response, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to perform request: %v", err)
	}

	c.responseBody, err = ioutil.ReadAll(c.response.Body)
	if nil != err {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	return nil
}

func (c *requestContext) iShouldReceiveResponseCode(code int) error {
	if c.response.StatusCode != code {
		return fmt.Errorf("expected response code %d, but actual is %d", code, c.response.StatusCode)
	}

	return nil
}

func (c *requestContext) headerShouldBe(header, value string) error {
	if actual := c.response.Header.Get(header); actual != value {
		return fmt.Errorf("expected header %s to be %s, but actual is %s", header, value, actual)
	}

	return nil

}

func (c *requestContext) headerShouldStartWith(header, value string) error {
	if !strings.HasPrefix(c.response.Header.Get(header), value) {
		actual := c.response.Header.Get(header)
		return fmt.Errorf("expected header %s to start with %s, but actual is %s", header, value, actual)
	}

	return nil

}

func (c *requestContext) responseShouldContain(text string) error {
	if !bytes.Contains(c.responseBody, []byte(text)) {
		return fmt.Errorf("expected response to contain %s, but actual is %s", text, string(c.responseBody))
	}
	return nil
}

func (c *requestContext) responseJSONBodyHasPathWithValue(path, value string) error {
	val := gjson.GetBytes(c.responseBody, path)
	if !val.Exists() {
		return fmt.Errorf("expected path %s in JSON response, but not found", path)
	}

	if val.String() != value {
		return fmt.Errorf("expected path %s in JSON response to be %s, but actual is %s; response: %s", path, value, val.String(), c.responseBody)
	}

	return nil
}

func (c *requestContext) responseJSONBodyHasPath(path string) error {
	val := gjson.GetBytes(c.responseBody, path)
	if !val.Exists() {
		return fmt.Errorf("expected path %s in JSON response, but not found", path)
	}

	return nil
}

func (c *requestContext) responseJSONBodyIsAnArrayOfLength(length int) error {
	var jsonResponse []interface{}
	err := json.Unmarshal(c.responseBody, &jsonResponse)
	if nil != err {
		return fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if len(jsonResponse) != length {
		return fmt.Errorf("expected JSON response array length is %d, but actual is %d", length, len(jsonResponse))
	}

	return nil
}

func (c *requestContext) requestJSONPayload(payload string) error {
	c.requestBody = bytes.NewBufferString(payload)
	return nil
}

func (c *requestContext) requestHeaderIsSetTo(header, value string) error {
	c.requestHeaders.Set(header, value)
	return nil
}

func (c *requestContext) requestJWTTokenIsNotSet() error {
	c.requestHeaders.Del(headerAuthorization)
	return nil
}

func (c *requestContext) requestJWTTokenIsValidAdminToken() error {
	jwtConfig := jwt.NewConfigWithHandlers(c.adminCred)
	accessToken, err := jwt.IssueAdminToken(jwtConfig.SigningAlgorithm, jwtConfig.SigningKey, c.adminCred.Username, jwtConfig.Timeout)
	if nil != err {
		return fmt.Errorf("Failed to issue JWT: %v", err)
	}

	c.requestHeaders.Set(headerAuthorization, "Bearer "+accessToken)

	return nil
}
