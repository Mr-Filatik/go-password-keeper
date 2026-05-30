// Package devconfig provides universal tools for strongly typed
// management of an application's (microservice's) infrastructure dependencies based
// on environment variables (.env files) and predefined operating modes (All/Some/None).
package devconfig

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mr-filatik/go-password-keeper/internal/platform/types"
)

const (
	dependsModeAll  = "all"
	dependsModeSome = "some"
	dependsModeNone = "none"

	// DependsModeAll denotes the mode when all external dependencies are included.
	DependsModeAll DependsMode = dependsModeAll
	// DependsModeSome denotes the mode when external dependencies are included selectively.
	DependsModeSome DependsMode = dependsModeSome
	// DependsModeNone denotes the mode when all external dependencies are disabled (replaced with mocks).
	DependsModeNone DependsMode = dependsModeNone
)

// DependsMode describes the external dependency modes for the project.
type DependsMode string

// DependName describes the name of a specific external dependency for the project.
type DependName string

// DevConfig describes the configuration for managing a project's external dependencies.
//
// In DependsModeAll mode, all external dependencies are enabled by default.
// In DependsModeNone mode, all external dependencies are disabled by default.
type DevConfig[Dependencies any] struct {
	Deps Dependencies
	Mode DependsMode `env:"DEPENDENCIES_MODE" envDefault:"all"`

	depsMap map[DependName]bool
}

// LoadConfig loads external dependency information from the env file.
//
// By default, it looks for the ".env" file in the project's root directory.
// If the file is not found, DependsModeAll is assumed to be enabled by default.
//
// Important! Recommended for local development and limited dev environments only.
// For other environments, disabling external dependencies unless absolutely necessary is not recommended.
func LoadConfig[T any](defaultDeps T, opts ...Option) (*DevConfig[T], error) {
	opt := defaultConfigOptions()
	for _, o := range opts {
		o(opt)
	}

	fullPath := filepath.Join(opt.dir, opt.filename)
	_ = godotenv.Load(fullPath)

	cfg := &DevConfig[T]{
		Deps:    defaultDeps,
		Mode:    DependsModeAll,
		depsMap: make(map[DependName]bool),
	}

	err := env.Parse(cfg)
	if err != nil {
		fullPath := opt.dir + opt.filename

		return nil, fmt.Errorf("%w: parse file (%s): %w",
			types.ErrParsing, fullPath, err)
	}

	cfg.buildDepsMap()

	return cfg, nil
}

// IsEnabled indicates whether an external dependency by a specific name is enabled.
//
// In DependsModeAll mode, all external dependencies are enabled by default.
// In DependsModeNone mode, all external dependencies are disabled by default.
func (c *DevConfig[Dependencies]) IsEnabled(name DependName) bool {
	if c.Mode == DependsModeAll {
		return true
	}

	if c.Mode == DependsModeNone {
		return false
	}

	return c.depsMap[name]
}

func (c *DevConfig[T]) buildDepsMap() {
	val := reflect.ValueOf(c.Deps)
	typ := reflect.TypeOf(c.Deps)

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Bool {
			envTag := fieldType.Tag.Get("env")
			if strings.HasPrefix(envTag, "DEPENDENCY_") && strings.HasSuffix(envTag, "_ENABLED") {
				envTag = strings.TrimPrefix(envTag, "DEPENDENCY_")
				envTag = strings.TrimSuffix(envTag, "_ENABLED")
			} else {
				envTag = toSnakeUpperCase(fieldType.Name)
			}

			c.depsMap[DependName(envTag)] = field.Bool()
		}
	}
}

func toSnakeUpperCase(s string) string {
	var result strings.Builder

	for i, rune := range s {
		if i > 0 && rune >= 'A' && rune <= 'Z' {
			result.WriteRune('_')
		}

		result.WriteRune(rune)
	}

	return strings.ToUpper(result.String())
}
