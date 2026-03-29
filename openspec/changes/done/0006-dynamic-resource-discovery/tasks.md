# Tasks: Dynamic Resource Discovery

## Phase 1 — pkg/util/discover.go (new file)
- [ ] Define `ResourceMeta` struct with fields: `Kind`, `Resource`, `APIGroup`, `APIVersion string`, `Namespaced bool`
- [ ] Add `Discover(dc discovery.DiscoveryInterface, kinds string) ([]ResourceMeta, error)`
  - Call `dc.ServerGroupsAndResources()`
  - On partial failure (err != nil but lists != nil): log warning, continue with partial results
  - Filter to resources whose `Verbs` includes `"list"`
  - If `kinds` non-empty: filter to comma-separated kind names (case-insensitive)
- [ ] Add `parseKindsFilter(kinds string) map[string]bool` helper
- [ ] Add `hasVerb(verbs []string, target string) bool` helper
- [ ] Import `"k8s.io/client-go/discovery"` and `"strings"`

## Phase 2 — pkg/util/util.go
- [ ] Remove `var resources = []string{...}` (lines 12–30)
- [ ] Remove `configuredResources(kinds string) []string` helper
- [ ] Update `Getter()` signature to `Getter(namespace string, clientset kubernetes.Interface, resources []ResourceMeta, c chan interface{})`
- [ ] Replace `for _, resource := range configuredResources(kinds)` with `for _, meta := range resources`
- [ ] Update `switch resource` to `switch meta.Kind`
- [ ] Normalise both singular and plural kind names in switch cases (e.g. `case "Pod", "Pods":`)
- [ ] Replace `log.Error("Given kind is not supported"); return` with `log.Debugf("kind %q not handled, skipping", meta.Kind); continue`
- [ ] Ensure `defer close(c)` is still the first line of `Getter()`

## Phase 3 — cmd/ksearch.go
- [ ] Remove `var defaultResources = []string{...}` (lines 33–51)
- [ ] Remove `effectiveResources(kinds string) []string` helper (lines 53–59)
- [ ] After `clientset := kubernetes.NewForConfigOrDie(cfg)`, call:
  ```go
  resources, err := util.Discover(clientset.Discovery(), kinds)
  if err != nil { cmd.PrintErrln(err); os.Exit(1) }
  ```
- [ ] Build `resourceOrder` from `resources`:
  ```go
  resourceOrder := make([]string, len(resources))
  for i, r := range resources { resourceOrder[i] = r.Kind }
  ```
- [ ] Pass `resources` to `util.Getter()`: `go util.Getter(namespace, clientset, resources, getter)`
- [ ] Remove `"strings"` import if no longer used (was used by `effectiveResources`)
- [ ] Ensure `results := make([]cache.SectionEntry, len(resources))` (not `len(resourceOrder)`)

## Phase 4 — pkg/util/discover_test.go (new file)
- [ ] `TestDiscover_AllWhenEmpty` — `kinds=""` → all listable resources returned
- [ ] `TestDiscover_FilterByKinds` — `kinds="ConfigMaps"` → only ConfigMaps in result
- [ ] `TestDiscover_SkipsNonListable` — resource without `"list"` in verbs excluded
- [ ] `TestDiscover_CaseInsensitive` — `kinds="configmaps"` matches `Kind="ConfigMaps"`
- [ ] `TestDiscover_MultipleKinds` — `kinds="ConfigMaps,Secrets"` → exactly those two kinds
- [ ] `TestDiscover_PartialFailureContinues` — injected partial failure still returns discovered resources

## Phase 5 — pkg/util/util_test.go (update)
- [ ] Update `TestGetter_CustomKinds` to pass `[]ResourceMeta{{Kind: "ConfigMaps", Resource: "configmaps", Namespaced: true}}`
- [ ] Update `TestGetter_UnknownKind` to pass `[]ResourceMeta{{Kind: "NonExistentKind", Resource: "nonexistentkinds", Namespaced: true}}`
- [ ] Update `TestGetter_ChannelAlwaysClosed` to pass `[]ResourceMeta{{Kind: "Pods", Resource: "pods", Namespaced: true}}`

## Phase 6 — Verification
- [ ] `go build ./...` — zero errors
- [ ] `go vet ./...` — zero warnings
- [ ] `go test -race ./pkg/util/...` — all tests pass
- [ ] `grep -rn "defaultResources" . --include="*.go"` — zero results
- [ ] `grep -rn "effectiveResources" . --include="*.go"` — zero results
- [ ] `grep -rn "var resources" pkg/util/ --include="*.go"` — zero results
- [ ] `grep -rn "configuredResources" . --include="*.go"` — zero results
- [ ] Smoke test: `ksearch -n kube-system` produces output equivalent to pre-refactor
- [ ] Smoke test: `ksearch -n kube-system -k configmap,secret` shows only those two types
