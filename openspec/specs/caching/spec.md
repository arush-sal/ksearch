# Caching Specification

## Purpose
Manages short-lived disk-based caching of Kubernetes API response output to
reduce redundant network calls on repeated ksearch invocations against the
same cluster context. The cache stores only display-safe rendered text —
never raw Kubernetes API objects or sensitive field values.

## Requirements

### Requirement: No sensitive data at rest
The cache SHALL store only the pre-rendered display output produced by the
Printer layer. Raw Kubernetes API objects, Secret values, ConfigMap data,
ServiceAccount tokens, and all other sensitive field values SHALL NOT be
written to disk.

#### Scenario: Secret data excluded from cache file
- GIVEN a namespace containing Secrets with non-empty Data fields
- WHEN ksearch runs and writes a cache entry
- THEN the cache file contains only the formatted table text (name, type, count)
- AND the actual secret values are absent from the file

#### Scenario: Cache package has no Kubernetes type imports
- GIVEN the pkg/cache package source files
- WHEN grepped for `k8s.io/api` imports or `.Data` field access
- THEN zero matches are found

### Requirement: Cache key uniqueness
The cache SHALL generate a unique, deterministic key per
(cluster-context, namespace, kinds, pattern) tuple.

#### Scenario: Same args produce the same key
- GIVEN the same context, namespace, kinds, and pattern on two invocations
- WHEN KeyFor() is called
- THEN both calls return the identical SHA-256 hex string

#### Scenario: Different args produce different keys
- GIVEN two invocations that differ in any single field (context, namespace, kinds, or pattern)
- WHEN KeyFor() is called for each
- THEN two distinct SHA-256 hex strings are returned

#### Scenario: Kind order does not affect key
- GIVEN kinds="secret,configmap" and kinds="configmap,secret"
- WHEN KeyFor() is called for each
- THEN both return the same SHA-256 hex string

### Requirement: TTL-based invalidation
The cache SHALL invalidate entries older than the configured TTL.

#### Scenario: Entry expired
- GIVEN a cache entry with written_at = now - 90s and TTL = 60s
- WHEN Read() is called
- THEN nil is returned (cache miss)

#### Scenario: Entry valid
- GIVEN a cache entry with written_at = now - 30s and TTL = 60s
- WHEN Read() is called
- THEN the cached entry is returned (cache hit)

### Requirement: Configurable TTL
The cache SHALL allow the TTL to be overridden via --cache-ttl flag
or KSEARCH_CACHE_TTL environment variable. The default TTL SHALL be 60s.

#### Scenario: Custom TTL via flag
- GIVEN --cache-ttl=10s is passed
- WHEN a cache entry is written and 15 seconds pass
- THEN the next invocation treats the entry as expired

### Requirement: Cache bypass
The cache SHALL skip both reading and writing when --no-cache is passed,
and SHALL overwrite any existing cache entry after a fresh fetch.

#### Scenario: --no-cache bypasses and refreshes
- GIVEN a valid cache entry exists
- WHEN ksearch is invoked with --no-cache
- THEN a fresh API call is made and the cache entry is overwritten

### Requirement: Persistent disk storage
The cache SHALL be stored as JSON files under ~/.kube/ksearch/cache/
at mode 0600, with the directory at mode 0700.

#### Scenario: Cache survives process restart
- GIVEN ksearch was run and produced a cache file
- WHEN the process exits and ksearch is run again within TTL
- THEN the cached result is served from disk without any API call

### Requirement: Safe concurrent access
The cache SHALL not corrupt cache files when two ksearch processes run
simultaneously against the same cache key.

#### Scenario: Atomic write
- GIVEN two processes writing the same cache key simultaneously
- WHEN both Write() calls complete
- THEN the resulting file is valid JSON (no partial writes)
