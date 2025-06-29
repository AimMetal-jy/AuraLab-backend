package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Log is the global logger.
var Log = logrus.New()

// InitLogger initializes the logger.
func InitLogger() {
	// Log as JSON instead of the default ASCII formatter.
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	Log.SetLevel(logrus.InfoLevel)
}