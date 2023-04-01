FROM golang:1.17.7 AS builder
WORKDIR /go/src/github.com/oradwell/verbose-twit-banner/
COPY *.go go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo \
    -o verbose-twit-banner

FROM alpine:3.17.3
RUN apk --no-cache add ca-certificates \
    && adduser -D twit
WORKDIR /home/twit/
USER twit
COPY fonts/* fonts/
COPY images/* images/
COPY --from=builder \
    /go/src/github.com/oradwell/verbose-twit-banner/verbose-twit-banner .
CMD ["./verbose-twit-banner"]
