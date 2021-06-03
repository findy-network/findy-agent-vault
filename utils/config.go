package utils

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/spf13/viper"
)

const defaultPort = "8085"
const localhost = "localhost"

var Version = "dev"

// TODO: do not allow default value in production mode
const defaultJWTSecret = "mySuperSecretKeyLol"

type Configuration struct {
	// true if this vault service is the main agency subscriber
	// TODO: this need to be rethought when we are scaling vault
	AgencyMainSubscriber bool   `mapstructure:"agency_main_subscriber"`
	AgencyCertPath       string `mapstructure:"agency_cert_path"`
	AgencyHost           string `mapstructure:"agency_host"`
	AgencyPort           int    `mapstructure:"agency_port"`
	AgencyAdminID        string `mapstructure:"agency_admin_id"`
	Address              string
	DBHost               string `mapstructure:"db_host"`
	DBPassword           string `mapstructure:"db_password"`
	DBPort               int    `mapstructure:"db_port"`
	DBTracing            bool   `mapstructure:"db_tracing"`
	GenerateFakeData     bool
	JWTKey               string `mapstructure:"jwt_key"`
	LogLevel             string `mapstructure:"log_level"`
	ServerPort           int    `mapstructure:"server_port"`
	UseMockDB            bool
	UsePlayground        bool `mapstructure:"use_playground"`
	Version              string
}

func LoadConfig() *Configuration {
	defer err2.Catch(func(err error) {
		panic(fmt.Errorf("failed to read the configuration file: %s", err))
	})
	var config Configuration

	v := viper.New()
	v.SetEnvPrefix("fav")
	v.SetDefault("agency_main_subscriber", true)
	v.SetDefault("agency_cert_path", "")
	v.SetDefault("agency_host", localhost)
	v.SetDefault("agency_port", 50051)
	v.SetDefault("agency_admin_id", "findy-root")
	v.SetDefault("db_host", localhost)
	v.SetDefault("db_password", "")
	v.SetDefault("db_port", 5432)
	v.SetDefault("db_tracing", false)
	v.SetDefault("jwt_key", defaultJWTSecret)
	v.SetDefault("log_level", "3")
	v.SetDefault("server_port", defaultPort)
	v.SetDefault("use_playground", false)

	viper.SetConfigName("config.yaml")
	viper.AddConfigPath(".")
	v.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.Info("Configuration file was not found, using environment/default variables only")
		} else {
			err2.Check(err)
		}
	}
	err2.Check(v.Unmarshal(&config))

	config.Address = fmt.Sprintf(":%d", config.ServerPort)
	SetLogConfig(&config)
	config.Version = Version

	// make sure we do not accidentally subscribe to the data pump when developing in local
	if config.AgencyMainSubscriber && config.AgencyHost != localhost && config.DBHost == localhost {
		config.AgencyMainSubscriber = false
	}
	return &config
}
