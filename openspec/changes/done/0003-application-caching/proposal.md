# Proposal: Application-Level Caching

## Purpose
Every ksearch invocation issues N Kubernetes API calls (one per resource type).
On large clusters or slow network paths, repeated queries for the same
namespace and kinds against the same cluster are expensive. A file-based
cache keyed on cluster context + namespace + kinds with a configurable TTL
eliminates redundant API traffic for short-lived re-queries.

## Requirements

### Requirement: Cache keyed on cluster context
The cache SHALL be scoped per kubeconfig context (cluster + user).

#### Scenario: Different contexts produce different cache entries
- GIVEN kubeconfig has contexts "prod" and "staging"
- WHEN ksearch is run against "prod" then "staging"
- THEN the cache stores two independent entries and they do not interfere

### Requirement: Configurable TTL with a sane default
The cache SHALL expire entries after a configurable duration (default: 60s).
Users SHALL override the TTL via --cache-ttl flag or KSEARCH_CACHE_TTL env var.

#### Scenario: Cache miss after TTL expiry
- GIVEN a cache entry written 90 seconds ago with TTL=60s
- WHEN ksearch is invoked with the same args
- THEN the cache entry is ignored and a fresh API call is made

#### Scenario: Cache hit within TTL
- GIVEN a cache entry written 30 seconds ago with TTL=60s
- WHEN ksearch is invoked with the same args
- THEN the cached result is returned without any Kubernetes API call

### Requirement: Cache bypass flag
Users SHALL be able to disable the cache for a single invocation via --no-cache.

#### Scenario: --no-cache bypasses and refreshes cache
- GIVEN a valid cache entry exists
- WHEN ksearch is invoked with --no-cache
- THEN a fresh API call is made and the cache entry is overwritten

### Requirement: Cache stored on disk per user
The cache SHALL be stored under ~/.kube/ksearch/cache/ as JSON files
at mode 0600 with the directory at mode 0700.

#### Scenario: Cache survives process restart
- GIVEN ksearch was run and produced a cache file
- WHEN the process exits and ksearch is run again within TTL
- THEN the cached result is served from disk without any API call
