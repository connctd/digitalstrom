package digitalstrom

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type requestMethod string

const (
	// DefautBaseURL Will be set as default Base Url for connections
	DefautBaseURL = "https://192.168.178.178:8080" // TODO: should be exchanged by default dS web API address
)

// RequestError structure to give a proper error that could occur during
// http requests
type RequestError struct {
	Err        string
	StatusCode int
	Request    http.Request
}

// Error function of structure RequestError
func (e *RequestError) Error() string {
	return fmt.Sprintf("connection error - status code %d", e.StatusCode)
}

// RequestResult represents all successful http request results received
// from the digitalStrom server. It always contains the OK value as idication
// whether the requested data could be delivered or not. When OK is true, RequestResult
// contains Result, the requested data as json (map[string]interface{}). When
// OK is false, RequestResult contains Message, the problem desription the digitalStrom had
// for not delivering the requested data.
type RequestResult struct {
	OK      bool                   `json:"ok"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result"`
}

const (
	get  requestMethod = "GET"
	put  requestMethod = "PUT"
	post requestMethod = "POST"
)

// Connection Holds access credentials and URL information for a specific account.
// Contains routines for GET, PUSH, PUT and DELETE
type Connection struct {
	SessionToken     string
	ApplicationToken string
	BaseURL          string
	HTTPClient       http.Client
}

// Get Performs a GET request and returns the response body as string
func (c *Connection) Get(url string) (*RequestResult, error) {
	return c.Request(url, get, "", nil)
}

// Post Performs a Post Request with the given content and returns the response body as string
func (c *Connection) Post(url string, body string) (*RequestResult, error) {
	return c.Request(url, post, body, nil)
}

// Put Performs a Put Request with the given content and returns the response body as string
func (c *Connection) Put(url string, body string) (*RequestResult, error) {
	return c.Request(url, put, body, nil)
}

// Request is performing an Http-Request. In case it receives an HTTP-Error 403, an application Login will be performed and the
// request will be repeated (only one time).
func (c *Connection) Request(url string, method requestMethod, body string, params map[string]string) (*RequestResult, error) {
	res, err := c.doRequest(url, method, body, params)
	if err != nil {
		if reqErr, ok := err.(*RequestError); ok {
			if reqErr.StatusCode == 403 {
				e := c.applicationLogin()
				if e != nil {
					return res, e
				}
				return c.doRequest(url, method, body, params)
			}
		}
		return res, err
	}
	return res, nil
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

// doRequest is performing an http request. It is not recommended to use method "Connection.Request". In case the session token has been expired
// doRequest is giving back the error and is not trying to login automatically
func (c *Connection) doRequest(url string, method requestMethod, body string, params map[string]string) (*RequestResult, error) {

	logger.Info("performing http-request: " + url)

	req, err := c.generateHTTPRequest(url, method, body, params)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		logger.Error(err, "unable to perform http request")
		return nil, err
	}

	defer res.Body.Close()

	// TODO: generate a better error message according to StatusCode
	if res.StatusCode != http.StatusOK {
		err := &RequestError{"Unexpected Resonse for " + url, res.StatusCode, *req}
		logger.Error(err, "unable to perform http request")
		return nil, err
	}
	if res.Header.Get("Content-Length") != "0" {
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Error(err, "unable to read response body (content length: "+res.Header.Get("Content-Length")+")")
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

	logger.Info("registering new application with name " + applicationName)

	// request an ApplicationToken
	res, err := c.doRequest(c.BaseURL+"/json/system/requestApplicationToken", get, "", map[string]string{"applicationName": applicationName})
	if err != nil {
		logger.Error(err, "registration has been aborted")
		return "", err
	}
	if !res.OK {
		e := errors.New(res.Message)
		logger.Error(e, "registration has been aborted")
		return "", e
	}

	applicationToken := res.Result["applicationToken"].(string)

	logger.Info("got application token for '" + applicationName)

	// performing a login in order to generate a temporary session token
	// HINT during the http request generation, parameters will be URL-encoded. This leads to a problem with the password
	// when it contains symbols that will change during the encoding problem. The digitalstrom server is not performing an
	// URL decoding to restore the password
	// This is a workaround, that might fail with passwords which are not url compatible
	logger.Info("request session token with user credentials")
	res, err = c.doRequest(c.BaseURL+"/json/system/login?user="+username+"&password="+password, get, "", nil)
	if err != nil {
		logger.Error(err, "registration has been aborted")
		return "", err
	}
	if !res.OK {
		e := errors.New(res.Message)
		logger.Error(e, "registration has been aborted")
		return "", e
	}

	sessionToken := res.Result["token"].(string)
	logger.Info("got session token, trying to enable the application token")
	// use the session token to enable the application token. Future logins wont need user credentials anymore, the application token will be used to
	// perform an application login
	res, err = c.doRequest(c.BaseURL+"/json/system/enableToken", get, "", map[string]string{"applicationToken": applicationToken, "token": sessionToken})
	if err != nil {
		logger.Error(err, "registration has been aborted")
		return "", err
	}
	logger.Info("application '" + applicationName + " has been registered successfully")
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
