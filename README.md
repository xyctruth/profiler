# Continuous profiling for golang program

[![Go Report Card](https://goreportcard.com/badge/github.com/xyctruth/profiler?x=xyctruth)](https://goreportcard.com/report/github.com/xyctruth/profiler)
[![codecov](https://codecov.io/gh/xyctruth/profiler/branch/master/graph/badge.svg?token=YWNYJK9KQW)](https://codecov.io/gh/xyctruth/profiler)
[![Build status](https://img.shields.io/github/workflow/status/xyctruth/profiler/Server-Build/master)](https://github.com/xyctruth/profiler/actions/workflows/server-build.yml)

## [Demo](https://profiling.jia-huang.com)

![profiler](assets/profiler.png)

### Click Point Open Profile UI 

![profiler-pprof](assets/profiler-pprof.png)

### Click Trace Charts Point Open Trace UI

![profiler-pprof](assets/profiler-trace.png)

## Quick Start

需要被收集分析的 `golang` 程序,需要提供 `net/http/pprof` 端点，并配置在 `./collector.yaml` 配置文件中

程序会 watch `collector.yaml` 配置文件变化, 实时加载变化的配置

### Dev
```bash
     # run server :8080
    go run server/main.go 
     # run ui :80
    cd ui && npm install --registry=https://registry.npm.taobao.org  &&  npm run dev --base_api_url=http://localhost:8080 
```

### In Docker
```bash
    # No persistence
    docker run -d -p 80:80 --name profiler xyctruth/profiler:latest

    # Bind mount a volume
    mkdir -vp ~/profiler/config/
    cp ./collector.yaml ~/profiler/config/
    docker run -d -p 80:80 -v ~/profiler/data/:/profiler/data/ -v ~/profiler/config/:/profiler/config/ --name profiler xyctruth/profiler:latest
```
