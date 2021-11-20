build:
	go build
lint:
	test -z $(gofmt -l .)
format:
	go fmt
run:
	go run github.com/oradwell/verbose-twit-banner
clean:
	rm -f verbose-twit-banner
