# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

`ksearch` is a kubectl plugin that lists and searches Kubernetes resources across both `core/v1` and `apps/v1` API groups — including resources that `kubectl get` omits by default (ConfigMaps, Secrets, Endpoints, etc.). It is installed as a binary named `ksearch` on the PATH so that `kubectl ksearch` works via kubectl's plugin discovery.

## Build and run

```bash
# Build the binary
go build -o ksearch .

# Run directly (requires a valid kubeconfig)
./ksearch
./ksearch -n <namespace>
./ksearch -n <namespace> -p <name-pattern>
./ksearch -k configmap,secret -n default

# Run tests
go test ./...

# Run tests for a single package
go test ./pkg/printers/...
go test ./pkg/util/...

# Vet
go vet ./...
```

There is no Makefile. The binary is .gitignored; always rebuild after changes.

## Architecture

The data flow is strictly linear:

```
cmd/ksearch.go  →  pkg/util/util.go  →  pkg/printers/printers.go
   (CLI flags)       (API fetching)          (stdout rendering)
```

**`cmd/ksearch.go`** — Cobra root command. Initialises the Kubernetes client via `sigs.k8s.io/controller-runtime/pkg/client/config.GetConfigOrDie()`, spawns `util.Getter()` in a goroutine, and iterates the results channel calling `printers.Printer()` for each resource.

**`pkg/util/util.go`** — `Getter(namespace, clientset, kinds string, c chan interface{})` iterates a fixed resource list (or a user-supplied comma-separated `kinds` override), calls the appropriate `clientset.CoreV1()` or `clientset.AppsV1()` list method for each kind, and sends the typed list objects (`*v1.PodList`, `*appsv1.DeploymentList`, etc.) over the channel. The channel is closed when all kinds are exhausted.

**`pkg/printers/printers.go`** — `Printer(resource interface{}, resName string)` type-switches on the received object and dispatches to a per-kind printing function. Each function writes an aligned table to stdout using `text/tabwriter` and filters rows by `strings.Contains(name, resName)` when a pattern is set.

## Key design constraints

- The channel between `Getter` and the print loop carries `interface{}` holding concrete typed Kubernetes list structs. The type switch in `Printer()` must handle every type that `Getter()` can send.
- Adding a new resource kind requires changes in **both** `util.go` (add to `resources` slice and `switch`) and `printers.go` (add a print function and a `case` in `Printer()`).
- The default resource list is defined as a package-level `var` in `util.go` and is mutated in-place when `--kinds` is supplied. This means `Getter` is not safe to call concurrently from multiple goroutines.

## Planned changes (openspec/)

Three upgrade tracks are tracked under `openspec/changes/`:

| Change | Summary |
|--------|---------|
| `0001-dependency-upgrade` | Bump Go to 1.22, k8s client-go to v0.32.x; adds `context.Context` to all `.List()` calls; removes deprecated `ComponentStatuses` |
| `0002-concurrent-printing` | Refactor `printers.go` to accept `io.Writer`; extract shared `matchesPattern` helper; fan-out print goroutines with ordered flush |
| `0003-application-caching` | New `pkg/cache` package; SHA-256 keyed disk cache under `~/.kube/ksearch/cache/`; `--cache-ttl` and `--no-cache` flags |

Implement in order: 0001 must land before 0002 or 0003.
