package builtins

import "os"

// DefaultEnvironment - реализация Environment по умолчанию
type DefaultEnvironment struct{}

func (e *DefaultEnvironment) Getwd() (string, error) {
	return os.Getwd()
}

func (e *DefaultEnvironment) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (e *DefaultEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

func (e *DefaultEnvironment) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (e *DefaultEnvironment) Environ() []string {
	return os.Environ()
}
