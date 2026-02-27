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
	flag.StringVar(&configPath, "config", "", "API Gateway microservice config path")
}

type Config struct {
	ServiceName  string             `mapstructure:"serviceName"`
	Logger       *logger.Config     `mapstructure:"logger"`
	KafkaTopics  KafkaTopics        `mapstructure:"kafkaTopics"`
	Http         Http               `mapstructure:"http"`
	Grpc         Grpc               `mapstructure:"grpc"`
	Kafka        *kafka.Config      `mapstructure:"kafka"`
	Probes       probes.Config      `mapstructure:"probes"`
	Resources    Resources          `mapstructure:"resources"`
	OTL          *tracing.OTLConfig `mapstructure:"otl"`
	DevelopeMode bool               `mapstructure:"developeMode"`
}

type Http struct {
	Port                string   `mapstructure:"port"`
	Development         bool     `mapstructure:"development"`
	Title               string   `mapstructure:"title"`
	BasePath            string   `mapstructure:"basePath"`
	V1Path              string   `mapstructure:"v1Path"`
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

type Resources struct {
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
			configPath = fmt.Sprintf("%s/api_gateway_service/config/config.yaml", getwd)
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
	utils.CheckHttpTitle(&cfg.Http.Title, constants.ShortMicroserviceName)
	utils.CheckOTLName(&cfg.OTL.ServiceName, constants.GATEWAY, constants.ShortMicroserviceName)

	utils.CheckEnvStr(&cfg.Http.Port, constants.HttpPort)
	utils.CheckEnvStr(&cfg.OTL.Endpoint, constants.OTLEndpoint)
	utils.CheckEnvStr(&cfg.Grpc.CoreServicePort, constants.CoreServicePort)
	utils.CheckEnvArrStr(&cfg.Kafka.Brokers, constants.KafkaBrokers)
	utils.CheckEnvStr(&cfg.Http.BasePath, constants.HttpBasePath)

	utils.CheckEnvBool(&cfg.DevelopeMode, constants.DevelopeMode)

	if cfg.DevelopeMode {
		cfg.KafkaTopics.WebhookExample.TopicName += "Dev"
	}

	return cfg, nil
}
