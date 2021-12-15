# Build Server
FROM golang:1.16 as builder
WORKDIR /workspace
COPY ./ ./
RUN GOPROXY="https://goproxy.io,direct"  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o profiler ./server/main.go

# Build Ui
FROM node:16 as ui-builder
WORKDIR /workspace
COPY ui ./
ARG BASE_API_URL
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build --base_api_url=${BASE_API_URL}


# profiler image
FROM nginx:alpine
WORKDIR /profiler

RUN apk add graphviz

# server
COPY --from=builder /workspace/profiler .
COPY collector.yaml ./collector.yaml

# ui
COPY --from=ui-builder /workspace/dist /usr/share/nginx/html/
COPY --from=ui-builder /workspace/nginx.conf /etc/nginx/conf.d/default.conf

COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x entrypoint.sh
EXPOSE 80 8080

ENTRYPOINT ["./entrypoint.sh"]
