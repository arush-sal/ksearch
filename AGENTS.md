# Repository Guidelines

## Project Structure & Module Organization

`ksearch` is a Go CLI kubectl plugin for comprehensive Kubernetes resource discovery.
Entry points are `main.go` (calls `cmd.Execute()`) and `cmd/ksearch.go` (Cobra flags,
client init, cache handling, discovery, result dispatch). Core resource fetching is in
`pkg/util/`; stdout rendering is in `pkg/printers/`; local cache logic is in
`pkg/cache/`. Dependency snapshots are checked into `vendor/`. Design proposals and
planned changes live under `openspec/`; treat those as implementation guidance, not
runtime code.

### Data flow

```text
cmd/ksearch.go  ──► pkg/util/util.go  ──► pkg/printers/printers.go
 (Cobra flags,       (clientset List       (type-switch → tabwriter
  client init,        calls over a          tables to stdout)
  result loop)        chan interface{})
```

The fetch path now supports both typed Kubernetes list structs (`*v1.PodList`,
`*appsv1.DeploymentList`, etc.) and dynamically discovered unstructured lists.
The type switch in `Printer()` must cover every typed result `Getter()` can emit,
and unstructured output must remain sensible for dynamically discovered resources.

### Planned upgrades (openspec/)

Pending changes live in `openspec/changes/`. Completed and merged changes are
archived in `openspec/changes/done/`. When a change is merged, move its directory
into `done/` and commit the move before starting the next change.

**Currently pending:** none

Full proposals, designs, and task checklists are in `openspec/changes/<id>/`.
Completed change specs are in `openspec/changes/done/<id>/` for reference.

---

## Build, Test, and Development Commands

Use the `Makefile` for common tasks, or standard Go tooling directly.

```bash
# Build
make build
go build -o ksearch .

# Run
./ksearch
./ksearch -n <namespace>
./ksearch -n <namespace> -p <name-pattern>
./ksearch -k configmap,secret -n default

# Test
make test
go test ./...                         # all packages
go test ./pkg/printers/...            # single package
go test -run TestMatchesPattern ./pkg/printers/...   # single test
go test -v -count=1 ./...             # verbose, no cache

# Vet and format
make vet
go vet ./...
make fmt
gofmt -l .                            # list files that need formatting
gofmt -w .                            # apply formatting

# Lint
make lint
```

The binary is `.gitignored`; rebuild after every source change.

---

## Coding Style & Naming Conventions

- Tabs for indentation; `gofmt` formatting enforced.
- Exported identifiers only when cross-package use requires them.
- Keep CLI wiring in `cmd/`, Kubernetes API fetch logic in `pkg/util/`, formatting
  logic in `pkg/printers/`, cache logic in `pkg/cache/`.
- Error handling: log via `logrus` and continue rather than `log.Fatal`/`panic`
  inside library packages; reserve `OrDie` helpers for `cmd/` init paths only.

---

## Testing Scope

Tests already exist across command, cache, printer, and util packages. New work
should follow and extend the existing test patterns below.

### pkg/printers

| Test                        | What to verify                                                                                                            |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `TestMatchesPattern`        | Empty pattern matches any name; non-empty pattern matches substring; non-matching returns false                           |
| `TestPrinter_<Kind>`        | Given a populated `*v1.XList`, `Printer()` writes expected header and rows to a `bytes.Buffer`; empty list writes nothing |
| `TestPrinter_PatternFilter` | Rows not matching the pattern are absent from output                                                                      |

Use `bytes.Buffer` as the `io.Writer` target.

### pkg/util

| Test                     | What to verify                                                                 |
|--------------------------|--------------------------------------------------------------------------------|
| `TestGetter_CustomKinds` | Passing `kinds="configmap"` restricts the resource list to ConfigMaps only     |
| `TestGetter_UnknownKind` | An unrecognised kind logs an error and closes the channel cleanly              |

Use `k8s.io/client-go/kubernetes/fake` to construct a fake clientset that returns
pre-populated lists without hitting a real cluster.

### pkg/cache

