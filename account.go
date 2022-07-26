package digitalstrom

import (
	"encoding/json"
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

// Default polling setup values
const (
	defaultSensorPollingInterval  = 30
	defaultCircuitPollingInterval = 5
	defaultChannelPollingInterval = 30
	defaultMaxSimultanousPolls    = 10
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

	// updating
	PollingSetup      PollingSetup
	pollingHelpers    pollingHelpers
	quitTickerChannel chan bool

	// events
	eventHelpers eventHelpers
}

// PollingSetup defines update settings for automated value polling
type PollingSetup struct {
	DefaultCircuitsPollingInterval int `json:"default_circuit_polling_interval"`
	DefaultSensorsPollingInterval  int `json:"default_sensors_polling_interval"`
	DefaultChannelsPollingInterval int `json:"default_channels_polling_interval"`
	MaxParallelPolls               int `json:"max_parallel_polls"`
}

type pollingHelpers struct {
	parallelPollCount int
	pollIntervalMap   map[string]int
	lastPollMap       map[string]time.Time
	activePollingMap  map[string]time.Time
	mapMutex          *sync.Mutex
}

// eventHelpers are used to handle async request responses for updating value changes. Requesting values
// from devices over the api takes several seconds. That is why multiple requests are made simultaneously
// with the help of gophers.
// this structure helps to lock maps when async events will be handled
type eventHelpers struct {
	valueChangeReceiver map[string]ValueChangeReceiver
	mapMutex            *sync.Mutex
}

var logger = stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))

//SetLogger sets a custom logger
func SetLogger(newLogger logr.Logger) {
	logger = newLogger.WithName("lib-digitalstrom")
}

// NewAccount sets connection baseURL to default, generates maps and returns
// empty Account instance
func NewAccount() *Account {
	return &Account{
		Connection: Connection{
			BaseURL: defautBaseURL,
		},
		Devices:           make(map[string]Device),
		Groups:            make(map[int]Group),
		Zones:             make(map[int]Zone),
		Floors:            make(map[int]Floor),
		Circuits:          make(map[string]Circuit),
		quitTickerChannel: make(chan bool),
		PollingSetup: PollingSetup{
			DefaultCircuitsPollingInterval: defaultCircuitPollingInterval,
			DefaultChannelsPollingInterval: defaultChannelPollingInterval,
			DefaultSensorsPollingInterval:  defaultSensorPollingInterval,
			MaxParallelPolls:               defaultMaxSimultanousPolls,
		},
		pollingHelpers: pollingHelpers{
			parallelPollCount: 0,
			pollIntervalMap:   nil,
			activePollingMap:  make(map[string]time.Time),
			lastPollMap:       make(map[string]time.Time),
			mapMutex:          &sync.Mutex{},
		},
		eventHelpers: eventHelpers{
			valueChangeReceiver: make(map[string]ValueChangeReceiver),
			mapMutex:            &sync.Mutex{},
		},
	}

}

