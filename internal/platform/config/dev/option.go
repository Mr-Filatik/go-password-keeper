package devconfig

import "path/filepath"

const (
	defaultFileName = ".env"
	defaultPath     = "." // ".", "./", "../".
)

// Option describes additional options for working with the dev config.
type Option func(*configOptions)

type configOptions struct {
	dir      string
	filename string
}

func defaultConfigOptions() *configOptions {
	return &configOptions{
		dir:      defaultPath,
		filename: defaultFileName,
	}
}

// WithPath specifies the directory to look for the configuration env file.
//
// It is recommended to place files in ".internal/[app]/config/dev/".
func WithPath(dir string) Option {
	return func(confOpts *configOptions) {
		if dir == "" {
			dir = defaultPath
		}

		confOpts.dir = dir
	}
}

// WithFileName specifies the name of the configuration file (e.g., ".env", "dev.env", etc.).
//
// It is recommended to use the following file names:
// - ".env" - for local settings;
// - "dev.env" - for settings in the development environment.
func WithFileName(filename string) Option {
	return func(confOpts *configOptions) {
		if filename == "" {
			filename = defaultFileName
		}

		confOpts.filename = filename
	}
}

// WithFullPath specifies the full path to the configuration file, including the directory and its name.
//
// It is recommended to place files in ".internal/[app]/config/dev/.env".
func WithFullPath(fullPath string) Option {
	return func(confOpts *configOptions) {
		dir := filepath.Dir(fullPath)
		filename := filepath.Base(fullPath)

		if dir == "" {
			dir = defaultPath
		}

		confOpts.dir = dir

		if filename == "" {
			filename = defaultFileName
		}

		confOpts.filename = filename
	}
}
