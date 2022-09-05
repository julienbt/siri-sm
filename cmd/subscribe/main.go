package main

import (
	"runtime"

	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/subscribe"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := getLogger()

	var cfg config.ConfigSubscribe
	err := envconfig.Process("SIRISM_SUBSCRIBE", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	subscribeResp, err := subscribe.Subscribe(cfg, logger)
	_ = err
	_ = subscribeResp

}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "subscribe",
		"runtime": runtime.Version(),
	})
}
