# Delta: Resource Fetching Spec — Dynamic Discovery

Applies to: `openspec/specs/resource-fetching/spec.md`

## Changes

### Replace: Requirement "Supported resource kinds" (static list)

**Remove:**
```
### Requirement: Supported resource kinds
The system SHALL support listing the following resource kinds by default:
Pods, ConfigMaps, Endpoints, Events, LimitRanges, Namespaces,
PersistentVolumes, PersistentVolumeClaims, PodTemplates, ResourceQuotas,
Secrets, Services, ServiceAccounts, DaemonSets, Deployments, ReplicaSets,
StatefulSets.

#### Scenario: Default resource list is fetched
- GIVEN no --kinds flag is provided
- WHEN ksearch is invoked
- THEN all default resource kinds are queried from the Kubernetes API
```

**Replace with:**
```
### Requirement: Dynamic resource discovery
The system SHALL discover listable resource kinds from the Kubernetes API
server at runtime using the discovery API (ServerGroupsAndResources),
rather than maintaining a static hardcoded list.

Only resource types that advertise the "list" verb SHALL be included.

#### Scenario: Default resource list is discovered
- GIVEN no --kinds flag is provided
- WHEN ksearch is invoked
- THEN all resource kinds that support the "list" verb are queried from
  the Kubernetes discovery API and used as the search set

#### Scenario: Partial discovery failure is tolerated
- GIVEN one API group is unavailable (e.g. a CRD webhook is down)
- WHEN ksearch is invoked
- THEN the error is logged, and all successfully discovered resource kinds
  are still fetched and printed
```

### Add: Requirement "No duplicate resource list"

```
### Requirement: Single source of truth for resource kinds
The list of resource kinds to query SHALL NOT be duplicated across packages.
A single `Discover()` function in `pkg/util` is the authoritative source.
`cmd/ksearch.go` SHALL NOT define its own resource list.
```

### Unchanged requirements

- Custom resource kinds via flag (`--kinds`)
- Namespace scoping (`--namespace`)
- Error resilience (single kind failure continues)
