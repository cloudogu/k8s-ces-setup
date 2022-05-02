# k8s-ces-setup Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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