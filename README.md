# Deprecation Notice for Kusible

Dear users and contributors of Kusible,

We want to inform you that Kusible is being deprecated and will no longer be actively maintained. After careful consideration, the project maintainers have decided to discontinue further development and support for Kusible.

**What does this mean for you?**

1. **No Further Updates:** There will be no further updates, enhancements, or bug fixes to Kusible. The repository will remain as is, but we will not actively maintain or contribute to it.

2. **Security and Compatibility:** As Kusible will no longer be actively maintained, it may become vulnerable to security issues or compatibility problems with newer dependencies or Kubernetes versions.

3. **Consider Alternatives:** We encourage you to consider alternatives for your Kubernetes deployment needs, as Kusible will no longer be a recommended solution. Popular alternatives in the Kubernetes deployment space include Helm, Kustomize, and other similar tools.

4. **Archival:** The Kusible Git repository will be archived, but it will remain accessible for historical purposes.

We understand that this decision may impact users who have found value in Kusible, and we sincerely thank you for your support and contributions. We recommend transitioning to alternative solutions that are actively maintained and aligned with your requirements.


![Logo of kusible](assets/images/kusible-0.0.4-small.png)

# Kusible

kusible is a multi-kubernetes-helm-chart-deployment tool. It uses YAML files which describe chart repositories,
chart versions, chart settings, kubernetes clusters and so on. It uses this information to deploy charts to multiple
clusters, each with its own version and settings.

File structure and certain keywords are intentionally similar to Ansible. Differences result from the different needs
of the tasks to be performed and the tools used.

Kusible makes heavy use of [Spruce](https://github.com/geofffranks/spruce) to manipulate YAML and JSON files so understanding
of [Spruce Operators](https://github.com/geofffranks/spruce/blob/master/doc/operators.md) and [YAML anchors](https://learnxinyminutes.com/docs/yaml/)
is useful. Usage of [ejson](https://github.com/Shopify/ejson) to store encrypted data is also supported.

# Sample application structure

```
.
├── group_vars
│   ├── all.yml
│   ├── dev.yml
│   ├── prod.yml
│   └── test.yml
├── inventory
│   └── clusters.yml
├── playbook.yml
└── kubeconfig_dev
```

playbook.yml:
```
---
plays:
  - name: install-prometheus-operator
    groups: [all]
    charts:
    - name: prometheus-operator
      repo: stable
      chart: prometheus-operator
      version: (( grab vars.VERSIONS.prometheus-operator ))
      namespace: kube-system
      values: (( grab vars.VALUES ))
    repos: (( grab vars.REPOS ))
```

inventory/clusters.yml:
```yaml
---
inventory:
  - name: k8s-dev
    groups: [dc1,dev]
    kubeconfig:
      backend: file
      params:
        path: kubeconfig_dev
  - name: k8s-test
    groups: [dc2,test]
    kubeconfig:
      backend: s3
      params:
        accesskey: "test_accesskey"
        secretkey: "test_secretkey"
        server: s3.somewhere.com
        decrypt_key: verysecretkey
        bucket: kubernetes
        path: k8s-test/kubeconfig.enc.7z
  - name: k8s-prod
    groups: [dc3,prod]
    kubeconfig:
      backend: s3
      params:
        accesskey: "prod_accesskey"
        secretkey: (( vault "secret/my/credentials/admin:password" ))
        server: s3.somewhere.com
        decrypt_key: (( vault "secret/my/credentials/admin:decryptkey" ))
        bucket: kubernetes
        path: k8s-prod/kubeconfig.enc.7z
```

group_vars/all.yml
```yaml
---
vars:
  REPOS:
    - name: stable
      url: https://charts.helm.sh/stable
  VALUES:
    defaultRules:
      create: false
    alertmanager:
      enabled: false
    grafana:
      enabled: false
    kubeApiServer:
      enabled: false
    kubelet:
      enabled: false
    kubeControllerManager:
      enabled: false
    coreDns:
      enabled: false
    kubeDns:
      enabled: false
    kubeEtcd:
      enabled: false
    kubeScheduler:
      enabled: false
    kubeStateMetrics:
      enabled: false
    nodeExporter:
      enabled: false
    prometheus:
      enabled: false
    prometheusOperator:
      kubeletService:
        enabled: true
      image:
        repository: coreos/prometheus-operator
      configmapReloadImage:
        repository: coreos/configmap-reload
      prometheusConfigReloaderImage:
        repository: coreos/prometheus-config-reloader
      hyperkubeImage:
        repository: google-containers/hyperkube
```

group_vars/dev.yml
```yaml
---
vars:
  VERSIONS:
    prometheus-operator: ">0.0.0-0"
```

group_vars/test.yml
```yaml
---
vars:
  VERSIONS:
    prometheus-operator: ">0.0.0"
```

group_vars/prod.yml
```yaml
---
vars:
  VERSIONS:
    prometheus-operator: "5.12.3"
```

# Command line

```
This is a CLI tool to render and deploy kubernetes resources to multiple clusters

Usage:
  kusible [command]

Available Commands:
  deploy      Deploy an application
  groups      List available groups based on given regex
  help        Help about any command
  inventory   Get inventory information
  render      Render an application as kubernetes manifests
  uninstall   Uninstall an application
  values      Compile values for a list of groups
  version     Print the version number of kusible

Flags:
  -h, --help               help for kusible
      --json-log           log as json
      --log-functions      log function names (performance impact!)
      --log-level string   log level (trace,debug,info,warn/warning,error,fatal,panic) (default "warning")
  -q, --quiet              Suppress all normal output

Use "kusible [command] --help" for more information about a command.
```

# Documentation

See [documentation](https://bedag.github.io/kusible/)

# License

Licensed under [the Apache 2.0 license](LICENSE)
