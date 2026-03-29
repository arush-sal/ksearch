# Proposal: Dynamic Resource Discovery

## Problem

ksearch hard-codes a static list of Kubernetes resource kinds in **two places**:

- `pkg/util/util.go:12–30` — `var resources = []string{...}` + `configuredResources()`
- `cmd/ksearch.go:33–51` — `var defaultResources = []string{...}` + `effectiveResources()`

The two lists are identical. Both require manual maintenance whenever a resource type
is added to or removed from Kubernetes. New CRDs, alpha resources, and types removed
in newer API versions are never reflected without a code change and release.

## Proposal

Replace both static lists with a live call to the Kubernetes discovery API
(`ServerPreferredResources`), filtered to resource types that support the `list` verb.

Expose a single `Discover(dc discovery.DiscoveryInterface, kinds string) ([]ResourceMeta, error)`
function from `pkg/util`. Remove both static lists and both helper functions entirely.

`cmd/ksearch.go` calls `Discover()` first to build the resource set, then passes
`[]ResourceMeta` to `Getter()`. No resource-kind logic lives in `cmd/` at all.

## Benefits

- Automatically correct across all Kubernetes versions without code changes
- Eliminates the duplicated static list
- Works with CRDs and operator-installed resource types automatically
- Cluster-specific: `ksearch` only queries resources that actually exist in the
  target cluster

## Risks and mitigations

| Risk | Mitigation |
|------|-----------|
| Discovery round-trip adds startup latency | Discovery is cheap (single HTTP call, cached by client-go) |
| Printer does not handle every discovered kind | Unrecognised kinds are logged at debug level and skipped; not an error |
| Discovery client unavailable | `Discover()` returns an error; `cmd/` exits with a clear message |

## Scope

- `pkg/util/discover.go` — new file, `ResourceMeta` struct, `Discover()` function
- `pkg/util/util.go` — remove static list, update `Getter()` signature
- `cmd/ksearch.go` — remove duplicate list, call `Discover()` before `Getter()`
- `pkg/util/discover_test.go` — unit tests using fake discovery client
- `openspec/specs/resource-fetching/spec.md` — update dynamic-discovery requirement
