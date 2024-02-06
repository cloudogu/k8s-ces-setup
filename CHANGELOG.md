# k8s-ces-setup Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.21.0] - 2024-02-06
### Changed
- **Breaking:** [#86] Passwords (Docker-, Dogu- & Helmregistry) have to be encoded in Base64 
(see [here](docs/operations/installation_guide_en.md) or [here](docs/operations/configuration_guide_en.md))

## [v0.20.2] - 2024-01-09
### Fixed
- [#84] Use default value for the urlschema in the dogu registry secret.

## [v0.20.1] - 2023-12-13
### Changed
- [#82] Update component-operator dependency to 0.7.0.

### Fixed
- [#82] Fix issues with helm template.

## [v0.20.0] - 2023-12-08
### Added
- [#80] Add component patch template file for mirroring this chart in offline environments.

## [v0.19.1] - 2023-11-22
### Fixed
- [#78] Remove timeout and wait indefinitely for components to get "ready".

## [v0.19.0] - 2023-11-16
### Added
- [#76] components can overwrite their values.yaml-default values

## [v0.18.0] - 2023-10-19
### Added
- [#74] Add functionality to install the component `k8s-cert-manager` before all other operators.

## [v0.17.1] - 2023-10-11
### Changed
- [#72] Update component-operator
- Update other dependencies
- Replace go-yaml with sigs.k8s.io/yaml

## [v0.17.0] - 2023-10-06
### Added
- [#70] Add struct for the components to specify attributes like the deployNamespace.
  - With this change it is possible to install longhorn as a component.

## [v0.16.2] - 2023-10-05
### Changed
- [#68] Change component setup to install CRDs separately

## [v0.16.1] - 2023-09-04
### Changed
- [#66] Use new helm registry config from the component-operator where the url is divided in host and schema.

### Changed
- [#62] Use `Info` as default log level.
- [#64] Match Makefile helm variable with those from a newer Makefile version 

### Added
- [#59] Add helm chart as release artifact.

## [v0.16.0] - 2023-08-14
### Added
- [#56] Allow to configure resource patches, a powerful way to modify Kubernetes resources during the setup process
  - please see the [docs](docs/operations/configuration_guide_en.md) for more information

### Changed
- Allows to configure the IP address placeholder `<<ip>>` in the `setup.json` section `naming/fqdn` as described in the official [setup docs](https://docs.cloudogu.com/de/docs/system-components/ces-setup/operations/setup-json/#fqdn)
- [#52] Use latest etcd release from dogu registry.

### Fixed
- Uses now singular context object for all Kubernetes requests

## [v0.15.0] - 2023-06-13
### Added 
- [#54] Use IP address as FQDN from load-balancer if it is missing
  - With this change, we are improving the development on external cloud providers by identifying the FQDN early on.

## [v0.14.0] - 2023-05-16
### Fixed
- [#50] Reduce technical debt

### Changed
- [#48] Deploy the etcd client as deployment instead of stateful set.

## [v0.13.2] - 2023-04-14
### Fixed
- [#46] Trim "dogus/" suffix only on URL "default" schema
  - this change avoids removing the endpoint suffix for the "index" schema

## [v0.13.1] - 2023-03-29
### Changed
- [#44] Improve logging in wait for pod step. API error doesn't throw an error now so that the wait functionality
will be canceled by the timeout.

## [v0.13.0] - 2023-03-24
### Removed
- [#41] Remove SSL API which generated selfsigned certificates. The API is made available in [`k8s-service-discovery`](https://github.com/cloudogu/k8s-service-discovery).

## [v0.12.0] - 2023-02-17
### Added
- Add optional volume mount for selfsigned cert of the dogu registry; #38

## [v0.11.1] - 2023-01-31
### Fixed
- [#36] Fixed an issue where the finisher cronjob starts infinite jobs if the pod e.g. can't pull an image.

### Changed
- Update makefiles to version 7.2.0
- Update `ces-build-lib` to 1.62.0

## [v0.11.0] - 2023-01-13
### Changed
- [#34] Add/Update label for consistent mass deletion of CES K8s resources
  - Select any k8s-ces-setup related resources like this: `kubectl get deploy,pod,... -l app=ces,app.kubernetes.io/name=k8s-ces-setup`
  - Select all CES components like this: `kubectl get deploy,pod,... -l app=ces`
  - Update `ces-build-lib` to 1.61.0

### Fixed
- [#32] Fixed a permission issue where the setup finisher cronjob was not allowed to execute his finisher script.

## [v0.10.0] - 2022-12-05
### Fixed
- [#30] The `ecosystem-certificate` TLS secret will now be created during setup.

## [v0.9.0] - 2022-11-30
### Fixed
- [#28] Setup wrongly assumed that all service accounts are of type dogu when creating step to wait for 
  them. Now only steps for dogu service accounts are created. 

## [v0.8.1] - 2022-11-23
### Fixed
- [#26] Use correct label for dogu resources
  - `dogu.name=name` is now valid

### Changed
- [#24] Read dogu registry URL schema from cluster secret instead of config.

## [v0.8.0] - 2022-08-30
### Changed
- [#22] If the resource urls from the k8s-components e.g. `dogu-operator` have the same host as the configured
  dogu registry, basic auth will be used for those components.
- [#22] Update `makefiles` to version 7.0.1
- [#22] Update `ces-build-lib` to version 1.56.0

## [v0.7.1] - 2022-08-30
### Fixed
- Internal release bugfix

## [v0.7.0] - 2022-08-30
### Changed
- [#20] Update internally used dependency versions
  - Update `cesapp-lib` to version v0.4.0
  - Update `k8s-apply-lib` to version v0.4.0
  - Update `k8s-dogu-operator` to version v0.11.0

## [v0.6.0] - 2022-06-13
### Changed
- [#17] Update makefiles to version 6.0.3

### Fixed
- [#17] Change order of certificate chain und use unique serial number in generation.

### Changed
- Extract client to apply k8s resources into own repository 
(https://github.com/cloudogu/k8s-apply-lib)

## [v0.5.0] - 2022-05-24
### Added
- [#12] Implement the registryConfigEncrypted section from the `setup.json`. Setup creates secrets for these values which 
- can be processed by the `k8s-dogu-operator`.

## [v0.4.0] - 2022-05-12
### Added
- [#10] Automatic setup process with `setup.json`. See [custom setup configuration](docs/operations/custom_setup_configuration_en.md) for more information.

## [v0.3.0] - 2022-05-03
### Added
- [#8] Setup installs `k8s-service-discovery` when performing a setup. Please see
  the [Configuration](docs/operations/configuration_guide_en.md) for more information

### Changed
- Update makefiles to version 5.1.0

## [v0.2.0] - 2022-04-28
### Added
- Setup installs vital Cloudogu EcoSystem (CES) K8s components to prepare namespace for setup and regular operation:
    - `etcd` (along with a development client)
    - the most important `k8s-dogu-operator` along with its own resources
    - Please see the [Configuration](docs/operations/configuration_guide_en.md)
      and [Installation Guides](docs/operations/installation_guide_en.md) for more information

### Changed
- Harmonize names of CES instance credential secrets with those of the
  consuming [Dogu Operator](https://github.com/cloudogu/k8s-dogu-operator)
- Development goodies
    - Make target for deployment to local cluster are more convienient
    - serve local resources with a simple HTTP server

## [v0.1.0] - 2022-03-11
### Added
- initial release of the basic setup skeleton
