# Spec Delta: resource-fetching

Modifies: openspec/specs/resource-fetching/spec.md

## Requirement: Cache-aware fetching (NEW)

```diff
+### Requirement: Cache-aware fetching
+The system SHALL check the application cache before issuing Kubernetes
+API calls and SHALL skip API calls entirely on a cache hit.
+
+#### Scenario: Cache hit skips API
+- GIVEN a valid cache entry exists for the current (context, namespace, kinds)
+- WHEN ksearch is invoked without --no-cache
+- THEN no Kubernetes API calls are made
+- AND the cached resource lists are passed to the Printer pipeline
+
+#### Scenario: Cache miss triggers normal fetch
+- GIVEN no cache entry exists (or the entry is expired)
+- WHEN ksearch is invoked
+- THEN Getter() is called as normal
+- AND the results are written to the cache before printing
```
