# Build Server
FROM golang:1.17 as builder
WORKDIR /workspace
COPY ./ ./
ARG VERSION
ARG GITVERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
    go build -a -ldflags "-s -w -X github.com/xyctruth/profiler/version.Version=${VERSION:-undefined} -X github.com/xyctruth/profiler/version.GitRevision=${GITVERSION:-undefined}" \
    -o profiler ./server/main.go

# Build Ui
FROM node:16 as ui-builder
WORKDIR /workspace
COPY ui ./
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build


# profiler image
FROM nginx:alpine
WORKDIR /profiler

RUN apk update
RUN apk add graphviz
RUN apk add dumb-init

# server
COPY --from=builder /workspace/profiler .
COPY collector.yaml ./config/collector.yaml
# go trace assets
COPY pkg/internal/v1175/assets/trace_viewer_full ./pkg/internal/v1175/assets/trace_viewer_full
COPY pkg/internal/v1175/assets/webcomponents.min.js ./pkg/internal/v1175/assets/webcomponents.min.js

# ui
COPY --from=ui-builder /workspace/dist /usr/share/nginx/html/
COPY --from=ui-builder /workspace/nginx.conf /etc/nginx/nginx.conf

#env
ENV PROFILER_API_URL="127.0.0.1:8080"
ENV DATA_PATH=/profiler/data
ENV CONFIG_PATH=/profiler/config/collector.yaml
ENV DATA_GC_INTERNAL=5m
ENV UI_GC_INTERNAL=1m

COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x entrypoint.sh
EXPOSE 80 8080

ENTRYPOINT ["dumb-init", "--"]
CMD ["./entrypoint.sh"]
