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
    cd ui &&  npm run dev --base_api_url=http://localhost:8080 
```

### In Docker
```bash
    # 简单启动
    docker run -d -p 80:80 --name profiler xyctruth/profiler:latest

    # 挂载目录启动
    mkdir -vp ~/profiler/config/
    cp ./collector.yaml ~/profiler/config/
    docker run -d -p 80:80 -v ~/profiler/data/:/profiler/data/ -v  ~/profiler/config/:/profiler/config/ --name profiler xyctruth/profiler:latest
```

### In Kubernetes
`./deploy/kubernetes/`