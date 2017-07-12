package logging

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test")
	logger.SetLevel(log.DebugLevel)
	logger.Debugln("DEBUG")
	logger.Warnln("WARN")
	logger.Infoln("INFO")
	logger.Errorln("ERROR")
	logger.WithField("test", "test").Printf("test")

}

func TestWirteOutput(t *testing.T) {
	logger := NewLogger("file")
	var file *os.File
	_, err := os.Stat("logging.log")
	if err != nil {
		// no file, create
		file, _ = os.Create("logging.log")
	} else {
		file, _ = os.OpenFile("logging.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	}
	SetLoggerOut(file)
	logger.Debugln("DEBUG")
	logger.Warnln("WARN")
	logger.Infoln("INFO")
	logger.Errorln("ERROR")
	logger.WithField("test", "test").Printf("test")
}
