IMG ?= profiler:latest
BASE_API_URL ?= http://localhost:8080

docker-build:
	docker build --build-arg=BASE_API_URL=$(BASE_API_URL) -t ${IMG} .

docker-run:
	docker run -d -p 80:80 -p 8080:8080 --name profiler ${IMG}

docker-stop:
	docker stop profiler && docker rm profiler

docker-push:
	docker push ${IMG}