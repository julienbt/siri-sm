package main

import (
	"fmt"
	"runtime"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/julienbt/siri-sm/internal/checkstatus"
	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/siri"
	"github.com/julienbt/siri-sm/internal/utils"
)

func main() {
	logger := getLogger()

	var cfg config.ConfigCheckStatus
	err := envconfig.Process("SIRISM_CHECKSTATUS", &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	checkStatusResult, htmlBody, err := checkstatus.CheckStatus(cfg, logger)
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
	if htmlBody != nil {
		fmt.Println(utils.GetPrettyPrintOfHtmlBody(htmlBody))
	}
}

func getLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"app":     "checkstatus",
		"runtime": runtime.Version(),
	})
}
