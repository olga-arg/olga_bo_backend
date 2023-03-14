.PHONY: build clean deploy

LAMBDAS := $(shell ls -d */)

build:
	mkdir -p bin
	for lambda in $(LAMBDAS) ; do \
		lambda_name=$$(basename $$lambda); \
		cd $$lambda; \
		env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/$$lambda_name/$$lambda_name cmd/main.go ; \
		cd ..; \
	done

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