// ApplicationLogin uses the assigned applicationToken to generate a session token. The timeout depends on server settings,
// default is 180 seconds. This timeout will be automatically reset by every performed request.
func (a *Account) ApplicationLogin() error {
	return a.Connection.applicationLogin()
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

//GetOutputChannel Returning the output channel with die index ID <channelIndex> of device with display ID <deviceID> or nil when either
// device with the given ID couldn't be found or the OutputChannel index is higher than the amount of channels the device has
func (a *Account) GetOutputChannel(deviceID string, channelIndex int) (*OutputChannel, error) {
	device, ok := a.Devices[deviceID]
	if !ok {
		return nil, errors.New("no device with id '" + deviceID + "' found")
	}
	if channelIndex >= len(device.OutputChannels) {
		return nil, errors.New("channelIndex " + strconv.Itoa(channelIndex) + " out of range for device " + deviceID)
	}
	return &device.OutputChannels[channelIndex], nil
}

// Init of the Account. ApplicationLogin will be performed and complete structure requested. ApplicationToken has to be set in advance.
func (a *Account) Init() error {
	logger.Info("account initialization")
	logger.Info("performing application login")
	err := a.ApplicationLogin()
	if err != nil {
		logger.Error(err, "initialisation has been aborted")
		return err
	}
	logger.Info("requesting complete structure")
	s, err := a.RequestStructure()
	if err != nil {
		logger.Error(err, "initialisation has been aborted")
		return err
	}
	a.setStructure(*s)
	logger.Info("requesting circuits")
	circuits, err := a.RequestCircuits()
	if err != nil {
		logger.Error(err, "initialisation has been aborted")
		return err
	}
	// fill the circuit map for fast access
	for i := range circuits {
		a.Circuits[circuits[i].DisplayID] = circuits[i]
	}
	logger.Info("account successfully initialized")
	return nil
}

// RegisterApplication an application with the given applicitonName. Performs a request to generate an application token. A second request requires the
// Username and Password in order to generate a temporary session token. A third request enables the application token to login without
// further user credentials (applicationLogin). Returns the application token or an error. The application token will not be assigned automatically.
// Thus, in order to use the generated application token, it has to be set afterwards (Account.SetApplicationToken).
func (a *Account) RegisterApplication(applicationName string, username string, password string) (string, error) {
	return a.Connection.register(username, password, applicationName)
}

// RequestCircuits performs a getCircuits request. The received circuit array
// has to be assigned to the account separately
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

	// there we are, return everything
	return circuits, nil
}

// RequestStructure performs a getStructure request and returns it or an error that might have been occured.
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
	err = json.Unmarshal(jsonString, &s)

	// assign the structure to our account
	if err != nil {
		return nil, err
	}
	// return the shit
	return &s, err
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
	err = json.Unmarshal(jsonString, &s)

	return &s, err
}

// ResetPollingIntervals will remove all intervals for sensors,
// circuits and output channels. When StartPolling() is called, default
// polling intervals for sensors, circuits and output channels will not
// be set to default. To set default polling intervals again, call
// SetDefaultPollingIntervas().
func (a *Account) ResetPollingIntervals() {
	a.pollingHelpers.mapMutex.Lock()
	a.pollingHelpers.pollIntervalMap = make(map[string]int)
	a.pollingHelpers.mapMutex.Unlock()
}

// StartPolling starts the update routine. It calls the internal prepareUpdates function.
// When no update intervals are given in advance, a complete list of update intervals will
// be generated automatically (including all sensors, output channesl and circuits) by using
// the related default intervals.
func (a *Account) StartPolling() {
	a.preparePolling()
	ticker := *time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for id, interval := range a.pollingHelpers.pollIntervalMap {
					if a.isPollingIntervalReached(id, interval) {
						if a.pollingHelpers.parallelPollCount < a.PollingSetup.MaxParallelPolls {
							a.pollingHelpers.mapMutex.Lock()
							_, ok := a.pollingHelpers.activePollingMap[id]
							a.pollingHelpers.mapMutex.Unlock()
							if !ok {
								a.pollingHelpers.parallelPollCount++
								go a.performPolling(id)
							}
						}
					}
				}
			case <-a.quitTickerChannel:
				ticker.Stop()
				return
			}
		}
	}()
}

// SetApplicationToken that will be used for ApplicationLogin
func (a *Account) SetApplicationToken(token string) {
	a.Connection.ApplicationToken = token
}

