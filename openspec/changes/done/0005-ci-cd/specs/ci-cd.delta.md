# Spec Delta: ci-cd

Creates: openspec/specs/ci-cd/spec.md (new capability — no prior spec exists)

The full specification is written directly to openspec/specs/ci-cd/spec.md
as part of this change. See that file for the authoritative requirements.

Summary of requirements introduced by this change:
- PR validation pipeline: lint (golangci-lint), vet, test, build matrix on every PR
- Cross-platform release binaries for linux/amd64, linux/arm64, darwin/amd64,
  darwin/arm64, windows/amd64 — triggered by semver tags
- Version string embedded via ldflags at build time; reported by --version flag
- Krew plugin manifest (ksearch.yaml) generated and attached to every release
- Go module cache in CI to avoid redundant downloads
- Semver release cadence: patch/minor/major as defined in spec
