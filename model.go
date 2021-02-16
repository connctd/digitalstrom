package digitalstrom

import (
	"encoding/json"
	//"errors"
	//"strconv"
	//"strings"
)

// Structure represents the digitalSTROM structure of
// an installation.
type Structure struct {
	Apart Apartment `json:"apartment"`
}

// Apartment is the ogical instance of a digitalSTROM instalation. This includes
// all rooms and any device.
type Apartment struct {
	Zones  []Zone  `json:"zones"`
	Floors []Floor `json:"floors"`
}

// Zone is a logical representation of one room, hall or
// other partial structural works of a building.
type Zone struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	IsPresent bool     `json:"isPresent"`
	FloorID   int      `json:"floorId"`
	Devices   []Device `json:"devices"`
	Groups    []Group  `json:"groups"`
}

// Floor ...
type Floor struct {
	ID    int    `json:"id"`
	Order int    `json:"order"`
	Name  string `json:"name"`
	Zones []int  `json:"zones"`
}

// Group is the representation of an pplication type
type Group struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Color           int             `json:"color"`
	ApplicationType ApplicationType `json:"applicationType"`
	IsPresent       bool            `json:"isPresent"`
	IsValid         bool            `json:"isValid"`
	Devices         []string        `json:"devices"`
}

// Circuit is the physical connection between the circuit breaker and digitalSTROM-Meter on one ende and each electic deice connected on this line on the other hand
type Circuit struct {
	Name                        string `json:"name"`
	DSID                        string `json:"dsid"`
	DSUID                       string `json:"dSUID"`
	DisplayID                   string `json:"DisplayID"`
	HwVersion                   int    `json:"hwVersion"`
	HwVersionString             string `json:"hvVersionString"`
	SwVersion                   string `json:"swVersion"`
	ArmSwVersion                int    `json:"armSwVersion"`
	DspSwVersion                int    `json:"dspSwVersion"`
	IsUpToDate                  bool   `json:"isUpToDate"`
	APIVersion                  int    `json:"apiVersion"`
	Authorized                  bool   `json:"authorized"`
	HwName                      string `json:"HwName"`
	IsPresent                   bool   `json:"isPresent"`
	IsValid                     bool   `json:"isValid"`
	BusMemberType               int    `json:"busMemberType"`
	HasDevices                  bool   `json:"hasDevices"`
	HasMetering                 bool   `json:"hasMetering"`
	HasBlinking                 bool   `json:"hasBlinking"`
	VdcConfigURL                string `json:"VdcConfigURL"`
	VdcModelUID                 string `json:"VdcModelUID"`
	VdcHardwareGUID             string `json:"vdcHardwareGuid"`
	VdcHardwareModelGUID        string `json:"vdcHardwareModelGuid"`
	VdcImplementationID         string `json:"vdcImplementationId"`
	VdcVendorGUID               string `json:"vdcVendorGuid"`
	VdcOemGUID                  string `json:"vdcOemGuid"`
	IgnoreActionsFromNewDevices bool   `json:"ignoreActionsFromNewDevices"`
	Consumption                 int    // not part of json, values have to be requested separately
	MeterValue                  int    // not part of json, values have to be requested separately
}