// SetDefaultPollingIntervals is setting for all sensors, channels and circuits
// the corresponding default interval. Intervals that were set manually before, will be overwritten.
func (a *Account) SetDefaultPollingIntervals() {
	a.ResetPollingIntervals()
	for devID, dev := range a.Devices {
		for i := range dev.Sensors {
			id := "sensor." + devID + "." + strconv.Itoa(i)
			a.pollingHelpers.pollIntervalMap[id] = a.PollingSetup.DefaultSensorsPollingInterval
		}
		for i := range dev.OutputChannels {
			id := "channel." + devID + "." + strconv.Itoa(i)
			a.pollingHelpers.pollIntervalMap[id] = a.PollingSetup.DefaultChannelsPollingInterval
		}

	}
	for circuitID := range a.Circuits {
		a.pollingHelpers.pollIntervalMap["circuit."+circuitID] = a.PollingSetup.DefaultCircuitsPollingInterval
	}
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

// SetSessionToken for manually setting the token. Be aware of a timout for each session token. It is recommended to perform
// an ApplicationLogin using the ApplicationToken. This will update the session token automatically.
func (a *Account) SetSessionToken(token string) {
	a.Connection.SessionToken = token
}

// SetPollingInterval sets the automatic polling interval for the element identified
// by given id. The id is a combination of eighter "sensor.<deviceID>.<sensor index>,
// channel.<deviceID>.<ChannelType> or circuit.<circuitID>. When setting an polling interval,
// only those elements will be polled, that were added. To set default polling intervals for
// all elements, call SetDefaultPollingIntervals()
func (a *Account) SetPollingInterval(id string, interval int) error {
	if a.pollingHelpers.pollIntervalMap == nil {
		a.pollingHelpers.pollIntervalMap = make(map[string]int)
	}

	s := strings.Split(id, ".")
	if len(s) < 2 {
		return errors.New(id + " is not a valid identifier")
	}

	// ToDo: do better id test (sensor existing, channel existing, circuit existing)

	a.pollingHelpers.pollIntervalMap[id] = interval
	return nil
}

// SetURL sets the BaseUrl. This method should only be called when another URL should be used than the default one
func (a *Account) SetURL(url string) {
	a.Connection.BaseURL = url
}

// StopPolling stops the autonomous updater. Values will not requested until
// updater will be started again
func (a *Account) StopPolling() {
	a.quitTickerChannel <- true
}

// SubsribeValueChangeReceiver adds a pointer to a ValueChangeReceiver with the given id. Everytime a value of a sensor, circuit or
// output channel changes, the corresponding callback function of the receiver will be called
func (a *Account) SubsribeValueChangeReceiver(id string, receiver ValueChangeReceiver) {
	if receiver == nil {
		return
	}
	a.eventHelpers.mapMutex.Lock()
	logger.Info(fmt.Sprintf("receiver %s added to value change receiver map (mapsize = %d)", id, len(a.eventHelpers.valueChangeReceiver)))
	a.eventHelpers.valueChangeReceiver[id] = receiver
	a.eventHelpers.mapMutex.Unlock()
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

// UnsubscribeValueChangeReceiver removes an event channel.
func (a *Account) UnsubscribeValueChangeReceiver(id string) {
	a.eventHelpers.mapMutex.Lock()
	delete(a.eventHelpers.valueChangeReceiver, id)
	a.eventHelpers.mapMutex.Unlock()
}

// PollCircuitMeterValue is performing a getEnergyMeterValue request in order to
// receive the acutal meter value. This value wil be assign to the circuit and additionally
// returned. In case an error occured during the request, -1 will be return as well as the
// error itself.
func (a *Account) PollCircuitMeterValue(circuitID string) (int, error) {
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
		oldValue := circuit.MeterValue
		circuit.MeterValue = newValue
		a.Circuits[circuitID] = circuit
		a.dispatchMeterValueChange(circuitID, oldValue, newValue)
	}

	return newValue, nil
}

// PollCircuitConsumptionValue is performing a getconsumption request for the circuit with the given display ID.
// The requested value will be assigned to the circuit object automatically. Additionally the requested Value will be
// return or an error (when ocurred)
func (a *Account) PollCircuitConsumptionValue(circuitID string) (int, error) {

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
		oldValue := circuit.Consumption
		circuit.Consumption = newValue
		a.Circuits[circuitID] = circuit
		a.dispatchConsumptionValueChange(circuitID, oldValue, newValue)
	}
	return newValue, nil
}

