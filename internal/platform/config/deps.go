package config

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	DefaultPath     = "."
	DefaultFileName = ".env"
)

type Option func(*configOptions)

type configOptions struct {
	dir      string
	filename string
}

// WithPath задает директорию, где искать env-файл (например, "./config" или "../")
func WithPath(dir string) Option {
	return func(o *configOptions) {
		if dir == "" {
			dir = DefaultPath
		}

		o.dir = dir
	}
}

// WithName задает конкретное имя файла (например, ".env.production")
func WithName(filename string) Option {
	return func(o *configOptions) {
		if filename == "" {
			filename = DefaultFileName
		}

		o.filename = filename
	}
}

func WithFullPath(fullPath string) Option {
	return func(o *configOptions) {
		dir := filepath.Dir(fullPath)
		filename := filepath.Base(fullPath)

		if dir == "" {
			dir = DefaultPath
		}

		o.dir = dir

		if filename == "" {
			filename = DefaultFileName
		}

		o.filename = filename
	}
}

type DevConfig[Dependencies any] struct {
	Mode Mode `env:"DEPENDENCIES_MODE" envDefault:"all"`
	Deps Dependencies

	depsMap map[DependencyName]bool
}

const (
	exampleDepNameAccountAPI DependencyName = "ACCOUNT_API"
	exampleDepNameClientAPI  DependencyName = "CLIENT_API"
	exampleDepNameDatabase   DependencyName = "DATABASE"
	exampleDepNameKafkaDWH   DependencyName = "KAFKA_DWH"
)

type exampleDependencies struct {
	AccountAPI bool `env:"DEPENDENCY_ACCOUNT_API" envDefault:"true"`
	ClientAPI  bool `env:"DEPENDENCY_CLIENT_API"  envDefault:"true"`
	Database   bool `env:"DEPENDENCY_DATABASE"    envDefault:"true"`
	KafkaDWH   bool `env:"DEPENDENCY_KAFKA_DWH"   envDefault:"true"`
}

func (c *DevConfig[Dependencies]) IsEnabled(name DependencyName) bool {
	if c.Mode == DepsModeAll {
		return true
	}

	if c.Mode == DepsModeNone {
		return false
	}

	return c.depsMap[name]
}

type DependencyName string

type Mode string

const (
	depsModeAll  = "all"
	depsModeSome = "some"
	depsModeNone = "none"

	// DepsModeAll - все зависимости включены
	DepsModeAll Mode = depsModeAll
	// DepsModeSome - зависимости включены выборочно
	DepsModeSome Mode = depsModeSome
	// DepsModeNone - все зависимости выключены
	DepsModeNone Mode = depsModeNone
)

func LoadConfig[T any](defaultDeps T, opts ...Option) (*DevConfig[T], error) {
	opt := configOptions{
		dir:      DefaultPath,
		filename: DefaultFileName,
	}
	for _, o := range opts {
		o(&opt)
	}

	// 1. Загружаем переменные из .env файла в окружение системы
	// Если файла нет (например, в Docker-контейнере на проде), пропускаем без ошибки
	fullPath := filepath.Join(opt.dir, opt.filename)
	_ = godotenv.Load(fullPath)

	cfg := &DevConfig[T]{
		Deps:    defaultDeps,
		depsMap: make(map[DependencyName]bool),
	}

	// 2. Парсим переменные окружения напрямую в структуру
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	cfg.buildDepsMap()

	return cfg, nil
}

func (c *DevConfig[T]) buildDepsMap() {
	val := reflect.ValueOf(c.Deps)
	typ := reflect.TypeOf(c.Deps)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if field.Kind() == reflect.Bool {
			envTag := fieldType.Tag.Get("env")
			cleanKey := envTag
			if strings.HasPrefix(envTag, "DEPENDENCY_") {
				cleanKey = strings.TrimPrefix(envTag, "DEPENDENCY_")
			} else {
				cleanKey = toSnakeUpperCase(fieldType.Name)
			}

			c.depsMap[DependencyName(cleanKey)] = field.Bool()
		}
	}
}

func toSnakeUpperCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}
