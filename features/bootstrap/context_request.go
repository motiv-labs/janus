package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
)

const (
	headerAuthorization = "Authorization"
)

func RegisterRequestContext(s *godog.Suite, port, apiPort int) {
	ctx := &requestContext{port: port, apiPort: apiPort}

	ctx.requestHeaders = make(http.Header)

	s.Step(`^I request "([^"]*)" path with "([^"]*)" method$`, ctx.iRequestPathWithMethod)
	s.Step(`^I request "([^"]*)" API path with "([^"]*)" method$`, ctx.iRequestAPIPathWithMethod)
	s.Step(`^I should receive (\d+) response code$`, ctx.iShouldReceiveResponseCode)
	s.Step(`^header "([^"]*)" should be "([^"]*)"$`, ctx.headerShouldBe)
	s.Step(`^header "([^"]*)" should start with "([^"]*)"$`, ctx.headerShouldStartWith)
	s.Step(`^the response should contain "([^"]*)"$`, ctx.responseShouldContain)
	s.Step(`^response JSON body has "([^"]*)" field with value \'([^']*)\'$`, ctx.responseJSONBodyHasFieldWithValue)
	s.Step(`^response JSON body has "([^"]*)" field$`, ctx.responseJSONBodyHasField)
	s.Step(`^response JSON body is an array of length (\d+)$`, ctx.responseJSONBodyIsAnArrayOfLength)
	s.Step(`^response JSON body array element (\d+) has "([^"]*)" field with value \'([^']*)\'$`, ctx.responseJSONBodyArrayElementHasFieldWithValue)
	s.Step(`^request JSON payload \'([^']*)\'$`, ctx.requestJSONPayload)
	s.Step(`^request header "([^"]*)" is set to "([^"]*)"$`, ctx.requestHeaderIsSetTo)
}

type requestContext struct {
	port    int
	apiPort int

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

func (c *requestContext) doRequest(url, method string) error {
	var req *http.Request
	var err error
	if method == http.MethodGet || method == http.MethodDelete {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, c.requestBody)
	}
	if nil != err {
		return err
	}

	req.Header = c.requestHeaders

	c.response, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.responseBody, err = ioutil.ReadAll(c.response.Body)
	if nil != err {
		return err
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

func (c *requestContext) responseJSONBodyHasFieldWithValue(field, value string) error {
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(c.responseBody, &jsonResponse)
	if nil != err {
		return err
	}

	fieldValue, ok := jsonResponse[field]
	if !ok {
		return fmt.Errorf("expected field %s in JSON response, but not found", field)
	}

	// convert interface{} to string simply by marshalling to json
	valueBytes, err := json.Marshal(fieldValue)
	if err != nil {
		return err
	}

	if string(valueBytes) != value {
		return fmt.Errorf("expected field %s in JSON response to be %s, but actual is %s", field, value, string(valueBytes))
	}

	return nil
}

func (c *requestContext) responseJSONBodyHasField(field string) error {
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(c.responseBody, &jsonResponse)
	if nil != err {
		return err
	}

	_, ok := jsonResponse[field]
	if !ok {
		return fmt.Errorf("expected field %s in JSON response, but not found", field)
	}

	return nil
}

func (c *requestContext) responseJSONBodyIsAnArrayOfLength(length int) error {
	var jsonResponse []interface{}
	err := json.Unmarshal(c.responseBody, &jsonResponse)
	if nil != err {
		return err
	}

	if len(jsonResponse) != length {
		return fmt.Errorf("expected JSON response array length is %d, but actual is %d", length, len(jsonResponse))
	}

	return nil
}

func (c *requestContext) responseJSONBodyArrayElementHasFieldWithValue(idx int, field, value string) error {
	var jsonResponse []map[string]interface{}
	err := json.Unmarshal(c.responseBody, &jsonResponse)
	if nil != err {
		return err
	}

	elementToTest := jsonResponse[idx]
	fieldValue, ok := elementToTest[field]
	if !ok {
		return fmt.Errorf("expected field %s in JSON response array @ idx %d, but not found", field, idx)
	}

	// convert interface{} to string simply by marshalling to json
	valueBytes, err := json.Marshal(fieldValue)
	if err != nil {
		return err
	}

	if string(valueBytes) != value {
		return fmt.Errorf("expected field %s in JSON response array @ idx %d to be %s, but actual is %s", field, idx, value, string(valueBytes))
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
