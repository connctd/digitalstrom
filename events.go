package digitalstrom

import (
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

// ValueChangeReceiver combines callback methods for value change subsribtion. Methods will be
// called when values have been changed
type ValueChangeReceiver interface {
	OnOutputChannelValueChange(deviceID string, applicationType ApplicationType, oldValue int, newValue int)
	OnSensorValueChange(deviceID string, sensorIndex int, oldValue float64, newValue float64)
	OnCircuitConsumptionValueChange(circuitID string, oldValue int, newValue int)
	OnCircuitMeterValueChange(circuitID string, oldValue int, newValue int)
}

var logger = stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))

//SetLogger sets a custom logger
func SetLogger(newLogger logr.Logger) {
	logger = newLogger.WithName("lib-digitalstrom")
}
