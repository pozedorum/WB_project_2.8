package core

import (
	"io"
	"os"
)

type Environment interface {
	Getwd() (string, error) // Получить текущую директорию
	Chdir(dir string) error // Сменить директорию

	// Переменные окружения
	Getenv(key string) string       // Получить переменную окружения
	Setenv(key, value string) error // Установить переменную
	Environ() []string              // Получить все переменные окружения

	// Редиректы и потоки
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

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
