package digitalstrom

// ValueChangeReceiver combines callback methods for value change subsribtion. Methods will be
// called when values have been changed
type ValueChangeReceiver interface {
	OutputChannelValueChange(deviceID string, channelIndex int, oldValue int, newValue int)
	SensorValueChange(deviceID string, sensorIndex int, oldValue float64, newValue float64)
	CircuitConsumptionValueChange(circuitID string, oldValue int, newValue int)
	CircuitMeterValueChange(circuitID string, oldValue int, newValue int)
}
