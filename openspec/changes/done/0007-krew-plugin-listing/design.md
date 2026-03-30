# Design: Krew Plugin Index Listing

## 1. Binary name: keep `ksearch` — krew handles the symlink

The built binary stays `ksearch`. No rename needed.

Krew's plugin manifest has a `bin` field that tells krew which file inside the
archive is the executable. Krew then creates a `kubectl-ksearch` symlink
pointing to it. This means:

- **Standalone users**: download the archive, extract `ksearch`, put it on PATH,
  run `ksearch` directly.
- **Krew users**: `kubectl krew install ksearch` extracts `ksearch` from the
  archive, creates a `kubectl-ksearch` symlink, and `kubectl ksearch` works.

### .goreleaser.yml — no binary change needed

```yaml
builds:
  - main: .
    binary: ksearch               # stays as-is
    ldflags:
      - -s -w -X main.version={{ .Version }}
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ignore:
      - goos: windows
        goarch: arm64
```

The krew manifest's `bin: ksearch` (or `bin: ksearch.exe` on Windows) tells krew
where to find the executable. Krew creates the `kubectl-ksearch` symlink itself.

### No changes needed to CLAUDE.md, AGENTS.md, or .gitignore

Build command remains `go build -o ksearch .`.

## 2. kubectl-aware help text

When kubectl invokes a plugin, `os.Args[0]` is the full path to the symlink
(e.g. `/home/user/.krew/bin/kubectl-ksearch`). We can detect this to produce
appropriate usage strings.

```go
// cmd/ksearch.go
func pluginName() string {
    base := filepath.Base(os.Args[0])
    base = strings.TrimSuffix(base, filepath.Ext(base)) // strip .exe
    if strings.HasPrefix(base, "kubectl-") {
        return "kubectl " + strings.TrimPrefix(base, "kubectl-")
    }
    return base
}
```

Set `rootCmd.Use` to the result of `pluginName()` in `init()`:

```go
func init() {
    rootCmd.Use = pluginName()
    // ... existing flag setup
}
```

This means:
- Via krew (`os.Args[0]` = `.../kubectl-ksearch`): help shows `kubectl ksearch [flags]`
- Standalone (`os.Args[0]` = `.../ksearch`): help shows `ksearch [flags]`

## 3. GoReleaser krew config refinements

```yaml
krews:
  - name: ksearch
    ids:
      - default
    homepage: https://github.com/arush-sal/ksearch
    short_description: "Search and list all Kubernetes resources across API groups"
    description: |
      ksearch lists and searches Kubernetes resources across both core/v1 and
      apps/v1 API groups, including resources that kubectl get omits by default
      (ConfigMaps, Secrets, Endpoints, etc.). Supports pattern matching,
      namespace scoping, kind filtering, and result caching.
    caveats: |
      Requires a valid kubeconfig context. Uses the current-context by default.
      Results are cached for 60s; use --no-cache to bypass.
    skip_upload: true
    repository:
      owner: arush-sal
      name: krew-index
```

**Note on target index:** GoReleaser's `krews` section publishes to the repo
specified in `repository`. For the official `kubernetes-sigs/krew-index`, the
initial PR is manual (fork + PR to the official repo). After acceptance,
`krew-release-bot` handles updates. `skip_upload: true` keeps GoReleaser focused
on generating the local validation manifest in `dist/krew/ksearch.yaml` without
trying to push to a personal krew-index repository from CI.

## 4. krew-release-bot GitHub Action and template

Add a root `.krew.yaml` template for the bot to render on each release:

```yaml
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: ksearch
spec:
  version: {{ .TagName }}
  homepage: https://github.com/arush-sal/ksearch
  shortDescription: Search Kubernetes resources across API groups
  platforms:
    - selector:
        matchLabels:
          os: linux
          arch: amd64
      {{addURIAndSha "https://github.com/arush-sal/ksearch/releases/download/{{ .TagName }}/ksearch_linux_amd64.tar.gz" .TagName }}
      bin: ksearch
```

Add `.github/workflows/krew-release.yml`:

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

This action:
- Triggers when a GitHub Release is published (which GoReleaser does on tag push)
- Reads `.krew.yaml` to build the manifest for the official index update PR
- Opens a PR to `kubernetes-sigs/krew-index` automatically
- Trivial version bumps (only version, uri, sha256 changed) are auto-merged

**No secrets are needed** — the bot uses its own GitHub App credentials.

## 5. Local validation workflow

Before submitting to krew-index, validate locally:

```bash
# 1. Build a snapshot release
goreleaser release --snapshot --clean

# 2. Find the generated manifest
cat dist/krew/ksearch.yaml

# 3. Create a local linux/amd64 test archive from the built binary
mkdir -p /tmp/ksearch-local-archive
cp dist/ksearch_linux_amd64_v1/ksearch /tmp/ksearch-local-archive/ksearch
cp LICENSE /tmp/ksearch-local-archive/LICENSE
tar -C /tmp/ksearch-local-archive -czf /tmp/ksearch_local_linux_amd64.tar.gz \
  ksearch LICENSE

# 4. Update the linux/amd64 sha256 in dist/krew/ksearch.yaml to match the
#    local test archive. Krew still validates the manifest checksum even when
#    --archive overrides the download URI.
sha256sum /tmp/ksearch_local_linux_amd64.tar.gz

# 5. Test installation with the local archive
kubectl krew install --manifest=dist/krew/ksearch.yaml \
  --archive=/tmp/ksearch_local_linux_amd64.tar.gz

# 6. Verify
kubectl ksearch --help
kubectl ksearch -n default

# 7. Cleanup
kubectl krew uninstall ksearch
```

For local install tests, the manifest checksum must match the custom archive.
`--archive` replaces the download source, but Krew still verifies the selected
platform entry's `sha256` from the manifest before unpacking.

Cross-platform validation can still use the generated release archives directly:

```bash
# 8. Cross-platform validation (simulate darwin/arm64 on linux)
KREW_OS=darwin KREW_ARCH=arm64 kubectl krew install \
  --manifest=dist/krew/ksearch.yaml \
  --archive=dist/ksearch_darwin_arm64.tar.gz

# 9. Cleanup
kubectl krew uninstall ksearch
```

## 6. Initial krew-index submission

After the first tagged release with the updated GoReleaser config:

1. Fork `kubernetes-sigs/krew-index`
2. Copy the generated `dist/krew/ksearch.yaml` to `plugins/ksearch.yaml` in the fork
3. Verify the manifest matches the krew spec:
   - `apiVersion: krew.googlecontainertools.github.com/v1alpha2`
   - `metadata.name: ksearch`
   - All platform entries have valid `uri`, `sha256`, `bin`
   - `bin` field is `ksearch` (or `ksearch.exe` for Windows)
4. Open PR to `kubernetes-sigs/krew-index`
5. After merge, `krew-release-bot` handles all future updates

## Files changed

| File | Action |
|------|--------|
| `.goreleaser.yml` | **Update** — krew descriptions, caveats (binary name unchanged) |
| `.krew.yaml` | **New** — template consumed by `krew-release-bot` |
| `cmd/ksearch.go` | **Update** — add `pluginName()`, set `rootCmd.Use` dynamically |
| `.github/workflows/krew-release.yml` | **New** — krew-release-bot action |

## Security impact

None. This change only affects packaging, naming, and distribution metadata.
No runtime behaviour changes beyond help text formatting.