// Device  ...
type Device struct {
	ID                  string `json:"id"`
	DisplayID           string `json:"DisplayID"`
	UUID                string `json:"dSUID"`
	Gtin                string `json:"GTIN"`
	Name                string `json:"name"`
	DsUIDIndex          int    `json:"dSUIDIndex"`
	FunctionID          int    `json:"functionID"`
	ProductRevision     int    `json:"productRevision"`
	ProductID           int    `json:"productID"`
	HwInfo              string `json:"hwInfo"`
	OemStatus           string `json:"OemStatus"`
	OemEanNumber        string `json:"OemEanNumber"`
	OemSerialNumber     int    `json:"OemSerialNumber"`
	OemPartNumber       int    `json:"OemPartNumber"`
	OemProductInfoState string `json:"OemProductInfoState"`
	OemProductURL       string `json:"OemProductURL"`
	OemInternetState    string `json:"OemInternetState"`
	OemIsIndependent    bool   `json:"OemIsIndependent"`
	//ModelFeatures			ModelFeature `json:"modelFeatures"`
	IsVdcDevice bool `json:"isVdcDevice"`
	//SupportedBasicScenes []BasicScene
	ButtonUsage           string          `json:"buttonUsage"`
	MeterDSID             string          `json:"meterDSID"`
	MeterDSUID            string          `json:"meterDSUID"`
	MeterName             string          `json:"meterName"`
	BusID                 int             `json:"busID"`
	ZoneID                int             `json:"zoneID"`
	IsPresent             bool            `json:"isPresent"`
	IsValid               bool            `json:"isValid"`
	LastDiscovered        string          `json:"lastDiscovered"`
	FirstSeen             string          `json:"firstSeen"`
	InactiveSince         string          `json:"inactiveSince"`
	On                    bool            `json:"on"`
	Locked                bool            `json:"locked"`
	ConfigurationLocked   bool            `json:"configurationLocked"`
	IgnoreOperationLock   bool            `json:"ignoreOperationLock"`
	OutputMode            int             `json:"outputMode"`
	ButtonID              int             `json:"buttonID"`
	ButtonActiveGroup     int             `json:"buttonActiveGroup"`
	ButtonGroupMemberShip int             `json:"buttonMemberShip"`
	ButtonInputMode       int             `json:"buttonInputMode"`
	ButtonInputIndex      int             `json:"buttonInputIndex"`
	ButtonInputCount      int             `json:"buttonInputCount"`
	AKMInputProperty      string          `json:"AKMInputProperty"`
	BinaryInputCount      int             `json:"binaryInputCount"`
	BinaryInputs          []BinaryInput   `json:"binaryInputs"`
	SensorInputCount      int             `json:"sonsorInputCount"`
	Sensors               []Sensor        `json:"sensors"`
	SensorDataValid       bool            `json:"sensorDataValid"`
	OutputChannels        []OutputChannel `json:"outputChannels"`
	PairedDevices         []string        `json:"pairedDevices"`
	Groups                []int           `json:"groups"`
}

// BinaryInput ...
type BinaryInput struct {
	TargetGroup int             `json:"targetGroup"`
	InputType   BinaryInputType `json:"inputType"`
	InputID     int             `json:"inputId"`
	State       int             `json:"state"`
}

// OutputChannel ....
type OutputChannel struct {
	ChannelID    string            `json:"channelID"`
	ChannelType  OutputChannelType `json:"channelType"`
	ChannelIndex int               `json:"channelIndex"`
	ChannelName  string            `json:"channelName"`
	Value        int
}

// Sensor ...
type Sensor struct {
	Type   int     `json:"type"`
	Valid  bool    `json:"valid"`
	Value  float64 `json:"value"`
	Index  int
	device *Device
}

// OutputChannelType ...
type OutputChannelType string

// ApplicationType ....
type ApplicationType int

// ApplicationColor ...
type ApplicationColor string

// BinaryInputType ...
type BinaryInputType int

// Output Channel Types (OTC)
const (
	OCTbrightness               OutputChannelType = "brightness"
	OCThue                      OutputChannelType = "hue"
	OCTsaturation               OutputChannelType = "saturation"
	OCTcolortemp                OutputChannelType = "colortemp"
	OCTx                        OutputChannelType = "x"
	OCTy                        OutputChannelType = "y"
	OCTshadePositionOutside     OutputChannelType = "shadePositionOutside"
	OCTshadePositionIndoor      OutputChannelType = "shadePositionIndoor"
	OCTshadeOpeningAngleOutside OutputChannelType = "shadeOpeningAngleOutside"
	OCTshadeOpeningAngleInside  OutputChannelType = "shadeOpeningAngleInside"
	OCTtransparency             OutputChannelType = "transparency"
	OCTairFlowIntensity         OutputChannelType = "airFlowIntensity"
	OCTairFlowDirection         OutputChannelType = "airFlowDirection"
	OCTairFlapPosition          OutputChannelType = "airFlapPosition"
	OCTairLouverPosition        OutputChannelType = "airLouverPosition"
	OCTheatingPower             OutputChannelType = "heatingPower"
	OCTcoolingCapacity          OutputChannelType = "coolingCapacity"
	OCTaudioVolume              OutputChannelType = "audioVolume"
	OCTpowerState               OutputChannelType = "powerState"
)

