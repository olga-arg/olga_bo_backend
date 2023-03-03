.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/cmd/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/world world/cmd/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
