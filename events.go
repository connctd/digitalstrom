package digitalstrom

import (
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type ValueChangeEvent struct {
	valueID  string
	oldValue string
	newValue string
}

var logger = stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))

//SetLogger sets a custom logger
func SetLogger(newLogger logr.Logger) {
	logger = newLogger.WithName("lib-digitalstrom")
}
