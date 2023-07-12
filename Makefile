docker_build:
	docker build . -t test-ci-cd

lint:
	golangci-lint run ./...
