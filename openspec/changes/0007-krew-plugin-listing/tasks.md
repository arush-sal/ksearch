# Tasks: Krew Plugin Index Listing

## Phase 1 — kubectl-aware help text

Binary stays `ksearch`. Krew creates the `kubectl-ksearch` symlink automatically.
The `pluginName()` function detects the invocation context at runtime.

- [ ] `cmd/ksearch.go`: add `pluginName()` function:
  ```go
  func pluginName() string {
      base := filepath.Base(os.Args[0])
      base = strings.TrimSuffix(base, filepath.Ext(base))
      if strings.HasPrefix(base, "kubectl-") {
          return "kubectl " + strings.TrimPrefix(base, "kubectl-")
      }
      return base
  }
  ```
- [ ] `cmd/ksearch.go`: in `init()`, set `rootCmd.Use = pluginName()` (remove static `"ksearch"` from `rootCmd` declaration)
- [ ] Add `"path/filepath"` import to `cmd/ksearch.go` if not already present
- [ ] Verify: `go build -o ksearch . && ./ksearch --help` shows `ksearch` as usage
- [ ] Verify: `cp ksearch kubectl-ksearch && ./kubectl-ksearch --help` shows `kubectl ksearch` as usage

## Phase 2 — GoReleaser krew config refinements

- [ ] `.goreleaser.yml`: update `krews` section:
  - `short_description`: `"Search and list all Kubernetes resources across API groups"`
  - `description`: multi-line description covering core/v1, apps/v1, pattern matching, caching
  - Add `caveats` field: kubeconfig requirement, cache behaviour note
- [ ] `.goreleaser.yml`: verify `archives.files` includes `LICENSE`
- [ ] `.goreleaser.yml`: verify `krews.repository` points to `arush-sal/krew-index`
  (GoReleaser generates the manifest template here; the bot submits to the official index)
- [ ] Verify: `goreleaser release --snapshot --clean` succeeds and `dist/krew/ksearch.yaml`
  is generated with correct `bin: ksearch` on non-windows and `bin: ksearch.exe` on windows

## Phase 3 — krew-release-bot GitHub Action

- [ ] Create `.krew.yaml` with platform entries for linux/amd64, linux/arm64,
  darwin/amd64, darwin/arm64, and windows/amd64 using `bin: ksearch`
  (`ksearch.exe` on Windows) and release asset URLs under
  `https://github.com/arush-sal/ksearch/releases/download/{{ .TagName }}/...`
- [ ] Create `.github/workflows/krew-release.yml`:
  ```yaml
  name: Update krew-index

  on:
    release:
      types: [published]

  jobs:
    krew-update:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - uses: rajatjindal/krew-release-bot@v0.0.46
  ```
- [ ] Verify: the existing `release.yml` triggers on `push: tags: ['v*.*.*']` which
  creates a GitHub Release → krew-release.yml triggers on `release: [published]`

## Phase 4 — Local krew validation

- [ ] Build snapshot: `goreleaser release --snapshot --clean`
- [ ] Inspect generated manifest: `cat dist/krew/ksearch.yaml`
  - Verify `apiVersion: krew.googlecontainertools.github.com/v1alpha2`
  - Verify `metadata.name: ksearch`
  - Verify all platform entries have `uri`, `sha256`, `bin`
  - Verify `bin: ksearch` (non-windows) and `bin: ksearch.exe` (windows)
- [ ] Test local install (requires krew installed):
  ```bash
  kubectl krew install --manifest=dist/krew/ksearch.yaml \
    --archive=dist/ksearch_linux_amd64.tar.gz
  kubectl ksearch --help
  kubectl ksearch --version
  kubectl krew uninstall ksearch
  ```
- [ ] Cross-platform validation:
  ```bash
  KREW_OS=darwin KREW_ARCH=arm64 kubectl krew install \
    --manifest=dist/krew/ksearch.yaml \
    --archive=dist/ksearch_darwin_arm64.tar.gz
  kubectl krew uninstall ksearch
  ```

## Phase 5 — Tests

- [ ] `cmd/ksearch_test.go`: add `TestPluginName_WithKubectlPrefix` — set
  `os.Args[0] = "/usr/local/bin/kubectl-ksearch"`, assert `pluginName()` returns
  `"kubectl ksearch"`
- [ ] `cmd/ksearch_test.go`: add `TestPluginName_WithoutPrefix` — set
  `os.Args[0] = "/usr/local/bin/ksearch"`, assert `pluginName()` returns `"ksearch"`
- [ ] `cmd/ksearch_test.go`: add `TestPluginName_WindowsExe` — set
  `os.Args[0] = "C:\\Users\\foo\\kubectl-ksearch.exe"`, assert `pluginName()` returns
  `"kubectl ksearch"`

## Phase 6 — Submit to krew-index (manual, post-first-release)

- [ ] Tag and push first release with updated config: `git tag v0.x.x && git push origin v0.x.x`
- [ ] Wait for GoReleaser to complete and GitHub Release to be published
- [ ] Download the generated manifest from the release assets (or copy from `dist/krew/ksearch.yaml`)
- [ ] Fork `kubernetes-sigs/krew-index`
- [ ] Copy manifest to `plugins/ksearch.yaml` in the fork
- [ ] Open PR to `kubernetes-sigs/krew-index`
- [ ] After merge: verify `kubectl krew install ksearch` works

## Phase 7 — Verification

- [ ] `go build -o ksearch .` — zero errors
- [ ] `go vet ./...` — zero warnings
- [ ] `go test -race ./...` — all tests pass
- [ ] `./ksearch --help` — shows `ksearch` as usage prefix (standalone mode)
- [ ] `cp ksearch kubectl-ksearch && ./kubectl-ksearch --help` — shows `kubectl ksearch` (plugin mode)
- [ ] `./ksearch --version` — prints injected version
- [ ] `goreleaser release --snapshot --clean` — succeeds, manifest generated
- [ ] `cat dist/krew/ksearch.yaml` — `bin: ksearch` (not `kubectl-ksearch`)
- [ ] `kubectl krew install --manifest=dist/krew/ksearch.yaml --archive=dist/ksearch_linux_amd64.tar.gz` — installs cleanly
- [ ] `kubectl ksearch --help` — shows `kubectl ksearch` (krew symlink detected)
- [ ] `kubectl ksearch -n kube-system` — produces output
- [ ] `kubectl krew uninstall ksearch` — removes cleanly
