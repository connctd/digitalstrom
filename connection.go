package digitalstrom

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type requestMethod string

const (
	// DEFAULT_BASE_URL Will be set as default Base Url for connection
	DEFAULT_BASE_URL = "https://192.168.178.178:8080"
)

type RequestError struct {
	Err        string
	StatusCode int
	Request    http.Request
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("connection error - status code &status", strconv.Itoa(e.StatusCode))
}

type RequestResult struct {
	OK      bool                   `json:"ok"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result"`
}

const (
	get    requestMethod = "GET"
	put    requestMethod = "PUT"
	post   requestMethod = "POST"
	delete requestMethod = "DELETE"
)

// Connection Holds access credentials and URL information for a specific account. Contains routines for GET, PUSH, PUT and DELETE
type Connection struct {
	SessionToken     string
	ApplicationToken string
	BaseURL          string
	HTTPClient       http.Client
}

// Get Performs a GET request and returns the response body as string
func (c *Connection) Get(url string) (*RequestResult, error) {
	return c.doRequest(url, get, "", nil)
}

// Post Performs a Post Request with the given content and returns the response body as string
func (c *Connection) Post(url string, body string) (*RequestResult, error) {
	return c.doRequest(url, post, body, nil)
}

// Delete Performs a Delete Request and returns the response body as string
func (c *Connection) Delete(url string) (*RequestResult, error) {
	return c.doRequest(url, delete, "", nil)
}

// Put Performs a Put Request with the given content and returns the response body as string
func (c *Connection) Put(url string, body string) (*RequestResult, error) {
	return c.doRequest(url, put, body, nil)
}

func (c *Connection) generateHTTPRequest(url string, method requestMethod, body string, params map[string]string) (*http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var req *http.Request
	var err error

	if len(body) > 0 {
		req, err = http.NewRequest(string(method), url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(string(method), url, nil)
	}
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if c.checkAccessToken() {
		q.Add("token", c.SessionToken)
	}
	if params != nil {
		for p, v := range params {
			q.Add(p, v)
		}
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

// doRequest Performing an Http Request to the given url using the given method and (optionally) sends the body
func (c *Connection) doRequest(url string, method requestMethod, body string, params map[string]string) (*RequestResult, error) {

	req, err := c.generateHTTPRequest(url, method, body, params)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// TODO: generate a better error message according to StatusCode
	if res.StatusCode != http.StatusOK {
		return nil, &RequestError{"Unexpected Resonse for " + url, res.StatusCode, *req}
	}
	if res.Header.Get("Content-Length") != "0" {
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return convertToRequestResult(resp)
	}

	return nil, err
}

func convertToRequestResult(body []byte) (*RequestResult, error) {
	var requestResult RequestResult

	err := json.Unmarshal(body, &requestResult)
	if err != nil {
		return nil, err
	}

	return &requestResult, nil
}

func (c *Connection) applicationLogin() error {
	if !c.checkApplicationToken() {
		return errors.New("applicationToken is not set")
	}
	params := map[string]string{"loginToken": c.ApplicationToken}
	res, err := c.doRequest(c.BaseURL+"/json/system/loginApplication", get, "", params)
	if err != nil {
		return err
	}

	if res.OK {
		c.SessionToken = res.Result["token"].(string)
		return nil
	}
	return errors.New(res.Message)
}

// register an application with the given applicationName. Performs a request to generate an application token. A second request requires the
// Username and Password in order to generate a temporary session token. A third request enables the application token to login without
// further user credentials (applicationLogin). Returns the application token or an error. The application token will not be assigned automatically.
// Thus, in order to use the generated application token, it has to be set afterwards.
func (c *Connection) register(username string, password string, applicationName string) (string, error) {
	// request an ApplicationToken
	res, err := c.doRequest(c.BaseURL+"/json/system/requestApplicationToken", get, "", map[string]string{"applicationName": applicationName})
	if err != nil {
		return "", err
	}
	if !res.OK {
		return "", errors.New(res.Message)
	}
	applicationToken := res.Result["applicationToken"].(string)

	// performing a login in order to generate a temporary session token
	// HINT during the http request generation, parameters will be URL-encoded. This leads to a problem with the password
	// when it contains symbols that will change during the encoding problem. The digitalstrom server is not performing an
	// URL decoding to restore the password
	// This is a workaround, that might fail with passwords which are not url compatible
	res, err = c.doRequest(c.BaseURL+"/json/system/login?user="+username+"&password="+password, get, "", nil)
	if err != nil {
		return "", err
	}
	if !res.OK {
		return "", errors.New(res.Message)
	}

	sessionToken := res.Result["token"].(string)

	// use the session token to enable the application token. Future logins wont need user credentials anymore, the application token will be used to
	// perform an application login
	res, err = c.doRequest(c.BaseURL+"/json/system/enableToken", get, "", map[string]string{"applicationToken": applicationToken, "token": sessionToken})
	if err != nil {
		return "", err
	}

	return applicationToken, nil
}

func (c *Connection) checkApplicationToken() bool {
	// ToDo do propper checks
	return c.ApplicationToken != ""
}

func (c *Connection) checkAccessToken() bool {
	// ToDo do propper checks
	return c.SessionToken != ""
}
