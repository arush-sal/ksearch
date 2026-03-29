# Spec Delta: testing

Creates: openspec/specs/testing/spec.md (new capability — no prior spec exists)

The full specification is written directly to openspec/specs/testing/spec.md
as part of this change. See that file for the authoritative requirements.

Summary of requirements introduced by this change:
- Unit tests for all packages (pkg/printers, pkg/util, pkg/cache) — no live cluster required
- Secret output safety test — sentinel values must not appear in printer output
- Pattern filter correctness test — matching items shown, non-matching excluded
- Empty list silence test — zero output for empty resource lists
- Getter channel contract test — channel always closes, no deadlock on unknown kind
- Cache security test — raw cache JSON must not contain sentinel secret strings
- Integration tests gated by `integration` build tag, using kind if available
