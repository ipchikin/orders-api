package configs

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestConfig_CoreTest is the core test function
func TestConfig_CoreTest(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// ConfigTestSuite is testify suite
type ConfigTestSuite struct {
	suite.Suite
}

// TestConfig_LoadConfig tests loading config file
func (suite *ConfigTestSuite) TestConfig_LoadConfig() {
	config, err := LoadConfig("dev")
	suite.Nil(err)
	suite.NotEmpty(config.DBConfig)
}

// TestConfig_LoadConfigNoConfig tests loading no config file
func (suite *ConfigTestSuite) TestConfig_LoadConfigNoConfig() {
	config, err := LoadConfig("no")
	suite.NotNil(err)
	suite.Empty(config)
}
