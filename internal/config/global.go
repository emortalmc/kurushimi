package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	MongoDB   *MongoDBConfig
	Kafka     *KafkaConfig
	Namespace string

	PartyService *PartyServiceConfig

	LobbyFleetName string
	LobbyMatchRate time.Duration
	LobbyMatchSize int

	Port uint16

	Development bool
}

type MongoDBConfig struct {
	URI string
}

type KafkaConfig struct {
	Host string
	Port int
}

type PartyServiceConfig struct {
	ServiceHost string
	ServicePort uint16

	SettingsServiceHost string
	SettingsServicePort uint16
}

func LoadGlobalConfig() (config *Config, err error) {
	viper.SetEnvPrefix("mm")
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
