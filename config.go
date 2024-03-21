package config

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	ENV              string            `yaml:"Env"`
	AppName          string            `yaml:"AppName"`
	Port             string            `yaml:"Port"`
	ComponentConfigs ComponentConfigs  `yaml:"ComponentsConfigs"`
	Databases        DatabaseConfigMap `yaml:"Databases"`
	Services         ServiceConfigMap  `yaml:"Services"`
	Hash             string            `yaml:"Hash"`
}

type ComponentConfigs struct {
	Client *ClientConfig
}

type conflig int

func New(configPath string) (config *Config) {
	log.Tracef("config: %s\n", configPath)
	var errs []error
	if config, errs = new(builder).newConfig(configPath); len(errs) > 0 || config == nil {
		for _, err := range errs {
			log.Panicf("configuration error: %v\n", err.Error())
		}
		if config == nil {
			log.Panicln("configuration file not found")
		}
		log.Panicln("Exiting: failed to load the config file")
	}
	log.Tracef("env: %s\n", strings.ToUpper(config.ENV))
	return config
}

func (c *Config) Service(name string) (*ServiceConfig, error) {
	if service, ok := c.Services[name]; ok {
		return service, nil
	}
	//return error if the service was not found in config

	return nil, fmt.Errorf("Service : %s", fmt.Sprintf("%s not found"))

}

// Database returns an initialized database configuration by name

// Database returns an initialized database configuration by name
func (c *Config) Database(name string) (*DatabaseConfig, error) {
	if database, ok := c.Databases[name]; ok {
		return database, nil
	}
	// return error if the database not found in config
	return nil, fmt.Errorf("Database: %s", fmt.Sprintf("%s not found", name))
}
