package digitalstrom

type ChannelValueChangeEvent struct {
	DeviceID     string
	ChannelIndex int
	OldValue     int
	NewValue     int
}

type SensorValueChangeEvent struct {
	DeviceId    string
	SensorIndex int
	OldValue    float64
	NewValue    float64
}

type CircuitConsumptionValueChangeEvent struct {
	CircuitID string
	OldValue  int
	NewValue  int
}

type CircuitMeterValueChangeEvent struct {
	CircuitID string
	OldValue  int
	NewValue  int
}

type OnStateValueChangeEvent struct {
	DeviceId string
	OldValue bool
	NewValue bool
}

type ZoneTemperatureControlChangeEvent struct {
	ZoneId int
}

type BinaryInputStateChangeEvent struct {
	DeviceId string
	InputId  int
	OldValue int
	NewValue int
}
