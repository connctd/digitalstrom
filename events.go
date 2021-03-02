package digitalstrom

// ValueChangeReceiver combines callback methods for value change subsribtion. Methods will be
// called when values have been changed
type ValueChangeReceiver interface {
	OnOutputChannelValueChange(deviceID string, applicationType ApplicationType, oldValue int, newValue int)
	OnSensorValueChange(deviceID string, sensorIndex int, oldValue float64, newValue float64)
	OnCircuitConsumptionValueChange(circuitID string, oldValue int, newValue int)
	OnCircuitMeterValueChange(circuitID string, oldValue int, newValue int)
}
