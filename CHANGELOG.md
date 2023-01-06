# k8s-ces-setup Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- [#34] add/update label for consistent mass deletion of CES K8s resources
  - select any k8s-ces-setup related resources like this: `kubectl get deploy,pod,... -l app=ces,app.kubernetes.io/name=k8s-ces-setup`
  - select all CES components like this: `kubectl get deploy,pod,... -l app=ces`

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