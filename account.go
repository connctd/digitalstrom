package digitalstrom

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Default update intervals
const (
	DefaultSensorUpdateInterval  = 120
	DefaultCircuitUpdateInterval = 3
	DefaultChannelUpdateInterval = 30
)

// Account Main communication module to communicate with API. It caches and updates Devices for
// faster communication
type Account struct {
	Connection Connection
	structure  Structure
	Devices    map[string]Device
	Groups     map[int]Group
	Zones      map[int]Zone
	Floors     map[int]Floor
	Circuits   map[string]Circuit
	//Scenes     map[string]Scene

	// updating
	ticker      time.Ticker
	quitTicker  chan bool
	updateSetup updateSetup
}

type updateSetup struct {
	maxSimultanousUpdates int
	currentUpdates        int
	intervals             map[string]int
	lastUpdate            map[string]time.Time
}

// NewAccount sets connection baseURL to default, generates maps and returns
// empty Account instance
func NewAccount() *Account {

	return &Account{
		Connection: Connection{
			BaseURL: DEFAULT_BASE_URL,
		},
		Devices:    make(map[string]Device),
		Groups:     make(map[int]Group),
		Zones:      make(map[int]Zone),
		Floors:     make(map[int]Floor),
		Circuits:   make(map[string]Circuit),
		quitTicker: make(chan bool),
	}

}

// GetStructure returns the structure of the account.
func (a *Account) GetStructure() *Structure {
	return &a.structure
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

	return nil
}

// UpdateOnValue ...
func (a *Account) UpdateOnValue(device *Device) (bool, error) {
	return false, errors.New("not implemented yet")
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

	res, err := a.Connection.Request(a.Connection.BaseURL+"/json/circuit/getEnergyMeterValue", get, "", map[string]string{"dsuid": circuit.DSUID})
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

	newValue := int(value)

	if newValue != circuit.MeterValue {
		//	oldValue := circuit.MeterValue
		circuit.MeterValue = newValue
		a.Circuits[circuitID] = circuit
		//	dispatchValueChange(channel, conumption, oldvalue, newvalue)
		//a.OnCircuitMeterValueUpdate <- CircuitMeterValueUpdateEvent{CircuitID: circuitID, OldValue: oldValue, NewValue: newValue}
	}

	return newValue, nil
}

//GetSensor Returning the sensor with die index ID <sensorIndex> of device with display ID <deviceID> or nil when either
// device with the given ID couldn't be found or the sensor index is higher than the amount of sensors the device has
func (a *Account) GetSensor(deviceID string, sensorIndex int) (*Sensor, error) {
	device, ok := a.Devices[deviceID]
	if !ok {
		return nil, errors.New("no device with id '" + deviceID + "' found")
	}
	if sensorIndex >= len(device.Sensors) {
		return nil, errors.New("sensorIndex " + strconv.Itoa(sensorIndex) + " out of range for device " + deviceID)
	}
	return &device.Sensors[sensorIndex], nil
}

// SetOutputChannelValue sets the value for the given OutputChannel. Returns error
func (a *Account) SetOutputChannelValue(channel *OutputChannel, value string) error {
	params := make(map[string]string)
	params["dsid"] = channel.device.ID
	params["channelvalues"] = string(channel.ChannelType) + "=" + value

	res, err := a.Connection.Request(a.Connection.BaseURL+"/json/device/setOutputChannelValue", get, "", params)
	if err != nil {
		return err
	}

	if !res.OK {
		return errors.New(res.Message)
	}
	return nil

}

// UpdateSensorValue is requesting the current value the given sensor has. The value will be assigned
// the the sensor.
func (a *Account) UpdateSensorValue(sensor *Sensor) (float64, error) {
	params := make(map[string]string)
	params["dsid"] = sensor.device.ID
	params["sensorIndex"] = strconv.Itoa(sensor.Index)
	res, err := a.Connection.Request(a.Connection.BaseURL+"/json/device/getSensorValue", get, "", params)
	if err != nil {
		return 0, err
	}
	if !res.OK {
		return 0, errors.New(res.Message)
	}

	value, ok := res.Result["sensorValue"].(float64)
	fmt.Printf("sensor "+sensor.device.DisplayID+".%d = %f\r\n", sensor.Index, value)
	if !ok {
		return 0, errors.New("unable to extract sensorValue from request result")
	}
	sensor.Value = value

	return value, nil
}

