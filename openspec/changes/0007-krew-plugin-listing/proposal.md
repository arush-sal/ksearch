# Proposal: Krew Plugin Index Listing

## Problem

ksearch is a kubectl plugin, but it cannot be discovered or installed via
`kubectl krew install ksearch`. While `.goreleaser.yml` already contains a
`krews` section, it publishes to a personal `arush-sal/krew-index` repo —
not the official `kubernetes-sigs/krew-index`. Additionally, several gaps
prevent a clean krew experience:

1. **Binary name**: built as `ksearch`, not `kubectl-ksearch`. Without krew,
   placing the binary on PATH does not register it as a kubectl plugin.
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

1. Rename the built binary to `kubectl-ksearch` (works both with and without krew).
2. Detect invocation context (`kubectl ksearch` vs `kubectl-ksearch`) and adjust
   help output accordingly.
3. Add `krew-release-bot` GitHub Action to automate krew-index PRs on tag push.
4. Refine `.goreleaser.yml` krew config (caveats, description, target the official index).
5. Document local krew validation in the openspec tasks.
6. Submit the initial manifest PR to `kubernetes-sigs/krew-index`.

## Benefits

- Users discover and install via `kubectl krew install ksearch`
- Auto-updated on every tagged release (no manual krew-index PRs)
- Works without krew too: `kubectl-ksearch` on PATH is auto-discovered by kubectl
- Broader adoption through the official plugin index

## Risks and mitigations

| Risk | Mitigation |
|------|-----------|
| krew-index PR rejected on naming or quality | Follow naming guide strictly; test locally first |
| Binary rename breaks existing users | `ksearch` binary name was never published via krew; no backward compat concern |
| krew-release-bot misconfigured | Test with `--snapshot` before first real tag |
