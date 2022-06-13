# k8s-ces-setup Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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