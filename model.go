package digitalstrom

import (
	"encoding/json"
	//"errors"
	//"strconv"
	//"strings"
)

type Structure struct {
	Apart Apartment `json:"apartment"`
}

type Apartment struct {
	Zones  []Zone  `json:"zones"`
	Floors []Floor `json:"floors"`
}

type Zone struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	IsPresent bool     `json:"isPresent"`
	FloorID   int      `json:"floorId"`
	Devices   []Device `json:"devices"`
	Groups    []Group  `json:"groups"`
}

type Floor struct {
	ID    int    `json:"id"`
	Order int    `json:"order"`
	Name  string `json:"name"`
	Zones []int  `json:"zones"`
}

type Group struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Color           int      `json:"color"`
	ApplicationType int      `json:"applicationType"`
	IsPresent       bool     `json:"isPresent"`
	IsValid         bool     `json:"isValid"`
	Devices         []string `json:"devices"`
}

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
	Consumption                 int
	MeterValue                  int
}

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
	ButtonUsage           string `json:"buttonUsage"`
	MeterDSID             string `json:"meterDSID"`
	MeterDSUID            string `json:"meterDSUID"`
	MeterName             string `json:"meterName"`
	BusID                 int    `json:"busID"`
	ZoneID                int    `json:"zoneID"`
	IsPresent             bool   `json:"isPresent"`
	IsValid               bool   `json:"isValid"`
	LastDiscovered        string `json:"lastDiscovered"`
	FirstSeen             string `json:"firstSeen"`
	InactiveSince         string `json:"inactiveSince"`
	On                    bool   `json:"on"`
	Locked                bool   `json:"locked"`
	ConfigurationLocked   bool   `json:"configurationLocked"`
	IgnoreOperationLock   bool   `json:"ignoreOperationLock"`
	OutputMode            int    `json:"outputMode"`
	ButtonID              int    `json:"buttonID"`
	ButtonActiveGroup     int    `json:"buttonActiveGroup"`
	ButtonGroupMemberShip int    `json:"buttonMemberShip"`
	ButtonInputMode       int    `json:"buttonInputMode"`
	ButtonInputIndex      int    `json:"buttonInputIndex"`
	ButtonInputCount      int    `json:"buttonInputCount"`
	AKMInputProperty      string `json:"AKMInputProperty"`
	//Groups
	BinaryInputCount int `json:"binaryInputCount"`
	//BinaryInputs
	SensorInputCount int             `json:"sonsorInputCount"`
	Sensors          []Sensor        `json:"sensors"`
	SensorDataValid  bool            `json:"sensorDataValid"`
	OutputChannels   []OutputChannel `json:"outputChannels"`
	//PairedDevices
}

type OutputChannel struct {
	ChannelID    string `json:"channelID"`
	ChannelType  string `json:"channelType"`
	ChannelIndex int    `json:"channelIndex"`
	ChannelName  string `json:"channelName"`
	Value        int
}

type Sensor struct {
	Type   int     `json:"type"`
	Valid  bool    `json:"valid"`
	Value  float64 `json:"value"`
	Index  int
	device *Device
}

type OutputChannelType string

const (
	OCT_brightness               OutputChannelType = "brightness"
	OCT_hue                      OutputChannelType = "hue"
	OCT_saturation               OutputChannelType = "saturation"
	OCT_colortemp                OutputChannelType = "colortemp"
	OCT_x                        OutputChannelType = "x"
	OCT_y                        OutputChannelType = "y"
	OCT_shadePositionOutside     OutputChannelType = "shadePositionOutside"
	OCT_shadePositionIndoor      OutputChannelType = "shadePositionIndoor"
	OCT_shadeOpeningAngleOutside OutputChannelType = "shadeOpeningAngleOutside"
	OCT_shadeOpeningAngleInside  OutputChannelType = "shadeOpeningAngleInside"
	OCT_transparency             OutputChannelType = "transparency"
	OCT_airFlowIntensity         OutputChannelType = "airFlowIntensity"
	OCT_airFlowDirection         OutputChannelType = "airFlowDirection"
	OCT_airFlapPosition          OutputChannelType = "airFlapPosition"
	OCT_airLouverPosition        OutputChannelType = "airLouverPosition"
	OCT_heatingPower             OutputChannelType = "heatingPower"
	OCT_coolingCapacity          OutputChannelType = "coolingCapacity"
	OCT_audioVolume              OutputChannelType = "audioVolume"
	OCT_powerState               OutputChannelType = "powerState"
)

// Generates the Apartment structure out of a given json string
func GenerateApartment(j string) (*Apartment, error) {
	var apartement Apartment

	err := json.Unmarshal([]byte(j), &apartement)
	if err != nil {
		return nil, err
	}

	return &apartement, nil
}
