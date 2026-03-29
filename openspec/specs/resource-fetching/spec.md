# Resource Fetching Specification

## Purpose
Manages retrieval of Kubernetes resources from the API server for a given
namespace, context, and set of resource kinds.

## Requirements

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

### Requirement: Single source of truth for resource kinds
The list of resource kinds to query SHALL NOT be duplicated across packages.
A single `Discover()` function in `pkg/util` is the authoritative source.
`cmd/ksearch.go` SHALL NOT define its own resource list.

### Requirement: Custom resource kinds via flag
The system SHALL allow users to restrict fetching to a comma-separated
list of resource kinds via the --kinds flag.

#### Scenario: Only specified kinds are fetched
- GIVEN --kinds=configmap,secret is passed
- WHEN ksearch is invoked
- THEN only ConfigMaps and Secrets are queried

### Requirement: Namespace scoping
The system SHALL scope API calls to the namespace provided via --namespace,
or query all namespaces when the flag is omitted.

#### Scenario: Namespace-scoped fetch
- GIVEN --namespace=default is passed
- WHEN ksearch is invoked
- THEN all list calls are scoped to the "default" namespace

### Requirement: Error resilience
The system SHALL log an error and continue when a single resource kind
fails to list, rather than aborting the entire run.

#### Scenario: One kind fails, others succeed
- GIVEN the API server returns an error for Secrets
- WHEN ksearch is invoked
- THEN the error is logged and all other resource kinds are still fetched and printed