// UpdateCircuitConsumptionValue is performing a getconsumption request for the circuit with the given display ID.
// The requested value will be assigned to the circuit object automatically. Additionally the requested Value will be
// return or an error (when ocurred)
func (a *Account) UpdateCircuitConsumptionValue(circuitID string) (int, error) {

	circuit, ok := a.Circuits[circuitID]
	if !ok {
		return -1, errors.New("no circuit with display id '" + circuitID + "' found")
	}

	res, err := a.Connection.Request(a.Connection.BaseURL+"/json/circuit/getConsumption", get, "", map[string]string{"dsuid": circuit.DSUID})
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

	newValue := int(value)

	if newValue != circuit.Consumption {
		//oldValue := circuit.Consumption
		circuit.Consumption = newValue
		a.Circuits[circuitID] = circuit

	}
	return newValue, nil
}

// ApplicationLogin uses the assigned applicationToken to generate a session token. The timeout depends on server settings,
// default is 180 seconds. This timeout will be automatically reset by every performed request.
func (a *Account) ApplicationLogin() error {
	return a.Connection.applicationLogin()
}

// RegisterApplication an application with the given applicitonName. Performs a request to generate an application token. A second request requires the
// Username and Password in order to generate a temporary session token. A third request enables the application token to login without
// further user credentials (applicationLogin). Returns the application token or an error. The application token will not be assigned automatically.
// Thus, in order to use the generated application token, it has to be set afterwards (Account.SetApplicationToken).
func (a *Account) RegisterApplication(applicationName string, username string, password string) (string, error) {
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

	// assign the structure to our account
	a.setStructure(s)

	// return the shit
	return &s, nil
}

// RequestSystemInfo performs a get system/version request and returns the result or
// the error that has been occurred
func (a *Account) RequestSystemInfo() (*System, error) {
	res, err := a.Connection.Get(a.Connection.BaseURL + "/json/system/version")

	if err != nil {
		return nil, err
	}

	if !res.OK {
		return nil, errors.New(res.Message)
	}
	jsonString, _ := json.Marshal(res.Result)

	s := System{}
	json.Unmarshal(jsonString, &s)

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
	// get result as map[string]interface{}
	jsonString, err := json.Marshal(res.Result["circuits"])

	if err != nil {
		return nil, err
	}

	circuits := []Circuit{}
	// let json.Unmarshal do the job of mapping to Circuit
	err = json.Unmarshal(jsonString, &circuits)
	if err != nil {
		return nil, err
	}

	// fill the circuit map for fast access
	for i := range circuits {
		a.Circuits[circuits[i].DisplayID] = circuits[i]
	}
	// there we are, return everything
	return circuits, nil
}

// TurnOn sends eithe a turnOn or turnOff request for the given 'device', depending on value of paramter 'on'
func (a *Account) TurnOn(device *Device, on bool) error {

	var url = ""
	if on {
		url = "/json/device/turnOn"
	} else {
		url = "/json/device/turnOff"
	}

	res, err := a.Connection.Request(a.Connection.BaseURL+url, get, "", map[string]string{"dsid": device.ID})
	if err != nil {
		return err
	}

	if !res.OK {
		return errors.New(res.Message)
	}

	return nil
}

// SetStructure assigns a structure to the account, maps will be build and cross references assigned.
// It is highy recommended to use this function in order to assign a structure rather than assign them directly.
func (a *Account) setStructure(structure Structure) {
	a.structure = structure
	a.structure.assignCrossReferences()
	a.buildMaps()
}

