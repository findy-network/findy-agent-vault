package utils

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/spf13/viper"
)

const defaultPort = "8085"

// TODO: do not allow default value in production mode
const defaultJWTSecret = "mySuperSecretKeyLol"

type Configuration struct {
	Address          string
	ServerPort       int    `mapstructure:"server_port"`
	JWTKey           string `mapstructure:"jwt_key"`
	DBHost           string `mapstructure:"db_host"`
	DBPort           int    `mapstructure:"db_port"`
	DBPassword       string `mapstructure:"db_password"`
	AgencyHost       string `mapstructure:"agency_host"`
	AgencyPort       int    `mapstructure:"agency_port"`
	AgencyCertPath   string `mapstructure:"agency_cert_path"`
	UseMockDB        bool
	UseMockAgency    bool
	GenerateFakeData bool
	UsePlayground    bool
}

func LoadConfig() *Configuration {
	defer err2.Catch(func(err error) {
		panic(fmt.Errorf("failed to read the configuration file: %s", err))
	})
	var config Configuration

	v := viper.New()
	v.SetEnvPrefix("fav")
	v.SetDefault("server_port", defaultPort)
	v.SetDefault("jwt_key", defaultJWTSecret)
	v.SetDefault("db_host", "localhost")
	v.SetDefault("db_port", 5432)
	v.SetDefault("db_password", "")
	v.SetDefault("agency_host", "localhost")
	v.SetDefault("agency_port", 50051)
	v.SetDefault("agency_cert_path", "")

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
	return &config
}
