# CI/CD Specification

## Purpose
Defines the automated build, test, and release pipeline for ksearch.
Every pull request to `master` SHALL be validated by CI. Every semver
tag SHALL produce a GitHub Release with cross-platform binaries and a
kubectl krew manifest.

## Requirements

### Requirement: PR validation pipeline
Every push to a pull request targeting `master` SHALL trigger a pipeline
that lints, vets, tests, and builds the project.

#### Scenario: Failing test blocks merge
- GIVEN a PR where `go test ./...` returns non-zero
- WHEN the CI pipeline runs
- THEN the pipeline reports failure and the PR cannot be merged

#### Scenario: Build matrix covers all target platforms
- GIVEN a passing PR
- WHEN the build job runs
- THEN binaries are produced for linux/amd64, linux/arm64,
  darwin/amd64, darwin/arm64, and windows/amd64

### Requirement: Linting enforced in CI
`golangci-lint` SHALL run on every PR with at minimum the
`staticcheck`, `errcheck`, `govet`, and `gofmt` linters enabled.

#### Scenario: gofmt violation blocks merge
- GIVEN a PR with unformatted Go source
- WHEN the lint job runs
- THEN the pipeline reports failure

### Requirement: Release pipeline on semver tag
Pushing a tag matching `v*.*.*` to `master` SHALL trigger a release
pipeline that produces cross-platform archives and a GitHub Release.

#### Scenario: Release artefacts produced
- GIVEN a tag `v1.1.0` is pushed
- WHEN the release pipeline runs
- THEN the GitHub Release contains:
  - `ksearch_linux_amd64.tar.gz`
  - `ksearch_linux_arm64.tar.gz`
  - `ksearch_darwin_amd64.tar.gz`
  - `ksearch_darwin_arm64.tar.gz`
  - `ksearch_windows_amd64.zip`
  - `checksums.txt` (SHA-256)
  - `ksearch.yaml` (krew plugin manifest)

### Requirement: Version injected at build time
The released binary SHALL embed the git tag as its version string,
reported by `ksearch --version`.

#### Scenario: Version matches tag
- GIVEN the binary was built from tag `v1.1.0`
- WHEN `ksearch --version` is run
- THEN the output contains `v1.1.0`

### Requirement: Semver release cadence
Releases SHALL follow semantic versioning:
- Patch (`v0.0.x`): bug fixes, dependency bumps, doc updates
- Minor (`v0.x.0`): new resource kinds, new flags, non-breaking changes
- Major (`v1.0.0`): breaking CLI changes (flag renames, removed kinds,
  output format changes)

### Requirement: Go module cache in CI
CI workflows SHALL cache Go module downloads and the vendor directory
to reduce pipeline duration.

#### Scenario: Second run uses cached modules
- GIVEN the vendor directory and module cache are unchanged
- WHEN CI runs again on the same dependency set
- THEN no network downloads occur for Go modules
