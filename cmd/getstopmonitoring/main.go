package main

import (
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/getstopmonitoring"
	"github.com/julienbt/siri-sm/internal/siri"
)

const MONITORING_REF string = "ILEVIA:StopPoint:BP:CCH002:LOC"

func main() {
	logger := getLogger()

	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	monitoredStopVisits, err := getstopmonitoring.GetStopMonitoring(cfg, logger, MONITORING_REF)
	if err != nil {
		switch e := err.(type) {
		case *siri.RemoteError:
			logger.Error(e)
		default:
			logger.Fatal(e)
		}
		return
	}
	logger.Infof("GetStopMonitoring response: %#v", monitoredStopVisits)
}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "getstopmonitoring",
		"runtime": runtime.Version(),
	})
}