// Application Types
const (
	ATlights               ApplicationType = 1
	ATblinds               ApplicationType = 2
	ATheating              ApplicationType = 3
	ATaudio                ApplicationType = 4
	ATvideo                ApplicationType = 5
	ATcooling              ApplicationType = 9
	ATventilation          ApplicationType = 10
	ATwindow               ApplicationType = 11
	ATrecirculation        ApplicationType = 12
	ATtemperatureControl   ApplicationType = 48
	ATapartmentVentilation ApplicationType = 64
	ATsingleDevice         ApplicationType = -1 // no id defined
	ATsecurity             ApplicationType = -2 // no id defined
	ATaccess               ApplicationType = -3 // no id defined
)

// Available Colors
const (
	ACwhite   ApplicationColor = "white"
	ACblack   ApplicationColor = "black"
	ACgreen   ApplicationColor = "green"
	ACgray    ApplicationColor = "gray"
	ACblue    ApplicationColor = "blue"
	ACcyan    ApplicationColor = "cyan"
	ACmagenta ApplicationColor = "magenta"
	ACred     ApplicationColor = "red"
)

// Binary Input Types
const (
	BITgeneric               BinaryInputType = 0
	BITpresence              BinaryInputType = 1
	BITbrightness            BinaryInputType = 2
	BITpresenceInDarkness    BinaryInputType = 3
	BITtwilight              BinaryInputType = 4
	BITmotion                BinaryInputType = 5
	BITmotionInDarkness      BinaryInputType = 6
	BITsmoke                 BinaryInputType = 7
	BITwindStrenghAboveLimit BinaryInputType = 8
	BITrain                  BinaryInputType = 9
	BITsunRadiation          BinaryInputType = 10
	BITtemperatureBelowLimit BinaryInputType = 11
	BITBatteryStatusIsLow    BinaryInputType = 12
	BITwindowIsOpen          BinaryInputType = 13
	BITdoorIsOpen            BinaryInputType = 14
	BITwindowIsTilted        BinaryInputType = 15
	BITgarageDoorIsOpen      BinaryInputType = 16
	BITsunProtection         BinaryInputType = 17
	BITfrost                 BinaryInputType = 18
	BITheatingSystemEnabled  BinaryInputType = 19
	BITchangeOverSignal      BinaryInputType = 20
	BITinitialization        BinaryInputType = 21
	BITmalfunction           BinaryInputType = 22
	BITservice               BinaryInputType = 23
)

// GetID returns the id of the binary input type
func (b BinaryInputType) GetID() int {
	return int(b)
}

// GetName returns the Name of the binary input type
func (b BinaryInputType) GetName() string {
	names := [...]string{
		"Generic",
		"Presence",
		"Brightness",
		"Presence",
		"Twilight",
		"Motion",
		"Motion in darkness",
		"Smoke",
		"Wind strength above limit",
		"Rain",
		"Sun radiation",
		"Temperature below limit",
		"Battery status is low",
		"Window is open",
		"Door is open",
		"Window is tilted",
		"Garage door is open",
		"Sun protection",
		"Frost",
		"Heating system enabled",
		"Change-over signal",
		"Initialization",
		"Malfunction",
		"Service"}

	if len(names) < b.GetID() {
		return "unknown binary input type"
	}

	return names[b.GetID()]
}

// GetID returns the identifier for the application type
func (at ApplicationType) GetID() int {
	return int(at)
}

// GetName returns the name of the application type
func (at ApplicationType) GetName() string {
	switch at {
	case ATlights:
		return "lights"
	case ATblinds:
		return "blinds"
	case ATheating:
		return "heating"
	case ATaudio:
		return "audio"
	case ATvideo:
		return "video"
	case ATcooling:
		return "cooling"
	case ATventilation:
		return "ventilation"
	case ATwindow:
		return "window"
	case ATrecirculation:
		return "recirculation"
	case ATtemperatureControl:
		return "temperature control"
	case ATapartmentVentilation:
		return "apartment ventilation"
	case ATsingleDevice:
		return "single device"
	case ATsecurity:
		return "security"
	case ATaccess:
		return "access"
	}
	return "unknown application type"

}

// GenerateApartment takes a json string and generates and returns an instance of structure Apartment
// or the error that may have occured.
func GenerateApartment(j string) (*Apartment, error) {
	var apartement Apartment

	err := json.Unmarshal([]byte(j), &apartement)
	if err != nil {
		return nil, err
	}

	return &apartement, nil
}
