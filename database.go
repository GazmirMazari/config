package config

import (
	"fmt"
	"io"
	"os"
)


type builder struct {
	config     *Config
	configPath string
}

func (b *builder) loadConfigFile(configPath string) (*os.File, []error) {
	if file, loadErr := b.Load(configPath); loadErr != nil {
		return nil, []error{fmt.Errorf("newConfig: %w", loadErr)}
	} else {
		return file, nil
	}
}

func (b *builder) ReadConfig(configData io.Reader) (err error) {
	if b.config, err = initialConfig(configData); err != nil {
		return err
	}

	if mergeErr := mergeServiceComponentConfigs(b.config); mergeErr != nil {
		return fmt.Errorf("ReadConfig: failed to merge component configs, error: %w", mergeErr)
	}
	return nil
}
