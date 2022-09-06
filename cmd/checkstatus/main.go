package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/julienbt/siri-sm/internal/checkstatus"
	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/siri"
	"github.com/julienbt/siri-sm/internal/utils"
)

var LOCATION_NAME = "Europe/Paris"

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

	checkStatusResult, htmlReqBody, htmlRespBody, err := checkstatus.CheckStatus(cfg, logger, &requestTimestamp)
	if len(htmlReqBody) > 0 {
		fmt.Println(htmlReqBody)
	}
	if htmlRespBody != nil {
		fmt.Println(utils.GetPrettyPrintOfHtmlBody(htmlRespBody))
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
	logger.Infof("CheckStatus response: %#v", checkStatusResult)
}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "checkstatus",
		"runtime": runtime.Version(),
	})
}
