package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/postgres"
	"github.com/testovoleg/5s-microservice-template/pkg/probes"
	"github.com/testovoleg/5s-microservice-template/pkg/redis"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Reader microservice config path")
}

type Config struct {
	ServiceName      string              `mapstructure:"serviceName"`
	Logger           *logger.Config      `mapstructure:"logger"`
	KafkaTopics      KafkaTopics         `mapstructure:"kafkaTopics"`
	GRPC             GRPC                `mapstructure:"grpc"`
	Postgresql       *postgres.Config    `mapstructure:"postgres"`
	Kafka            *kafkaClient.Config `mapstructure:"kafka"`
	Redis            *redis.Config       `mapstructure:"redis"`
	MongoCollections MongoCollections    `mapstructure:"mongoCollections"`
	Probes           probes.Config       `mapstructure:"probes"`
	ServiceSettings  ServiceSettings     `mapstructure:"serviceSettings"`
	OTL              *tracing.OTLConfig  `mapstructure:"otl"`
	Resources        Resources           `mapstructure:"resources"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type MongoCollections struct {
	Products string `mapstructure:"products"`
}

type KafkaTopics struct {
	ProductCreated kafkaClient.TopicConfig `mapstructure:"productCreated"`
	ProductUpdated kafkaClient.TopicConfig `mapstructure:"productUpdated"`
	ProductDeleted kafkaClient.TopicConfig `mapstructure:"productDeleted"`
}

type ServiceSettings struct {
	RedisProductPrefixKey string `mapstructure:"redisProductPrefixKey"`
}

type Resources struct {
	REDOCLY_JSON string `mapstructure:"redocly_json"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/core_service/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	utils.CheckEnvStr(&cfg.GRPC.Port, constants.GrpcPort)
	utils.CheckEnvStr(&cfg.OTL.Endpoint, constants.OTLEndpoint)
	utils.CheckEnvArrStr(&cfg.Kafka.Brokers, constants.KafkaBrokers)
	utils.CheckEnvStr(&cfg.Resources.REDOCLY_JSON, constants.RedoclyJSON)

	return cfg, nil
}
