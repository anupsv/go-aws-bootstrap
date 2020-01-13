package config

import (
	"github.com/spf13/viper"
)

type Init struct {
	ConfigLocation string
	ConfigFileName string
}

type Config struct {
	S3Info  S3Info
	SqsInfo SqsInfo
}

func (init Init) readConfig() {

	viper.SetConfigName(init.ConfigFileName)
	viper.AddConfigPath(init.ConfigLocation)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}
