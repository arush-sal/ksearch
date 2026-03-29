# Design: Concurrent Printing and Printers Refactor

## Current problem

`pkg/printers/printers.go` has two issues:
1. 17 printer functions each repeat the same guard:
   ```go
   if resName != "" {
       if strings.Contains(item.Name, resName) { ... }
   } else { ... }
   ```
2. `Printer()` is called serially in `cmd/ksearch.go`:
   ```go
   for resource := range getter {
       printers.Printer(resource, resName)  // blocks until printed
   }
   ```

## Approach: buffer-per-resource + fan-out goroutines

Each resource received from the `getter` channel is dispatched to a
goroutine that renders it into a `bytes.Buffer`. After all goroutines
complete (via `sync.WaitGroup`), the main thread flushes the buffers
to stdout in deterministic order.

```
getter channel
      │
      ▼
 fan-out loop ──► goroutine 1 → Printer(resource, w) → results[i]
                 goroutine 2 → Printer(resource, w) → results[j]
                 goroutine N → Printer(resource, w) → results[k]
                      │
                 wg.Wait()
                      │
                 flush results[0..N] to os.Stdout in order
```

The `results` slice is pre-allocated with one slot per position in the
`resources` list. Each goroutine writes to its own slot — no mutex needed.

To map a received resource to its slot index, the goroutine receives the
index from the fan-out loop (the index at which the resource was sent by
`Getter()`). Since `Getter()` sends in `resources` slice order, the index
can be tracked with a counter in the fan-out loop.

## Printers refactor

### 1. Shared filter helper

```go
// matchesPattern returns true if pattern is empty or name contains pattern.
func matchesPattern(name, pattern string) bool {
    return pattern == "" || strings.Contains(name, pattern)
}
```

All 17 printer functions replace their duplicated guard with a single call:
```go
if !matchesPattern(item.Name, pattern) {
    continue
}
```

### 2. io.Writer parameter

Every printer function signature changes from:
```go
func printConfigMaps(cms *v1.ConfigMapList, resName string)
```
to:
```go
func printConfigMaps(w io.Writer, cms *v1.ConfigMapList, pattern string)
```

`tabwriter.NewWriter` is constructed against `w` instead of `os.Stdout`.

### 3. Printer() signature change

```go
// Before
func Printer(resource interface{}, resName string)

// After
func Printer(w io.Writer, resource interface{}, pattern string)
```

### 4. cmd/ksearch.go fan-out loop

```go
results := make([][]byte, len(resources))
var wg sync.WaitGroup
i := 0
for resource := range getter {
    idx := i
    i++
    wg.Add(1)
    go func(idx int, res interface{}) {
        defer wg.Done()
        var buf bytes.Buffer
        printers.Printer(&buf, res, resName)
        results[idx] = buf.Bytes()
    }(idx, resource)
}
wg.Wait()
for _, b := range results {
    os.Stdout.Write(b)
}
```

## Files changed

| File | Change |
|------|--------|
| `pkg/printers/printers.go` | Add `matchesPattern`; change all printer functions to accept `io.Writer`; update `Printer()` signature |
| `cmd/ksearch.go` | Replace serial loop with fan-out + WaitGroup + ordered flush |

## No-change items
- `pkg/util/util.go` — `Getter()` and channel protocol unchanged
- CLI flags
- Column headers and data fields per resource type
