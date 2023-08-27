.PHONY: build clean deploy

LAMBDAS := $(shell ls -d go*)

build:
	mkdir -p bin
	for lambda in $(LAMBDAS) ; do \
		lambda_name=$$(basename $$lambda); \
		echo $$lambda_name; \
		cd $$lambda; \
		env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/$$lambda_name/bootstrap cmd/main.go ; \
		cd ../bin/$$lambda_name/ ; \
		zip $$lambda_name.zip bootstrap; \
		cd - ; \
		cd ../ ; \
	done

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
