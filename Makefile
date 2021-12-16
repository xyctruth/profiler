IMG ?= profiler:latest

docker-build:
	docker build -t ${IMG} .

docker-run:
	docker run -d -p 80:80 -p 8080:8080 --name profiler ${IMG}

docker-stop:
	docker stop profiler && docker rm profiler

docker-push:
	docker push ${IMG}