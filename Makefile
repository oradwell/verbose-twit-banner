build:
	go build
lint:
	test -z $(gofmt -l .)
run:
	go run main.go
clean:
	rm -f verbose-twit-banner
