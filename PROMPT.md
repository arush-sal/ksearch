# Agent Task: ksearch Full Upgrade

## Project overview

You are working on **ksearch** — a Go kubectl plugin that lists and searches
Kubernetes resources across `core/v1` and `apps/v1` API groups. The repository
is at the root of your working directory.

Before writing any code, read these files in full:

- `AGENTS.md` — coding conventions, build commands, testing strategy, CI/CD spec
- `CLAUDE.md` — architecture overview and data flow
- `openspec/` — all proposals, designs, and task checklists for every change

The data flow is strictly:

```
cmd/ksearch.go → pkg/util/util.go → pkg/printers/printers.go
```

The channel between `Getter()` and the print loop carries `interface{}` holding
typed Kubernetes list structs. The type switch in `Printer()` must cover every
type `Getter()` can emit. Adding a resource kind always requires changes in both
`pkg/util/util.go` AND `pkg/printers/printers.go`.

---

## Branch setup

All work MUST be based off the `develop` branch. Before writing any code:

```bash
git fetch origin
git checkout develop
git pull origin develop
```

Create a dedicated branch off `develop` for each change, use worktree for the same:

```bash
git checkout -b <id>-<short-description>
# or
git worktree add -b <id>-<short-description> <path>
```

When a change is complete and its gate passes, merge it back into `develop`
before branching for the next change. Do not base any work off `master` directly.

---

## openspec/changes/ convention

- **Pending changes** — directories directly under `openspec/changes/` (excluding `done/`)
- **Completed changes** — moved into `openspec/changes/done/` after PR merge

When a change Pull Request on GitHub is merged, run:

```bash
mv openspec/changes/<id>-<name> openspec/changes/done/
```

Then commit the move as a standalone commit before starting the next change.
Only directories directly under `openspec/changes/` (not inside `done/`) represent
work still to be done.

---

## Your mission

Implement all pending changes in strict order. Do not start a later change until
the earlier one builds, vets, and tests cleanly. Each change has a complete spec
in `openspec/changes/<id>/` — read `proposal.md`, `design.md`, and `tasks.md`
before touching any code for that change. Use the task checklists in `tasks.md`
to track your progress.

**Currently pending (in order):**

None — all planned changes are complete.

Changes 0001–0007 are complete and archived in `openspec/changes/done/`.

---

## Change 0006 — Dynamic Resource Discovery (COMPLETED)

**Spec:** `openspec/changes/done/0006-dynamic-resource-discovery/`

Replace the duplicate static resource lists in `pkg/util/util.go` and `cmd/ksearch.go`
with live discovery from the Kubernetes API server.

### What to do

1. **Create `pkg/util/discover.go`:**
   - Define `ResourceMeta` struct: `Kind`, `Resource`, `APIGroup`, `APIVersion string`, `Namespaced bool`
   - Add `Discover(dc discovery.DiscoveryInterface, kinds string) ([]ResourceMeta, error)`:
     - Call `dc.ServerGroupsAndResources()`
     - On partial failure (err != nil, lists != nil): log warning, continue
     - Filter to resources with `"list"` in `Verbs`
     - If `kinds` non-empty: filter by kind name (case-insensitive)

2. **Update `pkg/util/util.go`:**
   - Delete `var resources = []string{...}` and `configuredResources()`
   - Update `Getter()` signature to `Getter(namespace string, clientset kubernetes.Interface, resources []ResourceMeta, c chan interface{})`
   - Dispatch on `meta.Kind`; normalise singular/plural (`"Pod", "Pods"`)
   - Replace `log.Error + return` on unknown kind with `log.Debugf + continue`

3. **Update `cmd/ksearch.go`:**
   - Delete `var defaultResources` and `effectiveResources()`
   - After building `clientset`, call `util.Discover(clientset.Discovery(), kinds)`
   - Build `resourceOrder` from the returned `[]ResourceMeta`
   - Pass `resources` to `util.Getter()`

