IMG ?= profiler:latest

.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

.PHONY: docker-run
docker-run:
	docker run -d -p 80:80 -p 8080:8080 --name profiler ${IMG}

.PHONY: docker-stop
docker-stop:
	docker stop profiler && docker rm profiler

.PHONY: test
test:
	go test -race -v -coverprofile=cover.out  $(shell go list ./pkg/... | grep -v v1175)

.PHONY: cover-ui
cover-ui: test
	go tool cover -html=cover.out -o cover.html
	open cover.html

.PHONY: fmt
fmt:
	gofmt -w $(shell find . -name "*.go")

.PHONY: lint
lint:
	golangci-lint run

