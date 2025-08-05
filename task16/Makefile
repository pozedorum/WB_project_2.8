

all: clean build

build:
	go build -o mywget cmd/main.go

clean:
	rm -rf ./mywget

lint:
	golangci-lint run ./...
	go vet ./...