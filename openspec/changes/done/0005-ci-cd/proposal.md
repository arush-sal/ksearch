# Proposal: CI/CD Pipeline and Release Cycle

## Purpose
There is currently no automated validation on pull requests and no release
pipeline. This change introduces a GitHub Actions CI workflow for every PR
and a GoReleaser-based release pipeline triggered by semver tags, producing
cross-platform binaries and a kubectl krew manifest.

## Requirements

### Requirement: Automated PR validation
Every pull request to `master` SHALL be automatically built, linted, vetted,
and tested before it can be merged.

#### Scenario: Broken PR is blocked
- GIVEN a PR that fails `go test ./...`
- WHEN the CI workflow runs
- THEN the workflow reports failure and merge is blocked

### Requirement: Cross-platform release binaries
Every semver tag SHALL produce a GitHub Release containing binaries for
linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64,
plus a SHA-256 checksums file.

#### Scenario: Tag triggers release
- GIVEN tag `v1.1.0` is pushed
- WHEN the release workflow runs
- THEN all five platform archives and checksums.txt are attached to the
  GitHub Release

### Requirement: Version embedded in binary
The released binary SHALL report the git tag via `ksearch --version`.

#### Scenario: Version string matches tag
- GIVEN binary built from tag `v1.1.0`
- WHEN `ksearch --version` is run
- THEN output contains `v1.1.0`

### Requirement: Krew manifest generated on release
Each release SHALL produce a `ksearch.yaml` krew plugin manifest pointing
at the release archive URLs and their SHA-256 checksums.

#### Scenario: Krew manifest attached to release
- GIVEN a successful release pipeline run
- WHEN the GitHub Release assets are listed
- THEN `ksearch.yaml` is present and references the correct archive URLs
