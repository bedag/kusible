# Introduction

kusible is a multi-kubernetes-helm-chart-deployment tool. It uses YAML files which describe chart repositories,
chart versions, chart settings, kubernetes clusters and so on. It uses this information to deploy charts to multiple
clusters, each with its own version and settings.

File structure and certain keywords are intentionally similar to Ansible. Differences result from the different needs
of the tasks to be performed and the tools used.

Kusible makes heavy use of [Spruce](https://github.com/geofffranks/spruce) to manipulate YAML and JSON files so understanding
of [Spruce Operators](https://github.com/geofffranks/spruce/blob/master/doc/operators.md) and [YAML anchors](https://learnxinyminutes.com/docs/yaml/)
is useful. Usage of [ejson](https://github.com/Shopify/ejson) to store encrypted data is also supported.

## Main parts of a kusible project

### The inventory

The kusible inventory describes the different kubernetes clusters which serve as deploy targets for the different
helm charts. Each cluster has a name, a list of (arbitrary) groups it belongs to and a source for a kubeconfig
used to access the cluster.

```yaml
---
inventory:
  - name: %cluster-name%
    groups: [all, <cluster-name>]
    config_namespace: kube-config
    kubeconfig:
      backend: s3
      params:
        accesskey: $S3_ACCESSKEY
        secretkey: $S3_SECRETKEY
        server: $S3_SERVER
        decrypt_key: $EJSON_PRIVKEY
        bucket: kubernetes
        path: <cluster-name>/kubeconfig/kubeconfig.enc.7z
```

With the exception of "`%cluster-name%`" the values shown here are defaults used if no values are provided in the inventory (`$S3_ACCESSKEY`,
`$S3_SECRETKEY` and `$S3_SERVER` are environment variables). Given that the environment variables `$S3_ACCESSKEY`, `$S3_SECRETKEY` and `$S3_SERVER`
are set and no groups other than "`all`" and "`%cluster-name%`" are required, an inventory entry can be written as

```yaml
---
inventory:
  - name: <cluster-name>
```

Currently there are two kubeconfig backends: s3 and file. S3 is the default and supports plain yaml and encrypted tar.7z files kubeconfig files. The
inventory syntax for the s3 backend can be seen above. If the kubeconfig file is encrypted, it is assumed it uses the same key as the ejson
files in the group vars, which is provided using the `-e` cli option. Alternatively it can be specified in the `decrypt_key:` parameter.

The file backend has the following syntax:

```yaml
  kubeconfig:
    backend: file
    params:
      decrypt_key: $EJSON_PRIVKEY
      path:
```

#### Inventory location

The default inventory file is `inventory.yml`. This can be changed with the `-i` cli parameter. The inventory can be a file or a directory (including
subdirectories). Spruce operators, yaml anchors / references and ejson encrypted files work as expected.

#### Kubeconfig and Kubernetes cluster requirements

The kubeconfig is expected to only contain a single cluster and a single user.

Each cluster requires a "cluster-inventory" configmap in the config namespace provided to kusible with the `-c` cli flag. The kubeconfig for the
cluster must contain credentials allowing access to this configmap.

### The group variables

Group variables are stored in the `group_vars` directory. Each group assigned to a cluster in the inventory can have its own file or directory.
Group vars files must end in `.yml` or (for ejson encrypted files) in `.ejson`. `.yml` and `.ejson` files can exist at the same time (so
`group_vars/all.yml` and `group_vars/all.ejson` can exist at the same time). Instead of files a group can have its own subdirectory, for example
`group_vars/all/`. In this case all files (including files in subdirectories) will be used. File and directory group_vars can be used together.

All group variables belonging to a cluster will be merged in the order in which the groups are assigned to the cluster where the `all` group
has the lowest priority and the group named like the cluster has the highest priority.

If ejson encrypted files are present, the ejson privkey must be provided with the `-e` cli option.

Group vars can make use of spruce operators and can use this to access settings in the inventory config map of the given cluster.

All group variables should be inside the `vars` hash map e.g.:

```yaml
---
vars:
  var1: foo
  var2: bar
```

### Playbooks

Playbooks tie the group variables and the inventory together and define which chart gets deployed on which clusters. Each playbook consists
of one or multiple plays. Each play will be deployed to a list of groups (as defined in the inventory) and defines a list of charts and repositories
to be used. The general structure of a playbook is as follows:

```yaml
---
plays:
  - name:
    groups: []

    charts:
    - name:
      repo:
      chart:
      version:
      namespace:
      values:

    repos:
      - name:
        url:
```

With the exception of the `groups` field, spruce operators can be used. This is especially necessary to access the group variables, as they
must be accessed by using `(( grab vars. ))` (all group vars are in the `vars` hash map).

Groups in the `groups` field are **OR** associated, meaning that a playbook applies to all clusters in each group.

## The kusible CLI

TBD

### Limits

The `-l` parameters limits the operation to a subset of clusters in the inventory. For example using `-l foo` would
limit the operation to all clusters in the `foo` group. The `-l` parameter can be specified multiple times, all limits
are **AND** associated, meaning that only clusters that have all specified groups will be selected. The value of the parameter
is a regex implicitely wrapped in `^$` (eg. `^LIMIT$`).

Example:

```yaml
inventory:
  - name: cluster-01
    groups: [group-a, group-b]
  - name: cluster-02
    groups: [group-b, group-c]
  - name: cluster-03
    groups: [group-c, group-d]
  - name: cluster-04
    groups: [group-x]
```

```yaml
playbook:
  - name: test
    groups: [group-a, group-c]
```

* calling kusible without `-l` would execute the playbook on clusters 01 to 03 (each cluster is either in group-a or in group-c)
* so would calling it with `-l group-.*`
* calling it with `-l group-x` would not execute anything at all (as cluster-04 is neither in group-a nor group-c)
* calling it with `-l group-b` would execute it on cluster-01 and cluster-02
* calling it with `-l group-c -l group-d` would execute it only on cluster-03