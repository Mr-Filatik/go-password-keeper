package devconfig

import "github.com/mr-filatik/go-password-keeper/internal/platform/config"

const (
	AccountAPI config.DependencyName = "ACCOUNT_API"
	ClientAPI  config.DependencyName = "CLIENT_API"
	Database   config.DependencyName = "DATABASE"
	KafkaDWH   config.DependencyName = "KAFKA_DWH"
)

type Dependencies struct {
	AccountAPI bool `env:"DEPENDENCY_ACCOUNT_API" envDefault:"true"`
	ClientAPI  bool `env:"DEPENDENCY_CLIENT_API"  envDefault:"true"`
	Database   bool `env:"DEPENDENCY_DATABASE"    envDefault:"true"`
	KafkaDWH   bool `env:"DEPENDENCY_KAFKA_DWH"   envDefault:"true"`
}

// DEV_CONFIG_PATH - .env (local), dev.env, test.env, prod.env, ...

func LoadDevConfig(fileName string) (*config.DevConfig[Dependencies], error) {
	return config.LoadConfig(
		Dependencies{},
		config.WithPath("./internal/server/config/dev/"),
		config.WithName(fileName))
}
