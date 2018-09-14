package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var (
	config string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "",
		"config file")
}

func initConfig() {
	viper.SetEnvPrefix("BKIE")
	viper.AutomaticEnv()

	if config == "" {
		return
	}
	viper.SetConfigFile(config)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading configuration: %v", err)
	}
}
