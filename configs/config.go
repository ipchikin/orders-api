package configs

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config of the app
type Config struct {
	DBConfig  `mapstructure:"db_config"`
	APIConfig `mapstructure:"api_config"`
}

// DBConfig of connection details
type DBConfig struct {
	Driver       string
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

// APIConfig of all used APIs
type APIConfig struct {
	DistanceMatrixAPI `mapstructure:"distance_matrix_api"`
}

// DistanceMatrixAPI config
type DistanceMatrixAPI struct {
	URL    string
	Method string
	Key    string
}

// LoadConfig loads config file according to the mode
func LoadConfig(mode string) (config Config, err error) {
	viper.SetConfigName(mode)
	viper.SetConfigType("yaml")

	_, caller, _, ok := runtime.Caller(0)
	if !ok {
		err = errors.New("caller not found")
		return
	}
	viper.AddConfigPath(filepath.Dir(caller))
	viper.AddConfigPath(".")

	// Read config file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal to config struct
	err = viper.Unmarshal(&config)

	// Load env vars
	config.DBConfig.Password = os.Getenv("MYSQL_ROOT_PASSWORD")
	config.DBConfig.Database = os.Getenv("MYSQL_DATABASE")
	config.DistanceMatrixAPI.Key = os.Getenv("MAPS_API_KEY")

	return
}
