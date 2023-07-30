package config

import (
	"fmt"
	"os"

	"github.com/dominiclet/golang-base/init_server/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     string `yaml:"db"`
	Domain string `yaml:"domain"`
	Email  Email  `yaml:"email"`
}

type Email struct {
	ServerAddress string `yaml:"server_address"`
	ServerPort    int    `yaml:"server_port"`
	EmailAddress  string `yaml:"email_address"`
	AppPassword   string `yaml:"app_password"`
}

const (
	confPathEnvVar  = "CONFIG_PATH"
	defaultConfPath = "/opt/backend/config.yaml"
)

func InitConfig() *Config {
	logger := logger.GetLogger()

	var config Config
	configPath := defaultConfPath
	// Use provided config path if given in env var
	tmpConfPath, ok := os.LookupEnv(confPathEnvVar)
	if ok {
		configPath = tmpConfPath
	}
	logger.WithField("configFile", configPath).
		Info("Reading from config file")
	confData, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read config file: %v", err))
	}
	err = yaml.Unmarshal(confData, &config)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config yaml file: %v", err))
	}

	config.validateConfig()

	logger.WithFields(logrus.Fields{
		"DB": config.DB,
	}).Info("Initialized configs")

	return &config
}

func (c *Config) validateConfig() {
	if c.DB == "" {
		panic("DB field not set")
	}
	if c.Domain == "" {
		panic("domain not set")
	}
	if c.Email.ServerAddress == "" {
		panic("email.server_address not set")
	}
	if c.Email.ServerPort == 0 {
		panic("email.server_port not set")
	}
	if c.Email.EmailAddress == "" {
		panic("email.email_address not set")
	}
	if c.Email.AppPassword == "" {
		panic("email.app_password not set")
	}
}
