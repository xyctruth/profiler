# Build Server
FROM golang:1.16 as builder
WORKDIR /workspace
COPY ./ ./
RUN GOPROXY="https://goproxy.io,direct"  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o profiler ./server/main.go

# Build Ui
FROM node:16 as ui-builder
WORKDIR /workspace
COPY ui ./
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build


# profiler image
FROM nginx:alpine
WORKDIR /profiler

RUN apk add graphviz

# server
COPY --from=builder /workspace/profiler .
COPY collector.yaml ./config/collector.yaml

# ui
COPY --from=ui-builder /workspace/dist /usr/share/nginx/html/
COPY --from=ui-builder /workspace/nginx.conf /etc/nginx/nginx.conf

#env
ENV PROFILER_API_URL="127.0.0.1:8080"
ENV DATA_PATH=/profiler/data
ENV CONFIG_PATH=/profiler/config

COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x entrypoint.sh
EXPOSE 80 8080

ENTRYPOINT ["./entrypoint.sh"]
