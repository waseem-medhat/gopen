fmt:
	go fmt ./...

test: fmt
	go test ./...

lint: fmt
	golangci-lint run

build:
	GOOS=windows tinygo build -o ./bin/gopen-win-amd64 && strip ./bin/gopen-win-amd64
	GOOS=linux tinygo build -o ./bin/gopen-linux-amd64 && strip ./bin/gopen-linux-amd64
	GOOS=darwin tinygo build -o ./bin/gopen-darwin-amd64 && strip ./bin/gopen-darwin-amd64