// PollOnValues ...
func (a *Account) PollOnValues() error {
	s, err := a.RequestStructure()
	if err != nil {
		return err
	}
	for i := range s.Apartment.Zones {
		zone := s.Apartment.Zones[i]
		for n := range zone.Devices {
			device := zone.Devices[n]
			if (device.IsPresent) && (device.IsValid) {

				ad, ok := a.Devices[device.DisplayID]
				if ok {
					if ad.On != device.On {
						//oldval := ad.On

						fmt.Printf(" device %s.on = %t -> %t\n", device.DisplayID, ad.On, device.On)
						ad.On = device.On
						a.Devices[device.DisplayID] = ad
						// TODO: dispatch on value change a.dis
					}
				}
			}
		}
	}
	return nil
}

// PollOnValue ...
func (a *Account) PollOnValue(device *Device) (bool, error) {

	return false, errors.New("not implemented yet")
}

// PollSensorValue is requesting the current value the given sensor has. The value will be assigned
// the the sensor.
func (a *Account) PollSensorValue(sensor *Sensor) (float64, error) {
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
	//fmt.Printf("sensor "+sensor.device.DisplayID+".%d = %f\r\n", sensor.Index, value)
	if !ok {
		return 0, errors.New("unable to extract sensorValue from request result")
	}

	if sensor.Value != value {
		oldValue := sensor.Value
		sensor.Value = value
		a.dispatchSensorValueChange(sensor.device.DisplayID, sensor.Index, oldValue, value)
	}

	return value, nil
}

func (a *Account) PollChannelValue(channel *OutputChannel) (int, error) {
	params := make(map[string]string)
	params["dsid"] = channel.device.ID
	params["offset"] = strconv.Itoa(channel.ChannelIndex)
	res, err := a.Connection.Request(a.Connection.BaseURL+"/json/device/getOutputValue", get, "", params)
	if err != nil {
		return 0, err
	}
	if !res.OK {
		return 0, errors.New(res.Message)
	}

	value, ok := res.Result["value"].(float64)
	if !ok {
		return 0, errors.New("unable to extract channel output value from request result")
	}

	int_value := int(value)
	if channel.Value != int_value {
		//oldValue := channel.Value
		channel.Value = int_value
		// TODO dispatch change a.dispatchOutputChannelValueChange(channel.device.ID, channel.ChannelType, oldValue,value)
	}

	return int_value, nil
}

// SetStructure assigns a structure to the account, maps will be build and cross references assigned.
// It is highy recommended to use this function in order to assign a structure rather than assign them directly.
func (a *Account) setStructure(structure Structure) {
	a.Structure = structure
	a.Structure.assignCrossReferences()
	a.buildMaps()
}

// buldMaps is generating maps for devices, circuits, zones, groups
// and floors for fast access. It should be called whenever a structure,
// circuit or groups are requested
func (a *Account) buildMaps() {
	for i := range a.Structure.Apartment.Zones {
		zone := a.Structure.Apartment.Zones[i]
		a.Zones[zone.ID] = zone
		for j := range a.Structure.Apartment.Zones[i].Groups {
			group := a.Structure.Apartment.Zones[i].Groups[j]
			a.Groups[group.ID] = group
		}
		for j := range a.Structure.Apartment.Zones[i].Devices {
			device := a.Structure.Apartment.Zones[i].Devices[j]
			a.Devices[device.DisplayID] = device
		}
	}
	for i := range a.Structure.Apartment.Floors {
		floor := a.Structure.Apartment.Floors[i]
		a.Floors[floor.ID] = floor
	}
}

func (a *Account) preparePolling() {
	// create a new map to temporarily store the last updates for each value
	if a.pollingHelpers.pollIntervalMap == nil {
		a.SetDefaultPollingIntervals()
	}
	a.pollingHelpers.lastPollMap = make(map[string]time.Time)
	for key := range a.pollingHelpers.pollIntervalMap {
		a.pollingHelpers.lastPollMap[key] = time.Now()
	}
	a.pollingHelpers.parallelPollCount = 0
}

// isPollingIntervalReached checks whether the value with the given ID
// needs to be updated.
func (a *Account) isPollingIntervalReached(id string, interval int) bool {
	t, ok := a.pollingHelpers.lastPollMap[id]
	if !ok {
		return true
	}
	return time.Now().Sub(t).Seconds() > float64(interval)
}

