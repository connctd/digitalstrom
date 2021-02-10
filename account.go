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

// SetToken Sets the access token for future use.
func (a *Account) SetSessionToken(token string) {
	a.Connection.SessionToken = token
}

func (a *Account) SetApplicationToken(token string) {
	a.Connection.ApplicationToken = token
}

// SetURL sets the BaseUrl. This method should only be called when another URL should be used than the default one
func (a *Account) SetURL(url string) {
	a.Connection.BaseURL = url
}

// Init Inits the Account with Devices from the API
func (a *Account) Init() error {
	// everything seems to be fine, do not return an error
	return nil
}

func (a *Account) Login() error {
	if !a.Connection.checkApplicationToken() {
		return errors.New("Application Token is not set. Unable to perform application login.")
	}
	params := map[string]string{"loginToken": a.Connection.ApplicationToken}
	res, err := a.Connection.doRequest(a.Connection.BaseURL+"/json/system/loginApplication", get, "", params)
	if err != nil {
		return err
	}

	if res.OK {
		a.Connection.SessionToken = res.Result["token"].(string)
		return nil
	}
	return errors.New(res.Message)
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
