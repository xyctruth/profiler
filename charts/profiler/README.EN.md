# Profiler charts

## Usage

Install the Profiler chart:

```bash
helm install --create-namespace -n profiler-system profiler ./charts/profiler
```

Upgrade the Profiler chart:

```bash
helm upgrade -n profiler-system profiler ./charts/profiler
```

Uninstall the Profiler chart:

```bash
helm delete -n profiler-system profiler
```

## Parameters

### Configuration parameters

| Key             | Description                                                                                                                                             | Value |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|-------|
| `configuration` | Profiler configuration. Specify content for `collector.yaml`, ref: [sample config](https://github.com/xyctruth/profiler/blob/master/collector.dev.yaml) | `""`  |


### Ingress parameters

| Key                            | Description                                                                      | Value               |
|--------------------------------|----------------------------------------------------------------------------------|---------------------|
| `ingress.enabled`              | Enable ingress record generation for Profiler                                    | `false`             |
| `ingress.className`            | IngressClass that will be be used to implement the Ingress (Kubernetes 1.18+)    | `"nginx"`           |
| `ingress.annotations`          | To enable certificate auto generation, place here your cert-manager annotations. | `{}`                |
| `ingress.hosts.host`           | Default host for the ingress record.                                             | `"your-domain.com"` |
| `ingress.hosts.paths.path`     | Default path for the ingress record                                              | `"/"`               |
| `ingress.hosts.paths.pathType` | Ingress path type                                                                | `"Prefix"`          |
| `ingress.tls`                  | Enable TLS configuration for the host defined at ingress.hostname parameter      | `[]`                |