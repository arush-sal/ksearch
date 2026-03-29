# Design: CI/CD Pipeline and Release Cycle

## Files to create

```
.github/
  workflows/
    ci.yml       — runs on every PR push and push to master
    release.yml  — runs on tag push matching v*.*.*
  pull_request_template.md
.goreleaser.yml
```

## ci.yml

```yaml
name: CI
on:
  push:
    branches: [master, develop]
  pull_request:
    branches: [master]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: go test -race -coverprofile=coverage.out ./...
      - run: go vet ./...

  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux   goarch: amd64
          - goos: linux   goarch: arm64
          - goos: darwin  goarch: amd64
          - goos: darwin  goarch: arm64
          - goos: windows goarch: amd64
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} go build -o /dev/null .
```

## release.yml

```yaml
name: Release
on:
  push:
    tags: ['v*.*.*']

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## .goreleaser.yml

```yaml
version: 2

builds:
  - main: .
    binary: ksearch
    ldflags:
      - -s -w -X main.version={{.Version}}
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ignore:
      - goos: windows
        goarch: arm64

archives:
  - format: tar.gz
    name_template: "ksearch_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE

checksum:
  name_template: checksums.txt
  algorithm: sha256

release:
  github:
    owner: arush-sal
    name: ksearch
  draft: false
  prerelease: auto

krew:
  index:
    owner: arush-sal
    name: krew-index
  name: ksearch
  short_description: Search and list all Kubernetes resources in a namespace
  homepage: https://github.com/arush-sal/ksearch
  platforms:
    - selector:
        matchLabels:
          os: linux
          arch: amd64
    - selector:
        matchLabels:
          os: linux
          arch: arm64
    - selector:
        matchLabels:
          os: darwin
          arch: amd64
    - selector:
        matchLabels:
          os: darwin
          arch: arm64
    - selector:
        matchLabels:
          os: windows
          arch: amd64
```

## Version injection

`cmd/ksearch.go` declares the version as a package-level var:

```go
var version = "dev" // overridden at build time via ldflags
```

The Cobra root command uses it:

```go
var rootCmd = &cobra.Command{
    Version: version,
    ...
}
```

`ksearch --version` then prints the injected tag value.

## golangci-lint configuration (.golangci.yml)

```yaml
linters:
  enable:
    - staticcheck
    - errcheck
    - govet
    - gofmt
    - gosimple
    - ineffassign
    - unused

linters-settings:
  gofmt:
    simplify: true

issues:
  exclude-use-default: false
```

## pull_request_template.md

```markdown
## What changed
<!-- One sentence description of the user-visible behaviour change -->

## Verification
<!-- Commands run and expected output -->

## Checklist
- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] New/changed behaviour is covered by a test
- [ ] If CLI flags or output format changed: sample output included above
- [ ] If a new resource kind added: both `util.go` and `printers.go` updated
- [ ] Linked to `openspec/changes/<id>/` spec if applicable
```

## Release cadence

| Tag pattern | Meaning |
|-------------|---------|
| `v0.0.x` | Bug fixes, dep bumps, doc updates |
| `v0.x.0` | New resource kinds, new flags, non-breaking changes |
| `v1.0.0` | Breaking CLI changes (flag renames, output format changes) |

Branches:
- `master` — always buildable; direct target of PRs
- `develop` — feature integration (existing branch pattern in repo)
- Tags pushed from `master` only
