package digitalstrom

import "fmt"

//	"encoding/json"
//	"errors"

// Account Main communication module to communicate with API. It caches and updates Devices for
// faster communication
type Account struct {
	Connection Connection
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
func (a *Account) SetToken(token string) {
	a.Connection.AccessToken = token
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

//
func (a *Account) RequestAll() (string, error) {
	fmt.Println("requesting all")
	body, err := a.Connection.Get(a.Connection.BaseURL + "/json/apartment/getStructure")
	if err != nil {
		return "", err
	}
	return body, err
}
