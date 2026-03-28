# Spec Delta: caching

Creates/updates: openspec/specs/caching/spec.md

## New requirement: No sensitive data at rest

```diff
+### Requirement: No sensitive data at rest
+The cache SHALL store only the pre-rendered display output produced by the
+Printer layer. Raw Kubernetes API objects, Secret values, ConfigMap data,
+ServiceAccount tokens, and all other sensitive field values SHALL NOT be
+written to disk.
+
+#### Scenario: Secret data excluded from cache file
+- GIVEN a namespace containing Secrets with non-empty Data fields
+- WHEN ksearch runs and writes a cache entry
+- THEN the cache file contains only the formatted table text (name, type, count)
+- AND the actual secret values are absent from the file
+
+#### Scenario: Cache package has no Kubernetes type imports
+- GIVEN the pkg/cache package source files
+- WHEN grepped for `k8s.io/api` imports or `.Data` field access
+- THEN zero matches are found
```

## Updated requirement: Cache key uniqueness

```diff
-The cache SHALL generate a unique, deterministic key per
-(cluster-context, namespace, kinds) tuple.
+The cache SHALL generate a unique, deterministic key per
+(cluster-context, namespace, kinds, pattern) tuple.
+
+Pattern is included because the cache stores pre-rendered output which is
+already filtered by pattern; a different pattern produces different output
+and must be cached separately.

+#### Scenario: Kind order does not affect key
+- GIVEN kinds="secret,configmap" and kinds="configmap,secret"
+- WHEN KeyFor() is called for each
+- THEN both return the same SHA-256 hex string
```
