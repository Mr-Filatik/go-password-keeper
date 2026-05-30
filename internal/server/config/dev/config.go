package devconfig

import devconfig "github.com/mr-filatik/go-password-keeper/internal/platform/config/dev"

const (
	AccountAPI devconfig.DependName = "ACCOUNT_API"
	ClientAPI  devconfig.DependName = "CLIENT_API"
	Database   devconfig.DependName = "DATABASE"
	KafkaDWH   devconfig.DependName = "KAFKA_DWH"
)

type Dependencies struct {
	AccountAPI bool `env:"DEPENDENCY_ACCOUNT_API_ENABLED" envDefault:"true"`
	ClientAPI  bool `env:"DEPENDENCY_CLIENT_API_ENABLED"  envDefault:"true"`
	Database   bool `env:"DEPENDENCY_DATABASE_ENABLED"    envDefault:"true"`
	KafkaDWH   bool `env:"DEPENDENCY_KAFKA_DWH_ENABLED"   envDefault:"true"`
}

// DEV_CONFIG_PATH - .env (local), dev.env, test.env, prod.env, ...

func LoadDevConfig(fileName string) (*devconfig.DevConfig[Dependencies], error) {
	return devconfig.LoadConfig(
		Dependencies{},
		devconfig.WithPath("./internal/server/config/dev/"),
		devconfig.WithFileName(fileName))
}
