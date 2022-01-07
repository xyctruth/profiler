# Profiler

[![Go Report Card](https://goreportcard.com/badge/github.com/xyctruth/profiler?x=xyctruth)](https://goreportcard.com/report/github.com/xyctruth/profiler)
[![codecov](https://codecov.io/gh/xyctruth/profiler/branch/master/graph/badge.svg?token=YWNYJK9KQW)](https://codecov.io/gh/xyctruth/profiler)
[![Build status](https://img.shields.io/github/workflow/status/xyctruth/profiler/Server-Build/master)](https://github.com/xyctruth/profiler/actions/workflows/server-build.yml)
[![Release status](https://img.shields.io/github/v/release/xyctruth/profiler)](https://github.com/xyctruth/profiler/releases)
[![LICENSE status](https://img.shields.io/github/license/xyctruth/profiler)](https://github.com/xyctruth/profiler/LICENSE)

> [English](./README-EN.md) / [中文](./README-ZH.md)

## Profiler is a continuous profiling tool that base go pprof and go trace

- **Supported sample**
  - `trace` `fgprof` `profile` `mutex` `heap` `goroutine` `allocs` `block` `threadcreate`
- **Hot Reload**
  - Collect samples of the target service according to the configuration file
  - The collection program will monitor the changes of the configuration file and apply the changed configuration file immediately
- **Chart trend**
  - Provide charts to observe the trend of multiple service performance indicators and find the time point of performance problems
  - Each bubble is a sample file of Profile and Trace
- **Detailed analysis**
  - Click the bubbles in the charts to jump to the detailed page of Profile and Trace for further detailed analysis

### [Demo](https://profiling.jia-huang.com)

<table>
  <tr>
      <td width="50%" align="center"><b>Chart trend</b></td>
      <td width="50%" align="center"><b>Click the bubble to jump the detailed profile</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-pprof.png"/></td>
  </tr>
  <tr>
      <td width="50%" align="center"><b>Click the bubble to jump to the detailed trace</b></td>
      <td width="50%" align="center"><b>Click the bubble to jump to the detailed trance</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-trace.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-trace1.png"/></td>
  </tr>
</table>

## Get started quickly

### Local start
```bash
# run server :8080
go run server/main.go 

# run ui :80
cd ui 
npm install --registry=https://registry.npm.taobao.org 
npm run dev --base_api_url=http://localhost:8080 
```

### In Docker

```bash
# No persistence
docker run -d -p 80:80 --name profiler xyctruth/profiler:latest

# Persistence
mkdir -vp ~/profiler/config/
cp ./collector.yaml ~/profiler/config/
docker run -d -p 80:80 -v ~/profiler/data/:/profiler/data/ -v ~/profiler/config/:/profiler/config/ --name profiler xyctruth/profiler:latest
```

## Collector configuration

The `golang` program that needs to be collected and analyzed needs to provide the `net/http/pprof` endpoint and configure it in the `./collector.yaml` configuration file.

The configuration file can be updated online, and the collection program will monitor the change of the configuration file and apply the changed configuration file immediately.

### `collector.yaml`

```yaml
collector:
  targetConfigs:

    profiler-server:        # Target name
      interval: 15s         # Crawl time
      expiration: 0         # No expiration time
      host: localhost:9000  # Target service host
      profileConfigs:       # Use default configuration

    server2:
      interval: 10s
      expiration: 168h      # Expiration time seven days
      host: localhost:9000
      profileConfigs:       # Override some default configuration fields
        trace:
          enable: false
        fgprof:
          enable: false
        profile:
          path: /debug/pprof/profile?seconds=10
          enable: false
        heap:
          path: /debug/pprof/heap

```

### default configuration of `profileConfigs`

```yaml
profileConfigs:
  profile:
    path: /debug/pprof/profile?seconds=10
    enable: true
  fgprof:
    path: /debug/fgprof?seconds=10
    enable: true
  trace:
    path: /debug/pprof/trace?seconds=10
    enable: true
  mutex:
    path: /debug/pprof/mutex
    enable: true
  heap:
    path: /debug/pprof/heap
    enable: true
  goroutine:
    path: /debug/pprof/goroutine
    enable: true
  allocs:
    path: /debug/pprof/allocs
    enable: true
  block:
    path: /debug/pprof/block
    enable: true
  threadcreate:
    path: /debug/pprof/threadcreate
    enable: true
```
