# Spec Delta: output-formatting

Modifies: openspec/specs/output-formatting/spec.md

## Requirement: Printer output target

```diff
-Each printer function SHALL write directly to os.Stdout.
+Each printer function SHALL accept an io.Writer and write to it,
+allowing callers to capture output into a buffer before flushing.
```

## Requirement: Pattern filter ownership

```diff
-Each printer function SHALL implement its own name-pattern filter check
+A single matchesPattern(name, pattern string) bool helper SHALL own
+all name-pattern filtering; printer functions SHALL delegate to it.
```

## Requirement: Deterministic output order (NEW)

```diff
+The system SHALL output resource sections in the same order as the
+resource kinds list, regardless of the order in which goroutines complete.
+
+#### Scenario: Ordered output under concurrency
+- GIVEN Pods and Deployments are rendered concurrently
+- WHEN both goroutines finish
+- THEN Pods section always appears before Deployments section in stdout
```
