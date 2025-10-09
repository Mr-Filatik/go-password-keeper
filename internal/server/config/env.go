// Package config provides functionality for loading configuration from command-line flags and environment variables.
package config

import "os"

const (
	envNameAddress string = "ADDRESS"
)

// configEnvs - a structure containing the main environment variables for the application.
type configEnvs struct {
	address        string
	addressIsValue bool
}

// envReader is an interface for reading environment variables.
type envReader func(key string) (string, bool)

// getEnvsConfig gets values ​​from the store.
func getEnvsConfig(getenv envReader) *configEnvs {
	config := &configEnvs{
		address:        "",
		addressIsValue: false,
	}

	envAddress, ok := getenv(envNameAddress)
	if ok && envAddress != "" {
		config.address = envAddress
		config.addressIsValue = true
	}

	return config
}

// getEnvsConfigFromOS gets values ​​from environment variables.
func getEnvsConfigFromOS() *configEnvs {
	return getEnvsConfig(func(key string) (string, bool) {
		value, ok := os.LookupEnv(key)

		return value, ok
	})
}

// overrideConfigFromEnvs overrides the main config with new values.
func (c *Config) overrideConfigFromEnvs(conf *configEnvs) {
	if conf == nil {
		return
	}

	if conf.addressIsValue {
		c.Address = conf.address
	}
}
