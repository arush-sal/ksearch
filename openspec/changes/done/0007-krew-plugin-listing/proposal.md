# Proposal: Krew Plugin Index Listing

## Problem

ksearch is a kubectl plugin, but it cannot be discovered or installed via
`kubectl krew install ksearch`. While `.goreleaser.yml` already contains a
`krews` section, it publishes to a personal `arush-sal/krew-index` repo —
not the official `kubernetes-sigs/krew-index`. Additionally, several gaps
prevent a clean krew experience:

1. **Binary packaging contract**: the revised design keeps the built binary as
   `ksearch`, but parts of the spec still incorrectly require `kubectl-ksearch`.
   The packaging docs and manifests need to be reconciled around krew's symlink
   behavior.
2. **No automated krew-index updates**: each new release requires a manual PR
   to whichever krew-index repo is used.
3. **Help text is not kubectl-aware**: `ksearch --help` shows `ksearch` as the
   command name, not `kubectl ksearch`.
4. **No local validation workflow**: no documented way to test the krew manifest
   + archive locally before publishing.
5. **GoReleaser krew config incomplete**: missing `caveats` field, archive
   `name_template` does not match what krew-release-bot expects.

## Proposal

Make ksearch a first-class krew plugin:

1. Keep the built binary as `ksearch`; rely on the krew manifest `bin` field and
   the symlink krew creates for `kubectl ksearch`.
2. Detect invocation context (`kubectl ksearch` vs direct `ksearch`) and adjust
   help output accordingly.
3. Add `krew-release-bot` GitHub Action to automate krew-index PRs on tag push.
4. Refine `.goreleaser.yml` krew config (caveats, description, target the official index).
5. Add a `.krew.yaml` template for `krew-release-bot`, while keeping GoReleaser
   as the source for local manifest validation.
6. Document local krew validation in the openspec tasks.
7. Submit the initial manifest PR to `kubernetes-sigs/krew-index`.

## Benefits

- Users discover and install via `kubectl krew install ksearch`
- Auto-updated on every tagged release (no manual krew-index PRs)
- Works with krew without changing the standalone `ksearch` binary name
- Broader adoption through the official plugin index

## Risks and mitigations

| Risk | Mitigation |
|------|-----------|
| krew-index PR rejected on naming or quality | Follow naming guide strictly; test locally first |
| Stale specs/docs keep referring to `kubectl-ksearch` | Update proposal, task list, and distribution spec to match the revised design |
| krew-release-bot misconfigured | Test with `--snapshot` before first real tag |
