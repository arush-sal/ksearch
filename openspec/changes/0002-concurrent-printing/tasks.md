# Tasks: Concurrent Printing and Printers Refactor

## Phase 1 — Refactor pkg/printers/printers.go
- [ ] Add `matchesPattern(name, pattern string) bool` helper function
- [ ] Update every printer function to accept `(w io.Writer, ..., pattern string)`
      instead of writing to `os.Stdout` directly
- [ ] Replace the duplicated `if resName != "" { if strings.Contains(...) }` guard
      in each function with a single `if !matchesPattern(item.Name, pattern) { continue }`
- [ ] Construct `tabwriter.NewWriter` against the passed-in `w` in each function
- [ ] Update `Printer()` signature to `Printer(w io.Writer, resource interface{}, pattern string)`
- [ ] Add `"io"` import; remove any now-unused `"os"` import

## Phase 2 — Fan-out loop in cmd/ksearch.go
- [ ] Add imports: `"bytes"`, `"sync"`
- [ ] Replace the serial `for resource := range getter { printers.Printer(...) }` loop with:
  - Pre-allocate `results := make([][]byte, len(effectiveResources))`
  - Track index `i` across the range loop
  - Spawn a goroutine per resource that writes `printers.Printer(&buf, res, resName)`
    into a `bytes.Buffer` and stores result in `results[idx]`
  - Call `wg.Wait()` after the range loop
  - Flush `results` slice to `os.Stdout` in order

## Phase 3 — Tests
- [ ] Add `pkg/printers/printers_test.go` with table-driven tests for `matchesPattern`:
  - empty pattern matches everything
  - non-empty pattern matches substring
  - non-matching pattern returns false

## Phase 4 — Verification
- [ ] `go build ./...` — zero errors
- [ ] `go vet ./...` — zero warnings
- [ ] `go test ./pkg/printers/...` — all tests pass
- [ ] Smoke test: `ksearch -n kube-system` output is identical to pre-refactor output
- [ ] Smoke test: `ksearch -n kube-system -p coredns` shows only coredns-matched items
- [ ] Smoke test: `ksearch -n kube-system -k configmap,secret` shows only those two types
