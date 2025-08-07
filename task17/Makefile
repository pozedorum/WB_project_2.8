
all: clean build

build:
	go build -o mytelnet cmd/main.go

race_test:
	go run -race cmd/main.go tcpbin.com 4242

test:
	go run test_server/main.go

clean:
	rm -rf ./mytelnet

lint:
	golangci-lint run ./...
	go vet ./...