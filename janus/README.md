## Prerequisites

- Kubernetes 1.12+
- Helm 2.11+

## Installing the Chart

To install the chart with the release name `my-release`:

```console
$ helm install --name my-release janus
```

The command deploys Janus on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

The following table lists the configurable parameters of the Janus chart and their default values.

| Parameter                           | Description                                                   | Default                                                  |
|-------------------------------------|---------------------------------------------------------------|----------------------------------------------------------|
| `image.repository`                  | Janus Image name                                              | `quay.io/hellofresh/janus`                               |
| `image.tag`                         | Janus Image tag                                               | `latest`                                                 |
| `image.pullPolicy`                  | Janus pull policy                                             | `IfNotPresent`                                           |
| `image.pullSecrets`                 | Specify docker-registry secret names as an array              | `[]` (does not add image pull secrets to deployed pods)  |
| `nameOverride`                      | String to partially override janus.fullname template with a string (will prepend the release name) | `nil`               |
| `fullnameOverride`                  | String to fully override janus.fullname template with a string                                     | `nil`               |
| `allowEmptyPassword`                | Allow DB blank passwords                                      | `yes`                                                    |
| `deployment.replicaCount`           | Number of Janus pod replicas                                  | `2`                                                      |
| `deployment.minAvailable`           | Creates PDB is min available (must be less than replicaCount) | `1`                                                      |
| `deployment.valuesFrom`             | Add needed env vars from Kubernetes metadata                  | `POD_NAME`                                               |
| `deployment.labels`                 | Add custom labels to the deployment                           | `app: janus`                                             |
| `deployment.databaseDSN`            | Database connection string                                    | `mongodb://janus-database:27017/janus`                   |
| `service.type`                      | Kubernetes Service type                                       | `LoadBalancer`                                           |
| `service.name`                      | Override service name                                         | ``                                                       |
| `service.ports[0].protocol`         | Service HTTP protocol                                         | `TCP`                                                    |
| `service.ports[0].port`             | Service HTTP port                                             | `80`                                                     |
| `service.ports[0].targetPort`       | Service HTTP target container port                            | `8080`                                                   |
| `service.ports[0].name`             | Service HTTP port name                                        | `http`                                                   |
| `ingress.enabled`                   | Enable ingress controller resource                            | `false`                                                  |
| `ingress.annotations`               | Ingress annotations                                           | `[]`                                                     |
| `ingress.name`                      | Override ingress name                                         | ``                                                       |
| `ingress.hosts[0].name`             | Hostname to your Janus installation                           | `janus.local`                                            |
| `ingress.hosts[0].paths[0].port`    | Port to service                                               | `80`                                                     |
| `ingress.hosts[0].paths[0].path`    | Path within the url structure                                 | `/`                                                      |
| `resources`                         | CPU/Memory resource requests/limits                           | Default values of the cluster                            |
| `affinity`                          | Map of node/pod affinities                                    | `{}`                                                     |
