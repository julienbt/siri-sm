package main

import (
	"fmt"
	"runtime"
	"time"

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

	subscribeResp, httpReqBody, httpRespBody, err := subscribe.Subscribe(cfg, logger, location)
	fmt.Println(httpReqBody)
	if err != nil {
		logger.Fatal(err)
	}
	_ = subscribeResp
	_ = httpRespBody

}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "subscribe",
		"runtime": runtime.Version(),
	})
}
