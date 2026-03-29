# Design: Dependency Upgrade

## Target versions

| Dependency                     | Current  | Target  |
|--------------------------------|----------|---------|
| Go                             | 1.13     | 1.22    |
| k8s.io/api                     | v0.17.0  | v0.32.x |
| k8s.io/apimachinery            | v0.17.0  | v0.32.x |
| k8s.io/client-go               | v0.17.0  | v0.32.x |
| sigs.k8s.io/controller-runtime | v0.4.0   | v0.20.x |
| github.com/spf13/cobra         | v0.0.5   | v1.8.x  |
| github.com/sirupsen/logrus     | v1.4.2   | v1.9.x  |

## Breaking changes requiring source edits

### 1. client-go List calls require context (since v0.21)

All 18 `.List(metav1.ListOptions{})` calls in `pkg/util/util.go` must become
`.List(ctx, metav1.ListOptions{})`. A single `ctx := context.Background()` is
created once at the top of `Getter()` and reused for all calls.

**Before:**
```go
list, err = clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
```
**After:**
```go
ctx := context.Background()
// ...
list, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
```

### 2. ComponentStatuses removed in Kubernetes 1.22+

`ComponentStatuses` was deprecated in 1.19 and removed in 1.22. Calls to it
return a 404/gone error on modern clusters. Remove it from:
- The default `resources` slice in `pkg/util/util.go`
- Its `switch` case in `pkg/util/util.go`
- The `printComponentStatuses()` function in `pkg/printers/printers.go`
- Its `case` in `Printer()` in `pkg/printers/printers.go`

### 3. Vendor directory

The current `vendor/` directory must be deleted and regenerated:
```bash
rm -rf vendor/ go.sum
go mod tidy
go mod vendor
```

## No-change items
- Public signatures of `Getter()` and `Printer()`
- CLI flags
- Output format
- Module path (`github.com/infracloudio/ksearch`)
