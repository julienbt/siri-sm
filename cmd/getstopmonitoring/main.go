package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/julienbt/siri-sm/internal/common/ioutils"
	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/getstopmonitoring"
	"github.com/julienbt/siri-sm/internal/siri"
)

var LOCATION_NAME = "Europe/Paris"

const MONITORING_REF string = "ILEVIA:StopPoint:BP:CAS001:LOC"

func main() {
	logger := getLogger()

	var cfg config.ConfigCheckStatus
	err := envconfig.Process("SIRISM_CHECKSTATUS", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	location, err := time.LoadLocation(LOCATION_NAME)
	if err != nil {
		logger.Fatal(err)
	}
	requestTimestamp := time.Now().In(location)
	monitoredStopVisits, htmlReqBody, htmlRespBody, err := getstopmonitoring.GetStopMonitoring(
		cfg,
		logger,
		&requestTimestamp,
		MONITORING_REF,
	)
	if len(htmlReqBody) > 0 {
		fmt.Println(htmlReqBody)
	}
	if htmlRespBody != nil {
		fmt.Println(ioutils.GetPrettyPrintOfHtmlBody(htmlRespBody))
	}
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