4. **Add `pkg/util/discover_test.go`:**
   - `TestDiscover_AllWhenEmpty`, `TestDiscover_FilterByKinds`, `TestDiscover_SkipsNonListable`,
     `TestDiscover_CaseInsensitive`, `TestDiscover_MultipleKinds`

5. **Update `pkg/util/util_test.go`:**
   - Adapt all `TestGetter_*` to pass `[]ResourceMeta` instead of a `kinds string`

### Gate before proceeding

```bash
go build ./...
go vet ./...
go test -race ./pkg/util/...
grep -rn "defaultResources" . --include="*.go"    # zero results
grep -rn "effectiveResources" . --include="*.go"  # zero results
grep -rn "configuredResources" . --include="*.go" # zero results
```

---

## Change 0007 — Krew Plugin Index Listing (COMPLETED)

**Spec:** `openspec/changes/done/0007-krew-plugin-listing/`

Make ksearch a first-class krew plugin listed on the official
`kubernetes-sigs/krew-index`.

### What to do

1. **Add kubectl-aware help text in `cmd/ksearch.go`:**
   - Binary stays `ksearch` — krew creates the `kubectl-ksearch` symlink automatically
   - Add `pluginName()` that detects `kubectl-` prefix in `os.Args[0]`
   - Set `rootCmd.Use = pluginName()` in `init()`
   - Result: `kubectl ksearch --help` shows `kubectl ksearch [flags]`;
     `ksearch --help` shows `ksearch [flags]`

3. **Refine `.goreleaser.yml` krew config:**
   - Improve `short_description` and `description`
   - Add `caveats` field (kubeconfig requirement, cache note)

4. **Add krew-release-bot GitHub Action:**
   - Create `.github/workflows/krew-release.yml`
   - Triggers on `release: [published]`
   - Uses `rajatjindal/krew-release-bot@v0.0.46`
   - No secrets needed — bot uses its own credentials

5. **Local krew validation:**
   - `goreleaser release --snapshot --clean`
   - `kubectl krew install --manifest=dist/krew/ksearch.yaml --archive=dist/ksearch_linux_amd64.tar.gz`
   - Verify `kubectl ksearch --help` and `kubectl ksearch -n default`

6. **Tests:**
   - `TestPluginName_WithKubectlPrefix`, `TestPluginName_WithoutPrefix`,
     `TestPluginName_WindowsExe`

7. **Submit to `kubernetes-sigs/krew-index`** (manual, after first tagged release)

### Gate before proceeding

```bash
go build -o ksearch .
go vet ./...
go test -race ./...
./ksearch --help               # usage shows "ksearch" (standalone)
goreleaser release --snapshot --clean
cat dist/krew/ksearch.yaml     # manifest has bin: ksearch
kubectl krew install --manifest=dist/krew/ksearch.yaml \
  --archive=dist/ksearch_linux_amd64.tar.gz
kubectl ksearch --help         # usage shows "kubectl ksearch" (via krew symlink)
kubectl krew uninstall ksearch
```

---

## Autonomous development loop

Each change MUST follow this loop before being considered complete:

1. **Dev agent** — implement the change per the openspec `tasks.md` checklist.
   Branch off `develop`: `git worktree add -b <id>-<short-description> <path>`