func (a *Account) setPollingTimeStamp(id string) {
	a.pollingHelpers.parallelPollCount--
	a.pollingHelpers.mapMutex.Lock()
	a.pollingHelpers.lastPollMap[id] = time.Now()
	// remove id from polling list
	// use mutex to prevent concurrent map writes

	delete(a.pollingHelpers.activePollingMap, id)
	a.pollingHelpers.mapMutex.Unlock()
}

func (a *Account) dispatchConsumptionValueChange(circuitID string, oldValue int, newValue int) {
	a.eventHelpers.mapMutex.Lock()

	for id, receiver := range a.eventHelpers.valueChangeReceiver {
		logger.Info(fmt.Sprintf("calling %s.OnConsumptionValueChange for ciruit %s (%d -> %d))", id, circuitID, oldValue, newValue))
		go receiver.CircuitConsumptionValueChange(circuitID, oldValue, newValue)
	}
	a.eventHelpers.mapMutex.Unlock()
}

func (a *Account) dispatchMeterValueChange(circuitID string, oldValue int, newValue int) {
	a.eventHelpers.mapMutex.Lock()

	for id, receiver := range a.eventHelpers.valueChangeReceiver {
		logger.Info(fmt.Sprintf("calling %s.OnMeterValueChange for ciruit %s (%d -> %d))", id, circuitID, oldValue, newValue))
		go receiver.CircuitMeterValueChange(circuitID, oldValue, newValue)
	}
	a.eventHelpers.mapMutex.Unlock()
}

func (a *Account) dispatchOutputChannelValueChange(deviceID string, at ApplicationType, oldValue int, newValue int) {
	a.eventHelpers.mapMutex.Lock()
	for id, receiver := range a.eventHelpers.valueChangeReceiver {
		logger.Info(fmt.Sprintf("calling %s.OnOutputChannelValueChange for channel %s.%s (%d -> %d))", id, deviceID, at.GetName(), oldValue, newValue))
		go receiver.OutputChannelValueChange(deviceID, at, oldValue, newValue)
	}
	a.eventHelpers.mapMutex.Unlock()
}

func (a *Account) dispatchSensorValueChange(deviceID string, sensorIndex int, oldValue float64, newValue float64) {
	a.eventHelpers.mapMutex.Lock()

	for id, receiver := range a.eventHelpers.valueChangeReceiver {
		logger.Info(fmt.Sprintf("calling %s.OnSensorValueChange for sensor %s.%d (%f -> %f))", id, deviceID, sensorIndex, oldValue, newValue))
		go receiver.SensorValueChange(deviceID, sensorIndex, oldValue, newValue)
	}
	a.eventHelpers.mapMutex.Unlock()
}

func (a *Account) performPolling(id string) {

	// independed from update result, set the current timestamp to reset the interval
	defer a.setPollingTimeStamp(id)
	// remember that value with this id will be polled now
	a.pollingHelpers.mapMutex.Lock()
	a.pollingHelpers.activePollingMap[id] = time.Now()
	a.pollingHelpers.mapMutex.Unlock()

	logger.Info(fmt.Sprintf("updating %s (%d/%d)", id, a.pollingHelpers.parallelPollCount, a.PollingSetup.MaxParallelPolls))

	// ids are separated by '.'
	s := strings.Split(id, ".")

	switch s[0] {
	case "circuit":
		if len(s) != 2 {
			logger.Info(fmt.Sprintf("WARNING: %s is not a valid curcuitID ", id))
			return
		}
		a.PollCircuitConsumptionValue(s[1])
		a.PollCircuitMeterValue(s[1])

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
		a.PollSensorValue(sensor)

	case "channel":
		if len(s) != 3 {
			return
		}
		number, err := strconv.Atoi(s[2])
		if err != nil {
			return
		}
		channel, err := a.GetOutputChannel(s[1], number)
		if err != nil {
			return
		}
		a.PollChannelValue(channel)

		return
	default:
		// place error logging for invalid id over here

	}

}
