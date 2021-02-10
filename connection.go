package digitalstrom

import (
	"crypto/tls"
	"encoding/json"
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

// doRequest Performing an Http Request to the given url using the given method and (optionally) sends the body
func (c *Connection) doRequest(url string, method requestMethod, body string, params map[string]string) (*RequestResult, error) {

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

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// TODO: generate a better error message according to StatusCode
	if res.StatusCode != http.StatusOK {

		return nil, &RequestError{"Unexpected Resonse for " + url, res.StatusCode}
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

func (c *Connection) checkApplicationToken() bool {
	// ToDo do propper checks
	return c.ApplicationToken != ""
}

func (c *Connection) checkAccessToken() bool {
	// ToDo do propper checks
	return c.SessionToken != ""
}
