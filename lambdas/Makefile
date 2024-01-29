.PHONY: build clean deploy

LAMBDAS := $(shell ls -d go*)
PYTHON_LAMBDAS := $(shell ls -d python-lambda-*/)

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

	# Package Python Lambdas
	for lambda in $(PYTHON_LAMBDAS) ; do \
		lambda_name=$$(basename $$lambda); \
		mkdir -p bin/$$lambda_name; \
		cd $$lambda; \
		zip -r ../bin/$$lambda_name/$$lambda_name.zip . ; \
		cd ..; \
	done

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
