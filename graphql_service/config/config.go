package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/probes"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"

	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "GraphQL microservice config path")
}

type Config struct {
	ServiceName string             `mapstructure:"serviceName"`
	Resources   Resources          `mapstructure:"resources"`
	Logger      *logger.Config     `mapstructure:"logger"`
	KafkaTopics KafkaTopics        `mapstructure:"kafkaTopics"`
	Http        Http               `mapstructure:"http"`
	Grpc        Grpc               `mapstructure:"grpc"`
	Kafka       *kafka.Config      `mapstructure:"kafka"`
	Probes      probes.Config      `mapstructure:"probes"`
	OTL         *tracing.OTLConfig `mapstructure:"otl"`
}

type Resources struct {
	GRAPHQL_QUERY string `mapstructure:"graphql_query"`
	API_ACCOUNTS  string `mapstructure:"api_accounts"`
}

type Http struct {
	Port                string   `mapstructure:"port"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath"`
	GraphQLPath         string   `mapstructure:"graphqlPath"`
	DebugHeaders        bool     `mapstructure:"debugHeaders"`
	HttpClientDebug     bool     `mapstructure:"httpClientDebug"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

type Grpc struct {
	CoreServicePort string `mapstructure:"coreServicePort"`
}

type KafkaTopics struct {
	WebhookExample kafka.TopicConfig `mapstructure:"webhookExample"`
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
			configPath = fmt.Sprintf("%s/graphql_service/config/config.yaml", getwd)
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

	utils.CheckKafkaGroup(&cfg.Kafka.GroupID, constants.ShortMicroserviceName)
	utils.CheckOTLName(&cfg.OTL.ServiceName, constants.GRAPHQL, constants.ShortMicroserviceName)

	utils.CheckEnvStr(&cfg.Http.Port, constants.HttpPort)
	utils.CheckEnvStr(&cfg.OTL.Endpoint, constants.OTLEndpoint)
	utils.CheckEnvStr(&cfg.Grpc.CoreServicePort, constants.CoreServicePort)
	utils.CheckEnvStr(&cfg.Resources.GRAPHQL_QUERY, constants.GraphQLQuery)
	utils.CheckEnvArrStr(&cfg.Kafka.Brokers, constants.KafkaBrokers)

	return cfg, nil
}
