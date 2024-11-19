docker_registry_secret:
  url: https://registry.cloudogu.com
  username: ""
  password: ""
dogu_registry_secret:
  url: https://dogu.cloudogu.com/api/v2/dogus
  urlschema: "https://dogu.cloudogu.com/api/v2/dogus_SCHEMA"
  username: ""
  password: ""
helm_registry_secret:
  host: k3ces.local:30098
  schema: oci
  plainHttp: "true"
  username: ""
  password: ""
components:
  k8s-longhorn:
    version: latest
    helmRepositoryNamespace: k8s
    deployNamespace: longhorn-system
    valuesYamlOverwrite: |
      longhorn:
        defaultSettings:
          storageOverProvisioningPercentage: 1000
        persistence:
          defaultClassReplicaCount: 2
        csi:
          attacherReplicaCount: 2
          provisionerReplicaCount: 2
          resizerReplicaCount: 2
          snapshotterReplicaCount: 2
        longhornUI:
          # Scale this up, if UI is needed
          replicas: 0