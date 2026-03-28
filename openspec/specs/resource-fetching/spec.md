# Resource Fetching Specification

## Purpose
Manages retrieval of Kubernetes resources from the API server for a given
namespace, context, and set of resource kinds.

## Requirements

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
