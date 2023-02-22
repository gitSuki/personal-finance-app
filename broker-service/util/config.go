package util

import "github.com/spf13/viper"

// Config stores all the configuration for the service's environmenment variables
type Config struct {
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

// LoadConfig uses viper to load the environment variables and returns a config struct
func LoadConfig(path string) (Config, error) {
	config := Config{}
	viper.AddConfigPath(path)
	viper.SetConfigName("broker-service")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
