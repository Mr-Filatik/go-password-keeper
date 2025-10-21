// Package config provides functionality for loading configuration from command-line flags and environment variables.
package config

import "os"

const (
	envNameServerAddress string = "SERVER_ADDRESS"
)

// configEnvs - a structure containing the main environment variables for the application.
type configEnvs struct {
	serverAddress        string
	serverAddressIsValue bool
}

// envReader is an interface for reading environment variables.
type envReader func(key string) (string, bool)

// getEnvsConfig gets values ​​from the store.
func getEnvsConfig(getenv envReader) *configEnvs {
	config := &configEnvs{
		serverAddress:        "",
		serverAddressIsValue: false,
	}

	envAddress, ok := getenv(envNameServerAddress)
	if ok && envAddress != "" {
		config.serverAddress = envAddress
		config.serverAddressIsValue = true
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

	if conf.serverAddressIsValue {
		c.Address = conf.serverAddress
	}
}
