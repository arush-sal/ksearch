# Tasks: Dependency Upgrade

## Phase 1 — Module file
- [ ] Update `go.mod`: set `go 1.22`, bump all direct dependency versions to targets in design.md
- [ ] Delete `go.sum` and `vendor/`
- [ ] Run `go mod tidy` to resolve transitive dependencies
- [ ] Run `go mod vendor` to regenerate the vendor tree

## Phase 2 — Source fixes
- [ ] `pkg/util/util.go`: add `"context"` import
- [ ] `pkg/util/util.go`: add `ctx := context.Background()` at top of `Getter()`
- [ ] `pkg/util/util.go`: add `ctx` as first argument to all 18 `.List()` calls
- [ ] `pkg/util/util.go`: remove `"ComponentStatuses"` from the default `resources` slice
- [ ] `pkg/util/util.go`: remove the `case "ComponentStatuses":` block from the switch
- [ ] `pkg/printers/printers.go`: remove `printComponentStatuses()` function
- [ ] `pkg/printers/printers.go`: remove `case resource.(*v1.ComponentStatusList):` from `Printer()`

## Phase 3 — Verification
- [ ] `go build ./...` — must produce zero errors
- [ ] `go vet ./...` — must produce zero warnings
- [ ] Smoke test: `ksearch -n default` against a live cluster — no panics, resources listed
- [ ] `grep -r ComponentStatuses . --include="*.go"` — must return no results
- [ ] `ksearch --help` — --pattern, --namespace, --kinds flags all present
