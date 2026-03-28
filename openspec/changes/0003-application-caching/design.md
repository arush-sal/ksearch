# Design: Application-Level Caching

## Cache key

The cache key is a lowercase hex SHA-256 hash of the concatenation:
  `<context>\x00<namespace>\x00<sorted-lowercased-kinds>\x00<pattern>`

Using the null byte as a separator prevents collisions between values.
Kinds are sorted before hashing so that `-k secret,configmap` and
`-k configmap,secret` produce the same key. Pattern is included because
the cache stores pre-rendered (already-filtered) output; a different
pattern produces different visible output and therefore needs a separate
cache entry.

Example:
```
context   = "my-cluster"
namespace = "default"
kinds     = "configmap,secret"  →  sorted: "configmap,secret"
pattern   = "nginx"

key = sha256("my-cluster\x00default\x00configmap,secret\x00nginx")
    = "3f2e1a..."
```

## Cache file format

Location: `~/.kube/ksearch/cache/<sha256-key>.json`
Mode: 0600

```json
{
  "written_at":   "2026-03-28T12:00:00Z",
  "ttl_seconds":  60,
  "context":      "my-cluster",
  "namespace":    "default",
  "kinds":        "configmap,secret",
  "pattern":      "nginx",
  "sections": [
    { "kind": "ConfigMaps", "output": "<base64(rendered tabwriter bytes)>" },
    { "kind": "Secrets",    "output": "<base64(rendered tabwriter bytes)>" }
  ]
}
```

The `output` field contains only the display-safe text that would have been
written to stdout — **never raw Kubernetes API objects or field values**.
This is the key security property: what is safe to show on a terminal is
safe to store on disk. No Kubernetes list objects are ever serialised.

## Security guarantee

The cache package MUST NOT import or reference any `k8s.io/api` types.
The only data that flows into `Write()` is `[]SectionEntry` — pre-rendered
text from the Printer layer. Secret values, ConfigMap data, service account
tokens, and all other sensitive field values are excluded by construction
because the printer never emits them.

Verification: `grep -r "v1\." pkg/cache/` and `grep -r "\.Data" pkg/cache/`
must both return zero results.

## New package: pkg/cache

```
pkg/cache/
  cache.go       — public API
  cache_test.go  — unit tests
```

### Public API

```go
package cache

type SectionEntry struct {
    Kind   string `json:"kind"`
    Output string `json:"output"` // base64(rendered []byte)
}

type CacheMeta struct {
    Context    string
    Namespace  string
    Kinds      string
    Pattern    string
    TTLSeconds int
}

type CacheEntry struct {
    WrittenAt  time.Time      `json:"written_at"`
    TTLSeconds int            `json:"ttl_seconds"`
    Context    string         `json:"context"`
    Namespace  string         `json:"namespace"`
    Kinds      string         `json:"kinds"`
    Pattern    string         `json:"pattern"`
    Sections   []SectionEntry `json:"sections"`
}

// KeyFor returns the SHA-256 cache key for the given invocation tuple.
// Kinds are sorted internally before hashing.
func KeyFor(context, namespace, kinds, pattern string) string

// Read loads a cache entry from disk; returns nil if missing or expired.
func Read(key string, ttl time.Duration) (*CacheEntry, error)

// Write saves pre-rendered section output to disk atomically.
// sections contains only display-safe text — no raw Kubernetes objects.
func Write(key string, meta CacheMeta, sections []SectionEntry) error
```

### Atomic write

To avoid partial writes under concurrent access, `Write()` writes to a
temp file in the same directory then calls `os.Rename()`. On POSIX
systems this is atomic.

```go
tmp, _ := os.CreateTemp(cacheDir, "*.tmp")
json.NewEncoder(tmp).Encode(entry)
tmp.Close()
os.Rename(tmp.Name(), filepath.Join(cacheDir, key+".json"))
```

## Integration in cmd/ksearch.go

```
startup
  │
  ├─ derive cache key (context + namespace + kinds + pattern)
  │
  ├─ if !noCache:
  │     entry = cache.Read(key, ttl)
  │     if entry != nil:
  │         range entry.Sections → write each decoded Output to os.Stdout
  │         return
  │
  └─ (cache miss or --no-cache)
       run util.Getter() as today
       fan-out Printer() per kind (change 0002) → collect []SectionEntry
       cache.Write(key, meta, sections)
       flush sections to os.Stdout in order
```

No Kubernetes objects ever leave the Printer pipeline into the cache layer.
The cache layer is intentionally kept ignorant of all Kubernetes types.

### Kubeconfig context detection

```go
import "k8s.io/client-go/tools/clientcmd"

rules := clientcmd.NewDefaultClientConfigLoadingRules()
cfg, _ := rules.Load()
currentContext := cfg.CurrentContext
```

## New flags

| Flag          | Default | Env var           | Description                   |
|---------------|---------|-------------------|-------------------------------|
| --cache-ttl   | 60s     | KSEARCH_CACHE_TTL | Duration before cache expires |
| --no-cache    | false   | —                 | Skip cache for this invocation|

## Cache directory bootstrap

```go
os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".kube", "ksearch", "cache"), 0700)
```

Called once at the top of `cache.Write()`. No separate init step.

## Eviction

Entries expire passively: `Read()` checks `written_at + ttl <= now` and
returns nil on expiry. Stale files are overwritten on the next cache miss.
No background daemon or cron job is required. A `ksearch cache clean`
subcommand for bulk eviction is deferred to a future change.
