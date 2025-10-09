// Package config provides functionality for loading configuration from command-line flags and environment variables.
package config

// Constants are default values.
const (
	defaultAddress string = ":8080"
)

// Config is a structure containing the main parameters of the application.
type Config struct {
	// Address - server startup address.
	Address string
}

// Initialize creates and initializes a *Config object.
//
// Values ​​are assigned (reassigned) in the following order:
// - default values;
// - values ​​from command-line flags;
// - values ​​from environment variables.
func Initialize() *Config {
	envsConf := getEnvsConfigFromOS()
	flagsConf, _ := getFlagsConfigFromOS()

	config := createAndOverrideConfig(flagsConf, envsConf)

	return config
}

func createAndOverrideConfig(flagsConf *configFlags, envsConf *configEnvs) *Config {
	config := &Config{
		Address: defaultAddress,
	}

	config.overrideConfigFromFlags(flagsConf)
	config.overrideConfigFromEnvs(envsConf)

	return config
}
