# Distribution Specification

## Purpose
Manages how ksearch is packaged, named, distributed, and discovered as a
kubectl plugin — both via the krew plugin index and via manual installation.

## Requirements

### Requirement: kubectl plugin discovery
The distributed binary SHALL remain named `ksearch`. The krew manifest SHALL
identify `ksearch` (or `ksearch.exe` on Windows) as the installed executable,
and krew SHALL provide the `kubectl-ksearch` symlink used for plugin discovery.

#### Scenario: Plugin discovered by kubectl
- GIVEN the plugin is installed by krew
- WHEN `kubectl plugin list` is run
- THEN `kubectl-ksearch` appears in the list via the symlink krew created

### Requirement: Krew index listing
The plugin SHALL be listed on the official `kubernetes-sigs/krew-index` so
users can install via `kubectl krew install ksearch`.

#### Scenario: Installation via krew
- GIVEN krew is installed
- WHEN `kubectl krew install ksearch` is run
- THEN the plugin is installed and `kubectl ksearch --help` works

### Requirement: Automated krew-index updates
Each tagged release SHALL automatically open a PR to update the krew-index
manifest via `krew-release-bot`, requiring no manual intervention.

#### Scenario: New release triggers krew update
- GIVEN a new `v*.*.*` tag is pushed
- WHEN GoReleaser publishes the GitHub Release
- THEN `krew-release-bot` opens a PR to `kubernetes-sigs/krew-index`
  updating the version, uri, and sha256 fields

### Requirement: Context-aware help text
The CLI help output SHALL show `kubectl ksearch` when invoked as a kubectl
plugin and show the bare binary name otherwise.

#### Scenario: Help text via kubectl
- GIVEN the binary is invoked as `kubectl ksearch`
- WHEN `--help` is passed
- THEN the usage line reads `kubectl ksearch [flags]`

#### Scenario: Help text via direct invocation
- GIVEN the binary is invoked directly as `./ksearch`
- WHEN `--help` is passed
- THEN the usage line reads `ksearch [flags]`

### Requirement: Cross-platform archives
Release archives SHALL be produced for linux/amd64, linux/arm64, darwin/amd64,
darwin/arm64, and windows/amd64. Each SHALL include the binary and a LICENSE file.

#### Scenario: Archive contents
- GIVEN a release archive for linux/amd64
- WHEN extracted
- THEN it contains `ksearch` and `LICENSE`
