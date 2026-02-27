package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/kafka"
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
	ServiceName     string              `mapstructure:"serviceName"`
	Logger          *logger.Config      `mapstructure:"logger"`
	Keycloak        Keycloak            `mapstructure:"keycloak"`
	KafkaTopics     KafkaTopics         `mapstructure:"kafkaTopics"`
	GRPC            GRPC                `mapstructure:"grpc"`
	Postgresql      *postgres.Config    `mapstructure:"postgres"`
	Kafka           *kafkaClient.Config `mapstructure:"kafka"`
	Redis           *redis.Config       `mapstructure:"redis"`
	Probes          probes.Config       `mapstructure:"probes"`
	ServiceSettings ServiceSettings     `mapstructure:"serviceSettings"`
	OTL             *tracing.OTLConfig  `mapstructure:"otl"`
	API             API                 `mapstructure:"5sApi"`
	DevelopeMode    bool                `mapstructure:"developeMode"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type KafkaTopics struct {
	WebhookExample kafka.TopicConfig `mapstructure:"webhookExample"`
}

type ServiceSettings struct {
	RedisMicroservicePrefixKey string `mapstructure:"redisMicroservicePrefixKey"`
}

type Keycloak struct {
	Host         string `mapstructure:"host"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
}

type API struct {
	AdminApiUrl        string `mapstructure:"adminApiUrl"`
	AuthApiUrl         string `mapstructure:"authApiUrl"`
	ExportApiUrl       string `mapstructure:"exportApiUrl"`
	StorageApiUrl      string `mapstructure:"storageApiUrl"`
	ApiUsername        string `mapstructure:"apiUsername"`
	ApiPassword        string `mapstructure:"apiPassword"`
	PropertyKeyAPIData string `mapstructure:"propertyKeyAPIData"`
	StorageBucket      string `mapstructure:"storageBucket"`
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

	utils.CheckKafkaGroup(&cfg.Kafka.GroupID, constants.ShortMicroserviceName)
	utils.CheckOTLName(&cfg.OTL.ServiceName, constants.GATEWAY, constants.ShortMicroserviceName)

	utils.CheckEnvStr(&cfg.GRPC.Port, constants.GrpcPort)
	utils.CheckEnvStr(&cfg.Postgresql.Host, constants.PostgresqlHost)
	utils.CheckEnvStr(&cfg.Postgresql.Port, constants.PostgresqlPort)
	utils.CheckEnvStr(&cfg.Postgresql.User, constants.PostgresqlUser)
	utils.CheckEnvStr(&cfg.Postgresql.Password, constants.PostgresqlPassword)
	utils.CheckEnvStr(&cfg.Postgresql.DBName, constants.PostgresqlDatabase)
	utils.CheckEnvStr(&cfg.Redis.Addr, constants.RedisAddr)
	utils.CheckEnvStr(&cfg.Redis.Password, constants.RedisPassword)
	utils.CheckEnvInt(&cfg.Redis.DB, constants.RedisDB)
	utils.CheckEnvInt(&cfg.Redis.PoolSize, constants.RedisPoolSize)
	utils.CheckEnvStr(&cfg.OTL.Endpoint, constants.OTLEndpoint)
	utils.CheckEnvStr(&cfg.Keycloak.Host, constants.KeycloakHost)
	utils.CheckEnvStr(&cfg.Keycloak.Realm, constants.KeycloakRealm)
	utils.CheckEnvStr(&cfg.Keycloak.ClientID, constants.KeycloakClientId)
	utils.CheckEnvStr(&cfg.Keycloak.ClientSecret, constants.KeycloakClientSecret)
	utils.CheckEnvArrStr(&cfg.Kafka.Brokers, constants.KafkaBrokers)
	utils.CheckEnvStr(&cfg.API.AdminApiUrl, constants.AdminAPIURL)
	utils.CheckEnvStr(&cfg.API.AuthApiUrl, constants.AuthAPIURL)
	utils.CheckEnvStr(&cfg.API.ApiUsername, constants.APIUsername)
	utils.CheckEnvStr(&cfg.API.ApiPassword, constants.APIPassword)

	utils.CheckEnvBool(&cfg.DevelopeMode, constants.DevelopeMode)

	if cfg.DevelopeMode {
		cfg.KafkaTopics.WebhookExample.TopicName += "Dev"
	}

	return cfg, nil
}
