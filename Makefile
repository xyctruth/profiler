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

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: test
test:
	go test -v -coverprofile=cover.out  ./pkg/...
	go test -v  ./...

.PHONY:
cover-ui:
	go test -v -coverprofile=cover.out  ./pkg/...
	go tool cover -html=cover.out -o cover.html
	open cover.html

.PHONY: fmt
fmt:
	gofmt -w $(shell find . -name "*.go")