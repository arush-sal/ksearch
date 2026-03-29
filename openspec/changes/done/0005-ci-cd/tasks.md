# Tasks: CI/CD Pipeline and Release Cycle

## Phase 1 — Lint configuration
- [ ] Create `.golangci.yml` with linters: staticcheck, errcheck, govet, gofmt,
      gosimple, ineffassign, unused
- [ ] Run `golangci-lint run ./...` locally and fix any existing violations

## Phase 2 — Version injection
- [ ] Add `var version = "dev"` to `cmd/ksearch.go`
- [ ] Set `Version: version` on the Cobra root command (replaces hardcoded `"v0.0.1"`)
- [ ] Verify: `go build -ldflags "-X main.version=v0.0.2" -o ksearch . && ./ksearch --version`
      outputs `v0.0.2`

## Phase 3 — CI workflow
- [ ] Create `.github/workflows/ci.yml` with jobs: lint, test, build
  - `lint`: golangci/golangci-lint-action@v6
  - `test`: `go test -race -coverprofile=coverage.out ./...` + `go vet ./...`
  - `build`: matrix over linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
- [ ] Pin Go version via `go-version-file: go.mod` (uses the version declared in go.mod)
- [ ] Enable `cache: true` on `actions/setup-go` for module caching

## Phase 4 — Release workflow and GoReleaser
- [ ] Create `.goreleaser.yml` with:
  - `builds`: ldflags injecting version, matrix of goos/goarch
  - `archives`: tar.gz for linux/darwin, zip for windows; include LICENSE
  - `checksum`: sha256, name `checksums.txt`
  - `krew`: generate manifest for all five platforms
- [ ] Create `.github/workflows/release.yml` using `goreleaser/goreleaser-action@v6`
      triggered on `tags: ['v*.*.*']`
- [ ] Test locally: `goreleaser release --snapshot --clean`

## Phase 5 — PR template
- [ ] Create `.github/pull_request_template.md` with checklist from design doc

## Phase 6 — Verification
- [ ] Push a branch, open a PR → CI workflow triggers and passes
- [ ] `go build ./...` + `go vet ./...` + `go test ./...` all pass in CI
- [ ] Create a test tag `v0.0.1-test`, push it → release workflow triggers
- [ ] Verify GitHub Release contains all five platform archives + `checksums.txt`
      + `ksearch.yaml` krew manifest
- [ ] Download a release binary, run `ksearch --version` → outputs the tag version
- [ ] Delete the test tag after verification
