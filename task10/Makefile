

all: mysort

small_test: mysort
	./mysort -r tests/small_test.txt

mysort: clean
	go build -o mysort ./cmd/main.go



clean:
	rm -rf mysort