// buldMaps is generating maps for devices, circuits, zones, groups
// and floors for fast access. It should be called whenever a structure,
// circuit or groups are requested
func (a *Account) buildMaps() {
	for i := range a.structure.Apartment.Zones {
		zone := a.structure.Apartment.Zones[i]
		a.Zones[zone.ID] = zone
		for j := range a.structure.Apartment.Zones[i].Groups {
			group := a.structure.Apartment.Zones[i].Groups[j]
			a.Groups[group.ID] = group
		}
		for j := range a.structure.Apartment.Zones[i].Devices {
			device := a.structure.Apartment.Zones[i].Devices[j]
			a.Devices[device.DisplayID] = device
		}
	}
	for i := range a.structure.Apartment.Floors {
		floor := a.structure.Apartment.Floors[i]
		a.Floors[floor.ID] = floor
	}
}

func (a *Account) RegisterEventChannel(channel chan ValueChangeEvent) {
	// register listener for consumption change events

}

func (a *Account) prepareUpdates() {
	// create a new map to temporarily store the last updates for each value
	if a.updateSetup.intervals == nil {
		a.generateDefaultUpdateSetups()
	}
	a.updateSetup.lastUpdate = make(map[string]time.Time)
	for key := range a.updateSetup.intervals {
		a.updateSetup.lastUpdate[key] = time.Now()
	}
	a.updateSetup.currentUpdates = 0
	a.updateSetup.maxSimultanousUpdates = 5
}

func (a *Account) generateDefaultUpdateSetups() {
	a.updateSetup.intervals = make(map[string]int)
	for devID, dev := range a.Devices {
		for i := range dev.Sensors {
			id := "sensor." + devID + "." + strconv.Itoa(i)
			a.updateSetup.intervals[id] = DefaultSensorUpdateInterval
		}
		for i := range dev.OutputChannels {
			id := "channel." + devID + "." + dev.OutputChannels[i].ChannelID
			a.updateSetup.intervals[id] = DefaultChannelUpdateInterval
		}
	}
	for circuitID := range a.Circuits {
		a.updateSetup.intervals["circuit."+circuitID] = DefaultCircuitUpdateInterval
	}
}

// isUpdateIntervalReached checks whether the value with the given ID
// needs to be updated.
func (a *Account) isUpdateIntervalReached(id string, interval int) bool {
	t, ok := a.updateSetup.lastUpdate[id]
	if !ok {
		return true
	}
	return time.Now().Sub(t).Seconds() > float64(interval)
}

func (a *Account) RunUpdates() {
	a.prepareUpdates()
	a.ticker = *time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-a.ticker.C:
				for id, interval := range a.updateSetup.intervals {
					if a.isUpdateIntervalReached(id, interval) {
						if a.updateSetup.currentUpdates < a.updateSetup.maxSimultanousUpdates {
							a.updateSetup.currentUpdates++
							go a.performUpdate(id)
						}
					}
				}
			case <-a.quitTicker:
				a.ticker.Stop()
				return
			}
		}
	}()
}

func (a *Account) decreaseUpdateCounter() {
	a.updateSetup.currentUpdates--
}

func (a *Account) performUpdate(id string) {

	defer a.decreaseUpdateCounter()

	fmt.Printf("\r\n%s (%d)", id, a.updateSetup.currentUpdates)

	s := strings.Split(id, ".")

	if len(s) < 2 {
		return
	}

	switch s[0] {
	case "circuit":
		if len(s) != 2 {
			return
		}
		a.UpdateCircuitConsumptionValue(s[1])
		a.UpdateCircuitConsumptionValue(s[1])
		a.updateSetup.lastUpdate[id] = time.Now()
	case "sensor":
		if len(s) != 3 {
			return
		}
		number, err := strconv.Atoi(s[2])
		if err != nil {
			return
		}
		sensor, err := a.GetSensor(s[1], number)
		if err != nil {
			return
		}
		a.UpdateSensorValue(sensor)
		a.updateSetup.lastUpdate[id] = time.Now()
	case "channel":
		// not implemented yet
		a.updateSetup.lastUpdate[id] = time.Now()
		return
	default:
		// place error logging for invalid id over here
		a.updateSetup.lastUpdate[id] = time.Now()
	}

}

func (a *Account) StopUpdates() {
	a.quitTicker <- true
}
