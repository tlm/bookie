package cmd

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logJson bool = false
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&logJson, "log-json", false,
		"enables json logging")
}

func setupLogging() {
	if logJson {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:    false,
			DisableSorting:   true,
			DisableTimestamp: false,
			FullTimestamp:    true,
			TimestampFormat:  time.RFC3339,
		})
	}
}
