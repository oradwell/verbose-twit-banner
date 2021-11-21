build: verbose-twit-banner

verbose-twit-banner: *.go
	go build

.PHONY: lint
lint:
	test -z $(gofmt -l .)

.PHONY: format
format:
	go fmt

.PHONY: test
test:
	go test

.PHONY: run
run:
	go run github.com/oradwell/verbose-twit-banner

.PHONY: clean
clean:
	rm -f verbose-twit-banner
