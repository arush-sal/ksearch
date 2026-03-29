# Contributing to ksearch

Thanks for evaluating or contributing to `ksearch`.

This document is for contributors. If you want to use the CLI, start with [README.md](/mnt/c/Users/micro/go/src/github.com/arush-sal/ksearch/README.md).

## Project Layout

- `main.go` wires the binary entrypoint to Cobra.
- `cmd/ksearch.go` owns CLI flags, cache handling, client initialization, discovery, and output flow.
- `pkg/util/` contains Kubernetes resource discovery and fetch logic.
- `pkg/printers/` contains formatted stdout rendering.
- `pkg/cache/` contains local cache read/write logic.
- `vendor/` stores checked-in dependencies.
- `openspec/` contains specs and completed change records.

## Local Development

Build the CLI:

```bash
make build
```

Run the full test suite:

```bash
make test
```

Run vet:

```bash
make vet
```

Format Go code:

```bash
make fmt
```

Run golangci-lint:

```bash
make lint
```

You can also use standard Go commands directly:

```bash
go test ./...
go vet ./...
gofmt -w .
```

## Development Notes

- Prefer idiomatic Go.
- Keep CLI wiring in `cmd/`.
- Keep fetch and discovery logic in `pkg/util/`.
- Keep output formatting in `pkg/printers/`.
- Keep cache behavior in `pkg/cache/`.
- Use `logrus` for recoverable errors in library packages instead of panics.

## Resource Support

`ksearch` discovers resources dynamically from the cluster and fetches them through typed clients where possible, falling back to unstructured handling where needed.

If you add or change printer behavior, verify the output still makes sense for both typed and unstructured resources.

## Testing

Current test coverage includes:

- cache key and cache persistence behavior
- printer filtering and rendering behavior
- command-level cache flow
- discovery and unstructured resource fetching behavior

Run:

```bash
go test ./...
```

For a fresh run without cache:

```bash
go test -count=1 ./...
```

Integration tests are intended to be run explicitly against a live cluster:

```bash
go test -tags integration ./...
```

## CI and Releases

GitHub Actions currently runs:

- lint
- test with race detection and coverage output
- multi-platform builds

Tagged releases use GoReleaser and publish archives for Linux, macOS, and Windows. Krew metadata is also generated from the release configuration.

## OpenSpec Workflow

- Active implementation guidance lives under `openspec/changes/`.
- Completed changes are moved to `openspec/changes/done/`.
- Specs under `openspec/specs/` describe the intended behavior at a higher level.

If your change maps to an OpenSpec change, keep the related docs aligned.

## Pull Requests

Keep changes scoped and easy to review.

Before opening a PR, verify:

- `make build`
- `make test`
- `make vet`
- any new or changed behavior has test coverage
- CLI behavior changes are reflected in `README.md`

Prefer short imperative commit subjects such as:

- `Add cache TTL override`
- `Fix unstructured resource printing`
- `Improve README installation docs`

## Safety

This tool talks to the current kubeconfig context and may be pointed at a live cluster.

- Test against a non-production namespace when possible.
- Be careful when validating cache behavior because cached output is stored locally.
