# Spec Delta: resource-fetching

Modifies: openspec/specs/resource-fetching/spec.md

## Requirement: Supported resource kinds

```diff
 The system SHALL support listing the following resource kinds by default:
-Pods, ComponentStatuses, ConfigMaps, Endpoints, Events, LimitRanges,
-Namespaces, PersistentVolumes, PersistentVolumeClaims, PodTemplates,
-ResourceQuotas, Secrets, Services, ServiceAccounts,
-DaemonSets, Deployments, ReplicaSets, StatefulSets.
+Pods, ConfigMaps, Endpoints, Events, LimitRanges, Namespaces,
+PersistentVolumes, PersistentVolumeClaims, PodTemplates, ResourceQuotas,
+Secrets, Services, ServiceAccounts, DaemonSets, Deployments, ReplicaSets,
+StatefulSets.
+
+Note: ComponentStatuses was removed — the API was deprecated in Kubernetes
+1.19 and removed in 1.22+. Querying it on modern clusters returns a gone error.
```
