# Tasks: Application-Level Caching

## Phase 1 — New pkg/cache package
- [ ] Create `pkg/cache/cache.go`
  - [ ] Define `SectionEntry`, `CacheMeta`, and `CacheEntry` structs with JSON tags
        — NO imports of `k8s.io/api` or any Kubernetes types
  - [ ] Implement `KeyFor(context, namespace, kinds, pattern string) string`
        using SHA-256 with null-byte separators; sort kinds before hashing;
        include pattern as the fourth component
  - [ ] Implement `Read(key string, ttl time.Duration) (*CacheEntry, error)`
        — return nil if file missing or `written_at + ttl < now`
  - [ ] Implement `Write(key string, meta CacheMeta, sections []SectionEntry) error`
        — base64-encode each section's output bytes, write to temp file,
        os.Rename for atomicity
  - [ ] Implement `cacheDir() string` helper returning `~/.kube/ksearch/cache/`
  - [ ] Call `os.MkdirAll(cacheDir(), 0700)` at the top of `Write()`
  - [ ] Set written cache files to mode 0600
- [ ] Create `pkg/cache/cache_test.go`
  - [ ] `TestKeyFor_Deterministic`: same inputs → same key
  - [ ] `TestKeyFor_KindsSorted`: `"secret,configmap"` and `"configmap,secret"` produce the same key
  - [ ] `TestKeyFor_Unique`: differing context, namespace, kinds, or pattern each produce a distinct key
  - [ ] `TestRead_Missing`: returns nil (not error) when file absent
  - [ ] `TestRead_Expired`: returns nil when `written_at + ttl < now`
  - [ ] `TestReadWrite_RoundTrip`: data written with `Write()` is recovered intact by `Read()`
  - [ ] `TestNoSensitiveData`: write a SectionEntry containing only printer-formatted text,
        read back raw JSON from disk, assert no known sensitive strings are present
  - [ ] `TestWrite_Atomic`: concurrent `Write()` calls to the same key leave a valid JSON file

## Phase 2 — Kubeconfig context detection in cmd/ksearch.go
- [ ] Import `k8s.io/client-go/tools/clientcmd`
- [ ] Load `clientcmd.NewDefaultClientConfigLoadingRules().Load()` to get `currentContext`
- [ ] Pass `currentContext` to `cache.KeyFor()`

## Phase 3 — Integration in cmd/ksearch.go
- [ ] Add `--cache-ttl` flag (type `time.Duration`, default `60s`) to `init()`
- [ ] Add `--no-cache` flag (type `bool`, default `false`) to `init()`
- [ ] Read `KSEARCH_CACHE_TTL` env var as fallback when `--cache-ttl` is not set
- [ ] Include `resName` (pattern) in the cache key derivation call
- [ ] Before calling `util.Getter()`:
  - [ ] Derive cache key via `cache.KeyFor(currentContext, namespace, kinds, resName)`
  - [ ] If `!noCache`: call `cache.Read(key, ttl)`
    - On hit: range `entry.Sections`, base64-decode each `Output`, write to `os.Stdout`; return early
  - [ ] On miss or `--no-cache`: proceed to `util.Getter()`
- [ ] After the fan-out Printer loop (change 0002) completes:
  - [ ] Collect `[]SectionEntry{Kind, base64(renderedBytes)}` from each goroutine's result
  - [ ] Call `cache.Write(key, meta, sections)` before flushing to stdout

## Phase 4 — Verification
- [ ] `grep -r "v1\." pkg/cache/` — must return zero results
- [ ] `grep -r "\.Data" pkg/cache/` — must return zero results
- [ ] `go test ./pkg/cache/... -run TestNoSensitiveData` — must pass
- [ ] `go test ./pkg/cache/...` — all unit tests pass
- [ ] `go build ./...` — zero errors
- [ ] `go vet ./...` — zero warnings
- [ ] Integration test: run `ksearch -n default` twice within TTL;
      second run must not make API calls (verify via `-v` log output)
- [ ] `ksearch -n default --no-cache` — hits API, overwrites cache file
- [ ] Switch kubeconfig context; verify a new distinct cache file is created
- [ ] `ksearch -n default -p nginx` and `ksearch -n default -p redis` produce
      different cache files (pattern is part of the key)
- [ ] Inspect cache file: `cat ~/.kube/ksearch/cache/*.json | jq '.sections[].output' | base64 -d`
      — output contains only table text, never secret values
- [ ] Verify file permissions: `stat ~/.kube/ksearch/cache/*.json` shows 0600
- [ ] Verify directory permissions: `stat ~/.kube/ksearch/cache/` shows 0700
- [ ] Run `ksearch -n default --cache-ttl=5s`, wait 6 seconds, run again;
      second run must miss the cache