| Test                       | What to verify                                                       |
|----------------------------|----------------------------------------------------------------------|
| `TestKeyFor_Deterministic` | Same (context, namespace, kinds) always returns the same SHA-256 key |
| `TestKeyFor_Unique`        | Differing namespace or context produces a different key              |
| `TestReadWrite_RoundTrip`  | Data written with `Write()` is recovered intact by `Read()`          |
| `TestRead_Expired`         | `Read()` returns nil when `written_at + ttl < now`                   |
| `TestRead_Missing`         | `Read()` returns nil (not an error) when no file exists              |
| `TestWrite_Atomic`         | Resulting file is valid JSON even under concurrent `Write()` calls   |

### Running tests against a live cluster (integration)

Look for kind binary, if it is present then create a test cluster and Set `KUBECONFIG` to that test cluster context and run:

```bash
go test -tags integration ./...
```

Integration tests are gated behind the `integration` build tag so they are
skipped in CI unless explicitly enabled.

---

## CI/CD and Release Cycle

### Current CI pipeline (.github/workflows/)

#### `ci.yml` — runs on every push and pull request to `master`

```yaml
jobs:
  lint:    golangci-lint (staticcheck, errcheck, govet, gofmt)
  test:    go test -race -coverprofile=coverage.out ./...
  build:   go build -o ksearch . (matrix: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
```

Minimum Go version is pinned from `go.mod` via `actions/setup-go` with Go module cache enabled.

#### `release.yml` — runs on git tag push matching `v*.*.*`

```yaml
jobs:
  goreleaser:  uses GoReleaser to build cross-platform binaries and create a GitHub Release
```

Artefacts per release:

- `ksearch_linux_amd64.tar.gz`
- `ksearch_linux_arm64.tar.gz`
- `ksearch_darwin_amd64.tar.gz`
- `ksearch_darwin_arm64.tar.gz`
- `ksearch_windows_amd64.zip`
- SHA-256 checksums file
- kubectl krew plugin manifest (`ksearch.yaml`) — required for krew submission

### `.goreleaser.yml`

Key settings:

- `builds`: set `ldflags` to inject version string from git tag (`-X main.version={{.Version}}`)
- `archives`: include `LICENSE` in every archive
- `checksum`: sha256
- `krew`: generate the krew manifest pointing at the GitHub release assets

### Release cadence

| Branch / tag               | Purpose                                               |
|----------------------------|-------------------------------------------------------|
| `master`                   | Integration branch; must always build and pass CI     |
| `develop`                  | Feature development (matches existing branch history) |
| `v<major>.<minor>.<patch>` | Release tags; trigger the release pipeline            |

Versioning follows **semver**:

- Patch (`v0.0.x`): bug fixes, doc updates, dependency bumps
- Minor (`v0.x.0`): new resource kinds, new flags, non-breaking behaviour changes
- Major (`v1.0.0`): breaking CLI changes (flag renames, output format changes, removed kinds)

### PR checklist

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] New/changed behaviour covered by a test
- [ ] If CLI flags or output format changed: sample output included in PR description
- [ ] If a new resource kind added: both `util.go` and `printers.go` updated
- [ ] Linked to the relevant `openspec/changes/<id>/` doc if applicable

---

## Commit & Pull Request Guidelines

Prefer short imperative commit subjects: `Add cache package`, `Fix pattern filter for Secrets`.
Keep each commit scoped to one logical change. PR descriptions must list:

1. The user-visible behaviour change
2. Verification steps (commands run + expected output)
3. Link to the relevant `openspec/changes/<id>/` spec if applicable

---

## Environment & Safety Notes

This tool talks to a live Kubernetes cluster through the current kubeconfig context.
During development, test against a dedicated non-production namespace.
The application cache (`~/.kube/ksearch/cache/`) stores command output locally.
Be careful when validating cache behavior against live clusters.

## Agent Orchestrator (ao) Session

You are running inside an Agent Orchestrator managed workspace.
Session metadata is updated automatically via shell wrappers.

If automatic updates fail, you can manually update metadata:

```bash
~/.ao/bin/ao-metadata-helper.sh  # sourced automatically
# Then call: update_ao_metadata <key> <value>
```
