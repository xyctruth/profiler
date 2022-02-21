# Profiler charts

## 使用

安装 Profiler chart:

```bash
helm install --create-namespace -n profiler-system profiler ./charts/profiler
```

升级 Profiler chart:

```bash
helm upgrade -n profiler-system profiler ./charts/profiler
```

拆卸 Profiler chart:

```bash
helm delete -n profiler-system profiler
```

## 参数

### Configuration 参数

| Key             | Description                                                                                                                           | Value |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------|-------|
| `configuration` | Profiler 配置. 为 Profiler 生成 `collector.yaml`, 参考: [sample config](https://github.com/xyctruth/profiler/blob/master/collector.dev.yaml) | `""`  |


### Ingress 参数

| Key                            | Description                                      | Value               |
|--------------------------------|--------------------------------------------------|---------------------|
| `ingress.enabled`              | 为 Profiler 启用 ingress                            | `false`             |
| `ingress.className`            | 为 Ingress 指定一个 Ingress class。 (Kubernetes 1.18+) | `"nginx"`           |
| `ingress.annotations`          | 配置 TLS 证书的自动签发注解等.                               | `{}`                |
| `ingress.hosts.host`           | 为 Profiler 指定主机名.                                | `"your-domain.com"` |
| `ingress.hosts.paths.path`     | 为 Profiler 配置匹配路径                                | `"/"`               |
| `ingress.hosts.paths.pathType` | Ingress 匹配类型, `Prefix` 或 `Exact`                 | `"Prefix"`          |
| `ingress.tls`                  | 为 Profiler 主机名配置 TLS secret                      | `[]`                |
