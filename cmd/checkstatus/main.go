package main

import (
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/julienbt/siri-sm/internal/checkstatus"
	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/siri"
)

func main() {
	logger := getLogger()

	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	checkStatusResult, err := checkstatus.CheckStatus(cfg, logger)
	if err != nil {
		switch e := err.(type) {
		case *siri.RemoteError:
			logger.Error(e)
		default:
			logger.Fatal(e)
		}
		return
	}
	logger.Infof("checkstatus response: %#v", checkStatusResult)
}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "checkstatus",
		"runtime": runtime.Version(),
	})
}
