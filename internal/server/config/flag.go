// Package config provides functionality for loading configuration from command-line flags and environment variables.
package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	flagNameAddress string = "address"
)

// configFlags - a structure containing the main application flags.
type configFlags struct {
	address        string
	addressIsValue bool
}

// getFlagsConfig gets the config from the specified arguments.
func getFlagsConfig(fs *flag.FlagSet, args []string) (*configFlags, error) {
	config := &configFlags{
		address:        "",
		addressIsValue: false,
	}

	argAddress := fs.String(flagNameAddress, "", "HTTP server endpoint")

	err := fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse argument %w", err)
	}

	if argAddress != nil && *argAddress != "" {
		config.address = *argAddress
		config.addressIsValue = true
	}

	return config, nil
}

// getFlagsConfigFromOS gets the flag values ​​from the application's startup arguments in the OS.
func getFlagsConfigFromOS() (*configFlags, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config, err := getFlagsConfig(fs, os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("get flag config %w", err)
	}

	return config, nil
}

// overrideConfigFromFlags overrides the main config with new values.
func (c *Config) overrideConfigFromFlags(conf *configFlags) {
	if conf == nil {
		return
	}

	if conf.addressIsValue {
		c.Address = conf.address
	}
}
