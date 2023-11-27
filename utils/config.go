package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	Port     string `mapstructure:"PORT"`
	LogLevel int    `mapstructure:"LOGLEVEL"`
}

func GetConfig() Config {
	once.Do(func() {
		err := factoryConfig()
		if err != nil {
			fmt.Println(err)
		}
	})
	return config
}

func factoryConfig() error {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)

	var loadConfig Config
	if err := v.ReadInConfig(); err != nil {
		loadConfig = defaultConfig()
	} else {
		err = v.Unmarshal(&loadConfig)
		if err != nil {
			return err
		}
	}
	config = loadConfig
	return nil
}

func defaultConfig() Config {
	return Config{
		Port:     "3000",
		LogLevel: 0,
	}
}
