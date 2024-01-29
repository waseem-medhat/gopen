fmt:
	go fmt ./...

test: fmt
	go test ./...

lint: fmt
	golangci-lint run

build:
	GOOS=windows go build -o ./bin/gopen-v$(ver)-win-amd64.exe && strip ./bin/gopen-v$(ver)-win-amd64.exe
	GOOS=linux go build -o ./bin/gopen-v$(ver)-linux-amd64 && strip ./bin/gopen-v$(ver)-linux-amd64
	GOOS=darwin go build -o ./bin/gopen-v$(ver)-darwin-amd64 && strip ./bin/gopen-v$(ver)-darwin-amd64
