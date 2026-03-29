# Design: Dynamic Resource Discovery

## Data structures

```go
// pkg/util/discover.go

// ResourceMeta holds the metadata for one listable resource type discovered
// from the Kubernetes API server.
type ResourceMeta struct {
    Kind       string // e.g. "Pod"
    Resource   string // plural form, e.g. "pods"
    APIGroup   string // "" for core, "apps" for Deployments, etc.
    APIVersion string // e.g. "v1", "apps/v1"
    Namespaced bool
}
```

## pkg/util/discover.go

```go
package util

import (
    "strings"

    "k8s.io/client-go/discovery"
)

// Discover returns all resource types from the API server that support the
// "list" verb. If kinds is non-empty, results are filtered to the
// comma-separated set of kind names (case-insensitive).
//
// On partial discovery failures (some API groups unavailable), Discover logs
// the errors and returns the successfully discovered resources rather than
// failing completely — consistent with kubectl behaviour.
func Discover(dc discovery.DiscoveryInterface, kinds string) ([]ResourceMeta, error) {
    _, lists, err := dc.ServerGroupsAndResources()
    // ServerGroupsAndResources returns a partial result + error on group failure.
    // We proceed with whatever was discovered.
    if err != nil && lists == nil {
        return nil, err
    }

    filter := parseKindsFilter(kinds)

    var result []ResourceMeta
    for _, list := range lists {
        gv := list.GroupVersion
        for _, r := range list.APIResources {
            if !hasVerb(r.Verbs, "list") {
                continue
            }
            if len(filter) > 0 && !filter[strings.ToLower(r.Kind)] {
                continue
            }
            result = append(result, ResourceMeta{
                Kind:       r.Kind,
                Resource:   r.Name,
                APIGroup:   r.Group,
                APIVersion: gv,
                Namespaced: r.Namespaced,
            })
        }
    }
    return result, nil
}

func parseKindsFilter(kinds string) map[string]bool {
    if kinds == "" {
        return nil
    }
    m := map[string]bool{}
    for _, k := range strings.Split(kinds, ",") {
        m[strings.ToLower(strings.TrimSpace(k))] = true
    }
    return m
}

func hasVerb(verbs []string, target string) bool {
    for _, v := range verbs {
        if v == target {
            return true
        }
    }
    return false
}
```

## pkg/util/util.go — updated Getter signature

Remove `var resources` and `configuredResources()`. Update `Getter`:

```go
// Getter fetches each resource kind in resources from the Kubernetes API and
// sends the result list objects onto c. The channel is always closed on return.
func Getter(namespace string, clientset kubernetes.Interface, resources []ResourceMeta, c chan interface{}) {
    defer close(c)
    ctx := context.Background()

    for _, meta := range resources {
        var list interface{}
        var err error

        switch meta.Kind {
        case "Pod", "Pods":
            list, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
        case "ConfigMap", "ConfigMaps":
            list, err = clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
        // ... all existing cases ...
        default:
            log.Debugf("kind %q is not handled by this printer version, skipping", meta.Kind)
            continue
        }
        handleError(err, meta.Kind)
        if list != nil {
            c <- list
        }
    }
}
```

Note: the `switch` normalises both singular and plural kind names (e.g. `"Pod"`
and `"Pods"`) because the discovery API returns singular forms while the old static
list used plural forms.

## cmd/ksearch.go — updated caller

Remove `var defaultResources` and `effectiveResources()`. Replace with:

```go
cfg := config.GetConfigOrDie()
clientset := kubernetes.NewForConfigOrDie(cfg)

resources, err := util.Discover(clientset.Discovery(), kinds)
if err != nil {
    cmd.PrintErrln(err)
    os.Exit(1)
}

resourceOrder := make([]string, len(resources))
for i, r := range resources {
    resourceOrder[i] = r.Kind
}

getter := make(chan interface{})
go util.Getter(namespace, clientset, resources, getter)

results := make([]cache.SectionEntry, len(resources))
// ... fan-out goroutines as before ...
```

## Cache key compatibility

The cache key already uses the `kinds` flag value (the raw comma-separated string
the user typed). Dynamic discovery does not change this — the key still hashes the
user-supplied `kinds` string, not the expanded resource list. This is correct because
two invocations with the same `kinds` string against the same cluster will discover
the same resources.

## Testing approach

### pkg/util/discover_test.go

Use `k8s.io/client-go/discovery/fake` or construct a `*fakediscovery.FakeDiscovery`
from `fake.NewSimpleClientset().Discovery()` to inject known resource lists.

```go
func TestDiscover_AllWhenEmpty(t *testing.T)      // kinds="" → all listable resources
func TestDiscover_FilterByKinds(t *testing.T)      // kinds="ConfigMaps" → only ConfigMaps
func TestDiscover_SkipsNonListable(t *testing.T)   // resource without "list" verb excluded
func TestDiscover_CaseInsensitive(t *testing.T)    // kinds="configmaps" matches "ConfigMaps"
func TestDiscover_PartialFailureContinues(t *testing.T) // one API group fails → rest returned
```

### pkg/util/util_test.go — updated

Existing `TestGetter_*` tests pass `[]ResourceMeta` instead of a `kinds string`.
Example:

```go
resources := []ResourceMeta{{Kind: "ConfigMaps", Resource: "configmaps", Namespaced: true}}
go Getter("default", fakeClient, resources, ch)
```

## Security impact

None. Discovery only returns kind names and API metadata — no resource data.
The Getter still fetches data through the existing typed client calls.
The Printer still controls what fields are displayed.
