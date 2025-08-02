# Makefile для проекта GShell

# Настройки

.PHONY: all build clean test lint install uninstall

all: build

# Сборка
build: clean
	go build -v -o GShell ./cmd/main.go


# Запуск тестов
test:
	go test -v ./... -coverprofile=coverage.out

# Покрытие кода тестами
coverage: test
	go tool cover -html=coverage.out

# Линтинг
lint:
	golangci-lint run ./...
	go vet ./...
# Очистка
clean:
	rm -rf GShell coverage.out *.txt

# Запуск (после сборки)
run: build
	./GShell

