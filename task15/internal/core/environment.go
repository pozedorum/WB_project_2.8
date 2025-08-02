package core

import (
	"errors"
	"os"
	"path/filepath"
)

type Environment interface {
	Getwd() (string, error) // Получить текущую директорию
	Chdir(dir string) error // Сменить директорию

	// Переменные окружения
	Getenv(key string) string       // Получить переменную окружения
	Setenv(key, value string) error // Установить переменную
	Environ() []string              // Получить все переменные окружения

	GetHomeDir() (string, error) // Получить домашнюю директорию
	GetBaseWd() (string, error)  // Получить текущую директорию (для красиовго вывода)
}

// DefaultEnvironment - реализация Environment по умолчанию
type DefaultEnvironment struct{}

func (e *DefaultEnvironment) Getwd() (string, error) {
	return os.Getwd()
}

func (e *DefaultEnvironment) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (e *DefaultEnvironment) GetHomeDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("HOME environment variable not set")
	}
	return home, nil
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

func (e *DefaultEnvironment) GetBaseWd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(wd), nil
}
