# Profiler

[![Go Report Card](https://goreportcard.com/badge/github.com/xyctruth/profiler?x=xyctruth)](https://goreportcard.com/report/github.com/xyctruth/profiler)
[![codecov](https://codecov.io/gh/xyctruth/profiler/branch/master/graph/badge.svg?token=YWNYJK9KQW)](https://codecov.io/gh/xyctruth/profiler)
[![Build status](https://img.shields.io/github/workflow/status/xyctruth/profiler/Server-Build/master)](https://github.com/xyctruth/profiler/actions/workflows/server-build.yml)
[![Release status](https://img.shields.io/github/v/release/xyctruth/profiler)](https://github.com/xyctruth/profiler/releases)
[![LICENSE status](https://img.shields.io/github/license/xyctruth/profiler)](https://github.com/xyctruth/profiler/LICENSE)

> [English](./README-EN.md) / [中文](./README-ZH.md)

# 简介

**Profiler 是一个基于 go pprof 与 go trace 持续性能剖析工具**

- **支持的样本**
  - `trace` `fgprof` `profile` `mutex` `heap` `goroutine` `allocs` `block` `threadcreate`
- **配置热更新**
  - 根据配置文件收集目标服务的样本
  - 收集程序会监听配置文件的变化,即时应用变化后的配置文件
- **图表趋势**
  - 提供图表观测多个服务性能指标的趋势，找出性能问题的时间点
  - 每一个气泡都是一份 Profile 与 Trace 的样本文件
- **详细分析**
  - 点击图表中的气泡跳转到 Profile 与 Trace 的详细页面进行进一步详细的分析

## [演示](https://profiler.jia-huang.com/)

<table>
  <tr>
      <td width="50%" align="center"><b>图表趋势</b></td>
      <td width="50%" align="center"><b>点击气泡跳转 Profile 详情</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-pprof.png"/></td>
  </tr>
  <tr>
      <td width="50%" align="center"><b>点击气泡跳转 Trace 详情</b></td>
      <td width="50%" align="center"><b>点击气泡跳转 Trace 详情</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-trace.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-trace1.png"/></td>
  </tr>
</table>

# 入门

## 本地

启动服务端 端口为:8080
```bash
go run server/main.go 
```

启动前端 端口为:80
```bash
cd ui
npm install --registry=https://registry.npm.taobao.org
npm run dev --base_api_url=http://localhost:8080
```

## Docker

```bash
docker run -d -p 80:80 --name profiler xyctruth/profiler:latest
```

使用自定义的配置文件

```bash
mkdir ~/profiler-config/
cp ./collector.yaml ~/profiler-config/
docker run -d -p 80:80 -v ~/profiler-config/:/profiler/config/ --name profiler xyctruth/profiler:latest
```

使用持久化数据

```bash
docker run -d -p 80:80 -v ~/profiler-data/:/profiler/data/ --name profiler xyctruth/profiler:latest
```

## Helm

安装 Profiler chart:

```bash
helm install --create-namespace -n profiler-system profiler ./charts/profiler
```

更多说明在 [Helm docs](https://github.com/xyctruth/profiler/blob/master/charts/profiler/README.ZH.md) 中

## 收集配置

需要被收集分析的 `golang` 程序,需要提供 `net/http/pprof` 端点，并配置在 `./collector.yaml` 配置文件中。

配置文件可以在线更新，收集程序会监听配置文件的变化，即时应用变化后的配置文件。

`collector.yaml`

```yaml
collector:
  targetConfigs:

    profiler-server:        # 目标名称
      interval: 15s         # 抓取间隔
      expiration: 0         # 无过期时间
      host: localhost:9000  # 目标服务host
      labels:
        namespace: f005
        type: gateway
      profileConfigs:       # 使用默认配置 

    server2:
      interval: 10s
      expiration: 168h      # 过期时间7天
      host: localhost:9000
      labels:
        namespace: f004
        type: svc
      profileConfigs:       # 覆盖部分默认配置字段
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

`profileConfigs` 默认配置

默认 trace 分析关闭, 因为 trace 文件过大,大约在(500KB ~ 2M), 需要开启 trace 分析在 `collector.yaml` 设置覆盖默认的 trace 配置.

```yaml
profileConfigs:
  profile:
    path: /debug/pprof/profile?seconds=10
    enable: true
  fgprof:
    path: /debug/fgprof?seconds=10
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
  trace:
    path: /debug/pprof/trace?seconds=10
    enable: false
```
