# 基于 pprof 的 Golang 程序连续分析

## [Demo](https://profiling.jia-huang.com)

![profiler](https://xtruth.oss-cn-shenzhen.aliyuncs.com/5.png)
 
### 点击 point
![speedscope](https://xtruth.oss-cn-shenzhen.aliyuncs.com/6.png)


## Quick Start

需要被收集分析的golang程序,需要提供`net/http/pprof`端点，并配置在`collector.yaml`配置文件中

```bash
     #run server :8080
    go run server/main.go 
     #run ui :80
    cd ui && npm run dev
```

### In Docker
```bash
    make docker-build docker-run
```

### In Kubernetes
`./deploy/kubernetes/`