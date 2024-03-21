package config

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Pool                     interface{}                   // This is a placeholder, adjust according to actual use
	DatabaseService          func() (interface{}, []error) // Placeholder, adjust as needed
	ComponentConfigOverrides ComponentConfigs
	mergedComponentConfigs   ComponentConfigs
}

// newConfig initializes a new configuration struct from a file path.
func (b *builder) newConfig(configPath string) (*Config, []error) {
	var errs []error
	if file, loadErrs := b.loadConfig(configPath); loadErrs != nil {
		return nil, loadErrs
	} else {
		if readErr := b.Read(file); readErr != nil {
			return nil, []error{fmt.Errorf("newConfig: failed to read file: %v; error: %w", file.Name(), readErr)}
		}
	}

	for _, service := range b.config.Services {
		service.setClient(*service.mergedComponents().Client)
	}

	var dbErrs []error
	// initialize the Collector for each crawler

	for _, database := range b.config.Databases {
		if database.Pool, dbErrs = database.DatabaseService(); dbErrs != nil {
			errs = appendAndLog(fmt.Errorf("newConfig: %v", dbErrs), errs)
		}
	}

	return b.config, errs
}

// loadConfig attempts to load the configuration file.
func (b *builder) loadConfig(configPath string) (*os.File, []error) {
	file, loadErr := b.Load(configPath)
	if loadErr != nil {
		return nil, []error{fmt.Errorf("loadConfig: %w", loadErr)}
	}
	return file, nil
}

// Load opens the configuration file.
func (b *builder) Load(path string) (*os.File, error) {
	log.Tracef("Loading config: %v", path)
	b.configPath = path
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Load: failed to open config file %v; %w", path, err)
	}
	return file, nil
}

// Read parses the configuration file.
func (b *builder) Read(configData io.Reader) error {
	config, err := initialConfig(configData)
	if err != nil {
		return err
	}

	b.config = config

	if mergeErr := mergeServiceComponentConfigs(b.config); mergeErr != nil {
		return fmt.Errorf("Read: failed to merge component configs, error: %w", mergeErr)
	}
	return nil
}

// initialConfig initializes the configuration from reader data.
func initialConfig(data io.Reader) (*Config, error) {
	buf := new(bytes.Buffer)
	if _, buffErr := io.Copy(buf, data); buffErr != nil {
		return nil, fmt.Errorf("initialConfig: failed to read config data; err: %w", buffErr)
	}

	config := new(Config)
	if err := yaml.Unmarshal(buf.Bytes(), config); err != nil {
		return nil, fmt.Errorf("initialConfig: failed unmarshalling config data; err: %w", err)
	}
	hash := sha256.Sum256(buf.Bytes())
	config.Hash = hex.EncodeToString(hash[:])

	return config, nil
}

// mergeServiceComponentConfigs merges service and component configurations.
func mergeServiceComponentConfigs(c *Config) error {
	for i, service := range c.Services {
		if mergeErr := mergeConfigs(&service.ComponentConfigOverrides, &c.ComponentConfigs, &service.mergedComponentConfigs); mergeErr != nil {
			return fmt.Errorf("mergeServiceComponentConfigs: failed to merging component config: %v; error %w", i, mergeErr)
		}
	}
	return nil
}

// mergeConfigs merges override and default component configurations.
func mergeConfigs(override, defaultC, mergedC *ComponentConfigs) error {
	if mergedC == nil {
		return errors.New("mergeConfigs: nil pointer passed for merged components")
	}

	// First, copy the default configuration to the merged configuration.
	if defaultC != nil {
		if err := copier.Copy(mergedC, defaultC); err != nil {
			return fmt.Errorf("mergeConfigs: failed to copy default configs; error: %w", err)
		}
	}

	// Then, merge the override configuration into the merged configuration.
	if override != nil {
		if err := mergo.MergeWithOverwrite(mergedC, *override); err != nil {
			return fmt.Errorf("mergeConfigs: failed to merge overrides; error: %w", err)
		}
	}

	return nil
}

func appendAndLog(err error, errs []error) []error {
	log.Error(err)           // Log the error.
	return append(errs, err) // Append the error to the slice and return it.
}
