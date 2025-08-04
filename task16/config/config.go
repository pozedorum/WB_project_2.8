package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config содержит параметры приложения
type Config struct {
	ResultDir    string
	WorkersCount int
	Timeout      time.Duration
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		ResultDir:    filepath.Join(".", "downloads"),
		WorkersCount: 5,
		Timeout:      time.Second * 60,
	}
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (*Config, error) {
	cfg := DefaultConfig()

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("syntax error in line %d: expected key=value", lineNum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ResultDir":
			// Удаляем кавычки если они есть
			cfg.ResultDir = strings.Trim(value, `"`)
		case "WorkersCount":
			count, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid WorkersCount in line %d: %w", lineNum, err)
			}
			cfg.WorkersCount = count
		case "Timeout":
			count, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid Timeout in line %d: %w", lineNum, err)
			}
			cfg.Timeout = time.Duration(count) * time.Second
		default:
			return nil, fmt.Errorf("unknown config key %q in line %d", key, lineNum)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return cfg, nil
}