2. **QA agent** — after dev agent signals completion, run all of:

   ```bash
   go build ./...          # must be zero errors
   go vet ./...            # must be zero warnings
   go test -race ./...     # must be zero failures
   golangci-lint run ./... # must be zero lint errors
   ```

   Plus all security gate commands specific to the change (see each change's "Gate" section).
   If any step fails: report failures to dev agent → loop back to step 1.

3. **PR agent** — after QA passes:
   - Push branch to origin
   - Create GitHub PR targeting `develop` with title `[<id>] <short description>`
   - PR body MUST include:
     - Link to `openspec/changes/<id>/` docs
     - `go test -race ./...` summary
     - Security gate results (where applicable)
     - Brief summary of what changed and why

4. **Senior Staff Engineer (SSE) review agent** — after PR is created:
   - Read all changed files in full
   - Apply the SSE review checklist below
   - Respond: **Approved** or **Changes Requested** (with specific line-level issues)
   - If Changes Requested: dev agent addresses each issue → QA re-runs → force-push → re-review
   - If Approved: Merge the PR into `develop` via `gh pr merge --squash`

5. **Post-merge cleanup** — dev agent:
   - `git pull origin develop`
   - Verify the changes has been merged
   - Move the completed change: `mv openspec/changes/<id>-<name> openspec/changes/done/`
   - Commit the move: `git add openspec/changes/ && git commit -m "Archive <id>: move to done"`
   - Delete feature branch and the worktree
   - Confirm `go build ./... && go test -race ./...` still passes on `develop`

### SSE review checklist

**Architecture:**

- No logic duplicated across packages (flag immediately)
- No raw Kubernetes types in `pkg/cache` (`grep -r "v1\." pkg/cache/` must be empty)
- Data flow strictly `cmd/ → pkg/util → pkg/printers` — no reverse imports

**Security:**

- `TestPrintSecrets_NoSensitiveDataInOutput` passes
- `TestNoSensitiveData` passes
- `pkg/cache` never imports `k8s.io/api`

**Code quality:**

- All error paths handled at system boundaries; internal errors logged and continued
- No new global mutable state (except existing package-level vars in `cmd/`)
- `golangci-lint` clean
- Flag any function longer than ~40 lines that could be split without losing clarity
- Flag repeated logic (3+ similar code blocks) that could be collapsed into a shared helper

**Tests:**

- Every new exported function has at least one test
- Table-driven tests where there are 3+ similar cases
- No test writes outside `t.TempDir()`
- No test requires a live cluster (unless `//go:build integration`)

---

## Global rules (from AGENTS.md)

- **Build tool:** `go build -o ksearch .` — no Makefile
- **Formatting:** `gofmt -w .` before every commit
- **Error handling:** `logrus` + continue in library packages; `OrDie` only in `cmd/` init paths
- **Commits:** short imperative subjects, one logical change per commit
  (`Add matchesPattern helper`, `Fix Getter context arg`)
- **Every PR** must link to the relevant `openspec/changes/<id>/` doc
- **Tests must pass without a live cluster** — use `fake.NewSimpleClientset()` for util,
  `t.TempDir()` for cache, `bytes.Buffer` for printers
- **On change completion:** move `openspec/changes/<id>-<name>/` into `openspec/changes/done/`
  and commit the move before starting the next change
- **Security invariants that must hold at all times:**
  - `pkg/cache` never imports `k8s.io/api` — enforced by `grep -r "v1\." pkg/cache/`
  - `TestPrintSecrets_NoSensitiveDataInOutput` and `TestNoSensitiveData` must always pass

---

## Final acceptance checklist

```bash
go build -o ksearch .
go vet ./...
go test -race ./...
go test ./pkg/printers/... -run TestPrintSecrets_NoSensitiveDataInOutput
go test ./pkg/cache/...    -run TestNoSensitiveData
grep -r "v1\." pkg/cache/              # zero results
grep -r "\.Data" pkg/cache/            # zero results
grep -rn "defaultResources" . --include="*.go"    # zero results (0006)
grep -rn "effectiveResources" . --include="*.go"  # zero results (0006)
grep -rn "configuredResources" . --include="*.go" # zero results (0006)
golangci-lint run ./...
goreleaser release --snapshot --clean
./ksearch --help               # usage shows "ksearch" (standalone, 0007)
./ksearch --version            # must print injected version, not "dev"
./ksearch -n default           # must produce output against a real cluster
./ksearch -n default           # second run within TTL must hit cache
./ksearch -n default --no-cache  # must bypass cache and refresh
kubectl krew install --manifest=dist/krew/ksearch.yaml \
  --archive=dist/ksearch_linux_amd64.tar.gz    # krew install works (0007)
kubectl ksearch --help         # usage shows "kubectl ksearch" via krew symlink
kubectl ksearch --version      # works via krew
kubectl krew uninstall ksearch # clean removal
ls openspec/changes/           # only done/ and any new pending changes
```
