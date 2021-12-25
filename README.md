# Profiler

[![Go Report Card](https://goreportcard.com/badge/github.com/xyctruth/profiler?x=xyctruth)](https://goreportcard.com/report/github.com/xyctruth/profiler)
[![codecov](https://codecov.io/gh/xyctruth/profiler/branch/master/graph/badge.svg?token=YWNYJK9KQW)](https://codecov.io/gh/xyctruth/profiler)
[![Build status](https://img.shields.io/github/workflow/status/xyctruth/profiler/Server-Build/master)](https://github.com/xyctruth/profiler/actions/workflows/server-build.yml)

> [Demo](https://profiling.jia-huang.com)

`profiler` 基于 `pprof` 与 `go trace` 持续性能剖析工具

- 根据配置文件收集目标服务的样本, 收集程序会监听配置文件变化实时更新收集目标
- 支持的样本 `trace` `fgprof` `profile` `mutex` `heap` `goroutine` `allocs` `block` `threadcreate`
- 提供图表观测服务性能指标的趋势，找出性能问题的时间点
- 点击图标中的气泡跳转到 `pprof` 与 `trace` 的详细页面进行进一步详细的分析

<table>
  <tr>
      <td width="50%" align="center"><b>pprof图表</b></td>
      <td width="50%" align="center"><b>点击气泡跳转pprof详情</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-pprof.png"/></td>
  </tr>
  <tr>
      <td width="50%" align="center"><b>trace图表</b></td>
      <td width="50%" align="center"><b>点击气泡跳转trace详情</b></td>
  </tr>
  <tr>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler.png"/></td>
     <td><img src="https://xtruth.oss-cn-shenzhen.aliyuncs.com/profiler-trace.png"/></td>
  </tr>
</table>

## 快速入门

### 配置文件

需要被收集分析的 `golang` 程序,需要提供 `net/http/pprof` 端点，并配置在 `./collector.yaml` 配置文件中

### 本地启动
```bash
# run server :8080
go run server/main.go 

# run ui :80
cd ui && npm install --registry=https://registry.npm.taobao.org && npm run dev --base_api_url=http://localhost:8080 
```

### In Docker

```bash
# 无持久化
docker run -d -p 80:80 --name profiler xyctruth/profiler:latest

# 持久化
mkdir -vp ~/profiler/config/
cp ./collector.yaml ~/profiler/config/
docker run -d -p 80:80 -v ~/profiler/data/:/profiler/data/ -v ~/profiler/config/:/profiler/config/ --name profiler xyctruth/profiler:latest
```