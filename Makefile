build:
	go build
lint:
	test -z $(gofmt -l .)
format:
	go fmt
run:
	go run main.go
clean:
	rm -f verbose-twit-banner
