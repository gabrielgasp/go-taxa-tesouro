package main

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func bootstrapConfig() error {
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if strings.ToLower(viper.GetString("ENV")) != "production" {
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	os.Setenv("TZ", viper.GetString("TZ"))
	viper.SetDefault("RATE_LIMIT_PER_MINUTE", 10)

	return nil
}
