# Tasks: Test Suite

## Phase 1 — pkg/printers/printers_test.go (write now, before 0001)
- [ ] Add `captureOutput(f func()) string` helper using `os.Pipe()`
- [ ] `TestPrintSecrets_NoSensitiveDataInOutput`
  - Construct `*v1.SecretList` with known sentinel values in `Data`
  - Assert sentinel strings are absent from captured output
  - Assert data count ("2") is present in output
- [ ] `TestPrinter_EmptyList`
  - For each of: `*v1.PodList`, `*v1.SecretList`, `*v1.ConfigMapList`, `*appsv1.DeploymentList`
  - Assert captured output is empty string
- [ ] `TestPrinter_PatternFilter`
  - Construct `*v1.ConfigMapList` with "nginx-config" and "redis-config"
  - Call Printer with pattern="nginx"
  - Assert "nginx-config" present, "redis-config" absent
- [ ] After 0002: replace `captureOutput` with `bytes.Buffer`; add `TestMatchesPattern`

## Phase 2 — pkg/util/util_test.go (write after 0001)
- [ ] Add `fake` import: `k8s.io/client-go/kubernetes/fake`
- [ ] `TestGetter_CustomKinds`
  - `fake.NewSimpleClientset` with a pre-populated ConfigMapList
  - Call `Getter("default", fakeClient, "ConfigMaps", ch)`
  - Assert exactly one `*v1.ConfigMapList` received and channel closes
- [ ] `TestGetter_UnknownKind`
  - Call `Getter("default", fakeClient, "NonExistentKind", ch)`
  - Assert channel closes within 2 seconds (no deadlock)
- [ ] `TestGetter_ChannelAlwaysClosed`
  - Call `Getter("", fakeClient, "Pods", ch)`
  - Drain channel and assert it closes within 2 seconds

## Phase 3 — pkg/cache/cache_test.go (write with 0003)
- [ ] All cache tests use `t.TempDir()` — never touch `~/.kube`
- [ ] Expose internal `writeToDir(dir, key, meta, sections)` and
      `readFromDir(dir, key, ttl)` helpers for test use
- [ ] `TestKeyFor_Deterministic`
- [ ] `TestKeyFor_KindsSorted`
- [ ] `TestKeyFor_Unique`
- [ ] `TestReadWrite_RoundTrip`
- [ ] `TestRead_Missing`
- [ ] `TestRead_Expired` — write entry with `written_at = now - 2m`, TTL = 60s
- [ ] `TestWrite_Atomic` — two goroutines write same key concurrently; result is valid JSON
- [ ] `TestNoSensitiveData`
  - Write a SectionEntry with printer-formatted output containing no sentinel strings
  - Read back the raw JSON file bytes
  - Assert none of: `"super-secret-value"`, `"my-api-token"`, `"password"` present

## Phase 4 — Integration test scaffold (write after 0003)
- [ ] Create `test/integration/integration_test.go` with `//go:build integration`
- [ ] Skip if `kind` binary not on PATH
- [ ] `TestIntegration_BasicRun`: create kind cluster, run ksearch, assert non-empty output
- [ ] `TestIntegration_CacheHit`: run ksearch twice within TTL, assert second run uses cache

## Phase 5 — Verification
- [ ] `go test ./pkg/printers/...` — Phase 1 tests pass
- [ ] `go test ./pkg/util/...` — Phase 2 tests pass (requires 0001)
- [ ] `go test ./pkg/cache/...` — Phase 3 tests pass (requires 0003)
- [ ] `go test -race ./...` — no data races
- [ ] `go test -tags integration ./...` — integration tests pass (requires kind)
- [ ] `go test ./pkg/printers/... -run TestPrintSecrets_NoSensitiveDataInOutput` — security gate passes
- [ ] `go test ./pkg/cache/... -run TestNoSensitiveData` — cache security gate passes
