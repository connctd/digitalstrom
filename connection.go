package digitalstrom

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type requestMethod string

const (
	// DEFAULT_BASE_URL Will be set as default Base Url for connection
	DEFAULT_BASE_URL = "https://192.168.178.178:8080"
)

const (
	get    requestMethod = "GET"
	put    requestMethod = "PUT"
	post   requestMethod = "POST"
	delete requestMethod = "DELETE"
)

// Connection Holds access credentials and URL information for a specific account. Contains routines for GET, PUSH, PUT and DELETE
type Connection struct {
	AccessToken string
	BaseURL     string
	httpClient  http.Client
}

// Get Performs a GET request and returns the response body as string
func (c *Connection) Get(url string) (string, error) {
	return c.doRequest(url+"?token="+c.AccessToken, get, "")
}

// Post Performs a Post Request with the given content and returns the response body as string
func (c *Connection) Post(url string, body string) (string, error) {
	return c.doRequest(url, post, body)
}

// Delete Performs a Delete Request and returns the response body as string
func (c *Connection) Delete(url string) (string, error) {
	return c.doRequest(url, delete, "")
}

// Put Performs a Put Request with the given content and returns the response body as string
func (c *Connection) Put(url string, body string) (string, error) {
	return c.doRequest(url, put, body)
}

// doRequest Performing an Http Request to the given url using the given method and (optionally) sends the body
func (c *Connection) doRequest(url string, method requestMethod, body string) (string, error) {
	if !c.checkAccessToken() {
		return "", errors.New("Access Token is not valid (maybe was not set?)")
	}
	var req *http.Request
	var err error

	if len(body) > 0 {
		req, err = http.NewRequest(string(method), url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(string(method), url, nil)
	}
	if err != nil {
		return "", err
	}
	req.Header.Set("token", c.AccessToken)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Unexpected Resonse Code for get " + url)
	}
	if res.Header.Get("Content-Length") != "0" {
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(resp), err
	}

	return "", err
}

func (c *Connection) checkAccessToken() bool {
	// ToDo do propper checks
	return c.AccessToken != ""
}
