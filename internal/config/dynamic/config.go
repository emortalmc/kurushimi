package dynamic

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

type Config struct {
	Redis     RedisConfig    `yaml:"redis"`
	RabbitMq  RabbitMQConfig `yaml:"rabbitmq"`
	Namespace string         `yaml:"namespace"`
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

func ParseConfig() (config Config, err error) {
	viper.SetEnvPrefix("kurushimi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	zap.S().Infow("Parsed config", "config", config)
	return
}
