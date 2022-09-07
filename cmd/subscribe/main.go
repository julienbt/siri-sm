package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/julienbt/siri-sm/internal/common/ioutils"
	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/subscribe"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var LOCATION_NAME = "Europe/Paris"

func main() {

	logger := getLogger()

	var cfg config.ConfigSubscribe
	err := envconfig.Process("SIRISM_SUBSCRIBE", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	location, err := time.LoadLocation(LOCATION_NAME)
	if err != nil {
		logger.Fatal(err)
	}
	requestTimestamp := time.Now().In(location)

	subscribeResp, htmlReqBody, htmlRespBody, err := subscribe.Subscribe(cfg, logger, &requestTimestamp)
	if len(htmlReqBody) > 0 {
		fmt.Println(htmlReqBody)
	}
	if htmlRespBody != nil {
		fmt.Println(ioutils.GetPrettyPrintOfHtmlBody(htmlRespBody))
	}
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Subscribe response: %#v", subscribeResp)
}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "subscribe",
		"runtime": runtime.Version(),
	})
}
