package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Redis     RedisConfig    `yaml:"redis"`
	RabbitMq  RabbitMQConfig `yaml:"rabbitmq"`
	Namespace string         `yaml:"namespace"`

	Development bool `yaml:"development"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func LoadGlobalConfig() (config *Config, err error) {
	viper.SetEnvPrefix("kurushimi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return
}
