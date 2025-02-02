all: build-arm build-x86

build-arm:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o bin/sofar-arm

build:
	go build -o bin/sofar-x86

test:
	go test -v ./...

lint:
	docker run -t --rm -v .:/app -v ~/.cache/golangci-lint/v1.63.4:/root/.cache -w /app golangci/golangci-lint:v1.63.4 golangci-lint run -v