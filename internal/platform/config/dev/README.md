To use it, you need to:

Рекомендуется все следующие файлы располагать в директории:
`[module]/internal/[app]/config/dev`.

1. Код config.go

1.1. Include the devconfig package.

```go
import devconfig "github.com/mr-filatik/go-password-keeper/internal/platform/config/dev"
```

1.2. Specify dependency names in uppercase (separate words with underscores).

```go
const (
	AccountAPI devconfig.DependName = "ACCOUNT_API"
	ClientAPI  devconfig.DependName = "CLIENT_API"
	Database   devconfig.DependName = "DATABASE"
	KafkaDWH   devconfig.DependName = "KAFKA_DWH"
)
```

Names must be written strictly in uppercase in snake letters (for example DEPEND_NAME).

1.3. Create a config file specifying the dependencies themselves.

```go
type Dependencies struct {
	AccountAPI bool `env:"DEPENDENCY_ACCOUNT_API_ENABLED" envDefault:"true"`
	ClientAPI  bool `env:"DEPENDENCY_CLIENT_API_ENABLED"  envDefault:"true"`
	Database   bool `env:"DEPENDENCY_DATABASE_ENABLED"    envDefault:"true"`
	KafkaDWH   bool `env:"DEPENDENCY_KAFKA_DWH_ENABLED"   envDefault:"true"`
}
```

Names in the env tag must be constructed according to the formula:
`DEPENDENCY_` + DEPEND_NAME + `_ENABLED`.

And have an envDefault tag with the "true" parameter.

1.4. Create a wrapper function that will search for files for the current project.

```go
func LoadDevConfig(fileName string) (*devconfig.DevConfig[Dependencies], error) {
	return devconfig.LoadConfig(
		Dependencies{},
		devconfig.WithPath("./internal/[app]/config/dev/"),
		devconfig.WithFileName(fileName))
}
```

2. Конфиг .env

Рекомендуется создать и .example.env на случай, если .env добавлен в .gitignore.

```
# DEPENDENCIES_MODE all, some and none
DEPENDENCIES_MODE=some

# Если MODE="some", эти значения перепишутся в false автоматически
DEPENDENCY_ACCOUNT_API_ENABLED=true
DEPENDENCY_CLIENT_API_ENABLED=false
DEPENDENCY_DATABASE_ENABLED=true
DEPENDENCY_KAFKA_DWH_ENABLED=false
```

3. Use in code as follows:

```go
import devconfig "[module]/internal/[app]/config/dev"

// ...

    devConfig, devErr := devconfig.LoadDevConfig("")
	if devErr != nil {
		panic(devErr)
	}

    // ...

    if devConfig.IsEnabled(devconfig.Database) {
	    // create database
    } else {
	    // create database mock
    }

// ...
```

Instead of an empty string, you can pass a parameter from the environment variables; it is recommended to name them DEV_CONFIG_PATH or DEV_CONFIG_FILENAME.


TODO !!!
Consider externalizing common dependencies, such as TRACER, DATABASE, KAFKA, REDIS, and other common names.
You'll need to create specific dependencies yourself, such as DATABASE_USER, DATABASE_OPERATION, and so on.

```go
const (
	AccountAPI devconfig.DependName = "ACCOUNT_API"
	ClientAPI  devconfig.DependName = "CLIENT_API"
	Database   devconfig.DependName = devconfig.Database
	Kafka      devconfig.DependName = devconfig.Kafka
)
```
