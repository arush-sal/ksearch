# Design: Test Suite

## Overview

Three test files, one per testable package. No external test framework —
stdlib `testing` only. `k8s.io/client-go/kubernetes/fake` is used for
Getter tests and requires 0001 (dep upgrade) to land first.

```
pkg/printers/printers_test.go   ← write now (pre-0001)
pkg/util/util_test.go           ← write after 0001
pkg/cache/cache_test.go         ← write with 0003
```

## pkg/printers/printers_test.go

Until change 0002 lands, `Printer()` writes directly to `os.Stdout`.
Tests capture stdout via `os.Pipe()`.

```go
func captureOutput(f func()) string {
    r, w, _ := os.Pipe()
    old := os.Stdout
    os.Stdout = w
    f()
    w.Close()
    os.Stdout = old
    var buf bytes.Buffer
    io.Copy(&buf, r)
    return buf.String()
}
```

After 0002, `Printer()` accepts an `io.Writer`. Replace `captureOutput`
with direct `bytes.Buffer` passing — no pipe needed.

### Tests

| Test | Purpose |
|------|---------|
| `TestPrintSecrets_NoSensitiveDataInOutput` | Security: known Data values must not appear in output |
| `TestPrinter_EmptyList` | Empty list of any type produces zero output |
| `TestPrinter_PatternFilter` | Matching items shown; non-matching items absent |
| `TestMatchesPattern` | Unit test for the `matchesPattern` helper (add after 0002) |

#### TestPrintSecrets_NoSensitiveDataInOutput (detail)

```go
list := &v1.SecretList{Items: []v1.Secret{{
    ObjectMeta: metav1.ObjectMeta{Name: "my-secret"},
    Type:       v1.SecretTypeOpaque,
    Data: map[string][]byte{
        "password": []byte("super-secret-value"),
        "token":    []byte("my-api-token"),
    },
}}}
out := captureOutput(func() { Printer(list, "") })
// must NOT contain the actual values
assertNotContains(t, out, "super-secret-value")
assertNotContains(t, out, "my-api-token")
// must contain the count
assertContains(t, out, "2")
```

## pkg/util/util_test.go

Requires 0001 because `clientset.CoreV1().Pods(ns).List()` changed
signature from `List(opts)` to `List(ctx, opts)` in client-go v0.21.

Uses `k8s.io/client-go/kubernetes/fake.NewSimpleClientset()`.

### Tests

| Test | Purpose |
|------|---------|
| `TestGetter_CustomKinds` | Only the requested kind is sent over the channel |
| `TestGetter_UnknownKind` | Channel closes cleanly; no deadlock |
| `TestGetter_ChannelAlwaysClosed` | Channel closes even when List calls succeed for all kinds |

## pkg/cache/cache_test.go

Written alongside change 0003. Uses `t.TempDir()` to isolate cache files
per test — no global state. Internal helpers `writeToDir` / `readFromDir`
accept a directory path so tests never touch `~/.kube`.

### Tests

| Test | Purpose |
|------|---------|
| `TestKeyFor_Deterministic` | Same inputs → same SHA-256 key |
| `TestKeyFor_KindsSorted` | `"secret,configmap"` == `"configmap,secret"` key |
| `TestKeyFor_Unique` | Each differing field produces a distinct key |
| `TestReadWrite_RoundTrip` | Written sections survive a Read() call intact |
| `TestRead_Missing` | Returns nil (not error) when file absent |
| `TestRead_Expired` | Returns nil when written_at + ttl < now |
| `TestWrite_Atomic` | Concurrent writes leave a valid JSON file |
| `TestNoSensitiveData` | Raw cache JSON contains no sentinel secret strings |

## Integration tests

Gated by `//go:build integration` tag. Require `kind` binary on PATH.
Test bootstrap creates a temporary kind cluster, runs ksearch against it,
and verifies output + cache behaviour end-to-end.

```go
//go:build integration

func TestIntegration_CacheHit(t *testing.T) {
    // create kind cluster, run ksearch twice, assert second run hits cache
}
```

Run with: `go test -tags integration ./...`
