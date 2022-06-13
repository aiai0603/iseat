package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var Log *logrus.Logger

const (
	logFileName = "server.log"
	log_dir     = "/log/"
)

// Need to use logger before all, so init it at the earliest
func init() {
	logr := logrus.New()
	logr.SetReportCaller(true)
	logr.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	})
	Log = logr
}

type LoggerConfig struct {
	Debug      bool `yaml:"debug"`
	ConsoleLog bool `yaml:"consoleLog"`
	FileLog    bool `yaml:"fileLog"`
}

func initLogConf(c *LoggerConfig) {
	if c.Debug {
		Log.SetLevel(logrus.DebugLevel)
	}
	logDir, err := GetLoc(log_dir)
	if err != nil {
		Log.Fatal(err)
	}
	file := logDir + logFileName
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		if c.ConsoleLog {
			if c.FileLog {
				mw := io.MultiWriter(os.Stdout, logFile)
				Log.SetOutput(mw)
			}
		} else {
			if c.FileLog {
				Log.SetOutput(logFile)
			}
		}
	} else {
		fmt.Println("Failed to init log file settings...")
		Log.Infof("Failed to log to file, using default stderr.")
	}
}
