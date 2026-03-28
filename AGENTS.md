# Repository Guidelines

## Project Structure & Module Organization

`ksearch` is a Go CLI kubectl plugin for comprehensive Kubernetes resource discovery.
Entry points are `main.go` (calls `cmd.Execute()`) and `cmd/ksearch.go` (Cobra flags,
client init, result dispatch). Core resource fetching is in `pkg/util/util.go`;
stdout rendering is in `pkg/printers/printers.go`. Dependency snapshots are checked
into `vendor/`. Design proposals and planned changes live under `openspec/`; treat
those as implementation guidance, not runtime code.

### Data flow

```
cmd/ksearch.go  ──► pkg/util/util.go  ──► pkg/printers/printers.go
 (Cobra flags,       (clientset List       (type-switch → tabwriter
  client init,        calls over a          tables to stdout)
  result loop)        chan interface{})
```

The channel carries concrete typed Kubernetes list structs (`*v1.PodList`,
`*appsv1.DeploymentList`, etc.). The type switch in `Printer()` must cover every
type `Getter()` can emit. **Adding a resource kind requires changes in both
`pkg/util/util.go` and `pkg/printers/printers.go`.**

### Planned upgrades (openspec/)

Three changes are tracked and must be implemented in order:

| ID | Summary |
|----|---------|
| `0001-dependency-upgrade` | Go 1.13→1.22; k8s client-go v0.17→v0.32.x; add `context.Context` to all `.List()` calls; remove deprecated `ComponentStatuses` |
| `0002-concurrent-printing` | Refactor `printers.go` to write to `io.Writer`; extract `matchesPattern` helper; fan-out goroutines with ordered flush |
| `0003-application-caching` | New `pkg/cache` package; SHA-256 keyed disk cache under `~/.kube/ksearch/cache/`; `--cache-ttl` and `--no-cache` flags |

Full proposals, designs, and task checklists are in `openspec/changes/<id>/`.

---

## Build, Test, and Development Commands

No Makefile. Use standard Go tooling.

```bash
# Build
go build -o ksearch .

# Run
./ksearch
./ksearch -n <namespace>
./ksearch -n <namespace> -p <name-pattern>
./ksearch -k configmap,secret -n default

# Test
go test ./...                         # all packages
go test ./pkg/printers/...            # single package
go test -run TestMatchesPattern ./pkg/printers/...   # single test
go test -v -count=1 ./...             # verbose, no cache

# Vet and format
go vet ./...
gofmt -l .                            # list files that need formatting
gofmt -w .                            # apply formatting
```

The binary is `.gitignored`; rebuild after every source change.

---

## Coding Style & Naming Conventions

- Tabs for indentation; `gofmt` formatting enforced.
- Exported identifiers only when cross-package use requires them.
- Keep CLI wiring in `cmd/`, Kubernetes API fetch logic in `pkg/util/`, formatting
  logic in `pkg/printers/`, cache logic in `pkg/cache/` (planned).
- Error handling: log via `logrus` and continue rather than `log.Fatal`/`panic`
  inside library packages; reserve `OrDie` helpers for `cmd/` init paths only.

---

## Testing Scope

Tests do not yet exist. The following test strategy should be followed when
implementing new code or the planned changes above.

### pkg/printers

| Test | What to verify |
|------|---------------|
| `TestMatchesPattern` | Empty pattern matches any name; non-empty pattern matches substring; non-matching returns false |
| `TestPrinter_<Kind>` | Given a populated `*v1.XList`, `Printer()` writes expected header and rows to a `bytes.Buffer`; empty list writes nothing |
| `TestPrinter_PatternFilter` | Rows not matching the pattern are absent from output |

Use `bytes.Buffer` as the `io.Writer` target once the `0002` refactor lands.
Until then, capture stdout with `os.Pipe()`.

### pkg/util

| Test | What to verify |
|------|---------------|
| `TestGetter_CustomKinds` | Passing `kinds="configmap"` restricts the resource list to ConfigMaps only |
| `TestGetter_UnknownKind` | An unrecognised kind logs an error and closes the channel cleanly |

Use `k8s.io/client-go/kubernetes/fake` to construct a fake clientset that returns
pre-populated lists without hitting a real cluster.

### pkg/cache (planned, change 0003)

| Test | What to verify |
|------|---------------|
| `TestKeyFor_Deterministic` | Same (context, namespace, kinds) always returns the same SHA-256 key |
| `TestKeyFor_Unique` | Differing namespace or context produces a different key |
| `TestReadWrite_RoundTrip` | Data written with `Write()` is recovered intact by `Read()` |
| `TestRead_Expired` | `Read()` returns nil when `written_at + ttl < now` |
| `TestRead_Missing` | `Read()` returns nil (not an error) when no file exists |
| `TestWrite_Atomic` | Resulting file is valid JSON even under concurrent `Write()` calls |

### Running tests against a live cluster (integration)

Look for kind binary, if it is present then create a test cluster and Set `KUBECONFIG` to that test cluster context and run:

```bash
go test -tags integration ./...
```

Integration tests are gated behind the `integration` build tag so they are
skipped in CI unless explicitly enabled.

---

## CI/CD and Release Cycle

### Planned CI pipeline (.github/workflows/)

#### `ci.yml` — runs on every push and pull request to `master`

```
jobs:
  lint:    golangci-lint (staticcheck, errcheck, govet, gofmt)
  test:    go test -race -coverprofile=coverage.out ./...
  build:   go build -o ksearch . (matrix: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
```

Minimum Go version pinned to whatever `go.mod` declares.
Use `actions/setup-go` with `cache: true` and `actions/cache` for the vendor directory.

#### `release.yml` — runs on git tag push matching `v*.*.*`

```
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

### .goreleaser.yml (to be created)

Key settings:

- `builds`: set `ldflags` to inject version string from git tag (`-X main.version={{.Version}}`)
- `archives`: include `LICENSE` in every archive
- `checksum`: sha256
- `release`: auto-generate changelog from commit messages
- `krew`: generate the krew manifest pointing at the GitHub release assets

### Release cadence

| Branch / tag | Purpose |
|---|---|
| `master` | Integration branch; must always build and pass CI |
| `develop` | Feature development (matches existing branch history) |
| `v<major>.<minor>.<patch>` | Release tags; trigger the release pipeline |

Versioning follows **semver**:

- Patch (`v0.0.x`): bug fixes, doc updates, dependency bumps
- Minor (`v0.x.0`): new resource kinds, new flags, non-breaking behaviour changes
- Major (`v1.0.0`): breaking CLI changes (flag renames, output format changes, removed kinds)

### PR checklist (to be added as `.github/pull_request_template.md`)

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
The planned application cache (`~/.kube/ksearch/cache/`) stores raw API responses
including Secret values in plaintext at mode 0600; document this in the README before
the 0003 change ships.
