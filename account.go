package digitalstrom

import (
	//"encode/json"
	"encoding/json"
	"errors"
)

// Account Main communication module to communicate with API. It caches and updates Devices for
// faster communication
type Account struct {
	Connection Connection
	Structure  Structure
}

// NewAccount set connection baseURL to default and returns Account
func NewAccount() *Account {
	return &Account{
		Connection: Connection{
			BaseURL: DEFAULT_BASE_URL,
		},
	}
}

// SetSessionToken for manually setting the token. Be aware of a timout for each session token. It is recommended to perform
// an ApplicationLogin using the ApplicationToken. This will update the session token automatically.
func (a *Account) SetSessionToken(token string) {
	a.Connection.SessionToken = token
}

// SetApplicationToken that will be used for ApplicationLogin
func (a *Account) SetApplicationToken(token string) {
	a.Connection.ApplicationToken = token
}

// SetURL sets the BaseUrl. This method should only be called when another URL should be used than the default one
func (a *Account) SetURL(url string) {
	a.Connection.BaseURL = url
}

// Init of the Account. ApplicationLogin will be performed and complete structure requested. ApplicationToken has to be set in advance.
func (a *Account) Init() error {
	err := a.ApplicationLogin()
	if err != nil {
		return err
	}
	_, err = a.RequestStructure()
	return err
}

// ApplicationLogin uses the assigned applicationToken to generate a session token. The timeout depends on server settings,
// default is 180 seconds. This timeout will be automatically reset by every performed request.
func (a *Account) ApplicationLogin() error {
	return a.Connection.applicationLogin()
}

// Register an application with the given applicitonName. Performs a request to generate an application token. A second request requires the
// Username and Password in order to generate a temporary session token. A third request enables the application token to login without
// further user credentials (applicationLogin). Returns the application token or an error. The application token will not be assigned automatically.
// Thus, in order to use the generated application token, it has to be set afterwards (Account.SetApplicationToken).
func (a *Account) Register(applicationName string, username string, password string) (string, error) {
	return a.Connection.register(username, password, applicationName)
}

//
func (a *Account) RequestStructure() (*Structure, error) {

	res, err := a.Connection.Get(a.Connection.BaseURL + "/json/apartment/getStructure")

	if err != nil {
		return nil, err
	}

	if !res.OK {
		return nil, errors.New(res.Message)
	}
	jsonString, _ := json.Marshal(res.Result)

	s := Structure{}
	json.Unmarshal(jsonString, &s)

	a.Structure = s

	return &s, nil
}
