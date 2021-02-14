package digitalstrom

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Account Main communication module to communicate with API. It caches and updates Devices for
// faster communication
type Account struct {
	Connection Connection
	Structure  Structure
	Devices    map[string]Device
	Groups     map[int]Group
	Zones      map[int]Zone
	Floors     map[int]Floor
	Circuits   map[string]Circuit
	//Scenes     map[string]Scene
}

// NewAccount set connection baseURL to default and returns Account
func NewAccount() *Account {
	return &Account{
		Connection: Connection{
			BaseURL: DEFAULT_BASE_URL,
		},
		Devices:  make(map[string]Device),
		Groups:   make(map[int]Group),
		Zones:    make(map[int]Zone),
		Floors:   make(map[int]Floor),
		Circuits: make(map[string]Circuit),
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
	if err != nil {
		return err
	}
	_, err = a.RequestCircuits()
	if err != nil {
		return err
	}
	return err
}

// UpdateCircuitMeterValue is performing a getEnergyMeterValue request in order to
// receive the acutal meter value. This value wil be assign to the circuit and additionally
// returned. In case an error occured during the request, -1 will be return as well as the
// error itself.
func (a *Account) UpdateCircuitMeterValue(circuitID string) (int, error) {
	circuit, ok := a.Circuits[circuitID]
	if !ok {
		return -1, errors.New("no circuit with display id '" + circuitID + "' found")
	}

	res, err := a.Connection.doRequest(a.Connection.BaseURL+"/json/circuit/getEnergyMeterValue", get, "", map[string]string{"dsuid": circuit.DSUID})
	if err != nil {
		return -1, err
	}
	if !res.OK {
		return -1, errors.New(res.Message)
	}
	value, ok := res.Result["meterValue"].(float64)
	if !ok {
		return -1, errors.New("unexpected response - no field 'meterValue' found in response")
	}
	circuit.MeterValue = int(value)
	return int(value), nil
}

//GetSensor Returning the sensor with die index ID <sensorIndex> of device with display ID <deviceID> or nil when either
// device with the given ID couldn't be found or the sensor index is higher than the amount of sensors the device has
func (a *Account) GetSensor(deviceID string, sensorIndex int) (*Sensor, error) {
	device, ok := a.Devices[deviceID]
	if !ok {
		return nil, errors.New("no device with id '" + deviceID + "' found")
	}
	if sensorIndex > len(device.Sensors) {
		return nil, errors.New("sensorIndex " + strconv.Itoa(sensorIndex) + " out of range for device " + deviceID)
	}
	return &device.Sensors[sensorIndex], nil
}

// UpdateSensorValue is requesting the current value the given sensor has. The value will be assigned
// the the sensor.
func (a *Account) UpdateSensorValue(sensor *Sensor) error {
	params := make(map[string]string)
	params["dsid"] = sensor.device.ID
	params["sensorIndex"] = strconv.Itoa(sensor.Index)
	res, err := a.Connection.doRequest(a.Connection.BaseURL+"/json/device/getSensorValue", get, "", params)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Message)
	}

	value, ok := res.Result["sensorValue"].(float64)
	if !ok {
		return errors.New("unable to extract sensorValue from result")
	}
	sensor.Value = value

	return nil
}

// UpdateCircuitConsumptionValue is performing a getconsumption request for the circuit with the given display ID.
// The requested value will be assigned to the circuit object automatically. Additionally the requested Value will be
// return or an error (when ocurred)
func (a *Account) UpdateCircuitConsumptionValue(circuitID string) (int, error) {
	circuit, ok := a.Circuits[circuitID]
	if !ok {
		return -1, errors.New("no circuit with display id '" + circuitID + "' found")
	}

	res, err := a.Connection.doRequest(a.Connection.BaseURL+"/json/circuit/getConsumption", get, "", map[string]string{"dsuid": circuit.DSUID})
	if err != nil {
		return -1, err
	}
	if !res.OK {
		return -1, errors.New(res.Message)
	}
	value, ok := res.Result["consumption"].(float64)
	if !ok {
		return -1, errors.New("unexpected response - no field consumption found in response")
	}
	circuit.Consumption = int(value)
	return int(value), nil
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

// RequestStructure performs a getStructure request. The complete Structure will be assigned
// to the account automatically when request was sucessfull, otherwise the occured error will
// be returned.
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
	a.buildMaps()
	return &s, nil
}

// RequestCircuits performs a getCircuits request. The received circuit array
// will be assigned to the account object and additionally returned or the error
//
func (a *Account) RequestCircuits() ([]Circuit, error) {
	res, err := a.Connection.Get(a.Connection.BaseURL + "/json/apartment/getCircuits")

	if err != nil {
		return nil, err
	}

	if !res.OK {
		return nil, errors.New(res.Message)
	}

	jsonString, err := json.Marshal(res.Result["circuits"])

	if err != nil {
		return nil, err
	}

	circuits := []Circuit{}
	err = json.Unmarshal(jsonString, &circuits)
	if err != nil {
		return nil, err
	}

	for i := range circuits {
		a.Circuits[circuits[i].DisplayID] = circuits[i]
	}

	return circuits, nil
}

// buldMaps is generating maps for devices, circuits, zones, groups
// and floors for fast access. It should be called whenever a structure,
// circuit or groups are requested
func (a *Account) buildMaps() {
	for i := range a.Structure.Apart.Zones {
		zone := a.Structure.Apart.Zones[i]
		a.Zones[zone.ID] = zone
		for j := range a.Structure.Apart.Zones[i].Groups {
			group := a.Structure.Apart.Zones[i].Groups[j]
			a.Groups[group.ID] = group
		}
		for j := range a.Structure.Apart.Zones[i].Devices {
			device := a.Structure.Apart.Zones[i].Devices[j]
			a.Devices[device.DisplayID] = device
			for n := range device.Sensors {
				device.Sensors[n].Index = n
				device.Sensors[n].device = &device
			}
		}
	}
	for i := range a.Structure.Apart.Floors {
		floor := a.Structure.Apart.Floors[i]
		a.Floors[floor.ID] = floor
	}
}
