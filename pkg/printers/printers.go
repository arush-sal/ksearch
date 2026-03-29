package printers

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func matchesPattern(name, pattern string) bool {
	return pattern == "" || strings.Contains(name, pattern)
}

func writef(w io.Writer, format string, values ...interface{}) {
	_, _ = fmt.Fprintf(w, format, values...)
}

func flushTabWriter(tabWriter *tabwriter.Writer) {
	_ = tabWriter.Flush()
}

func printPodDetails(w io.Writer, pods *v1.PodList, pattern string) {
	if len(pods.Items) > 0 {
		writef(w, "\nPods\n----\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "READY", "STATUS", "RESTARTS")

		for _, pod := range pods.Items {
			if !matchesPattern(pod.Name, pattern) {
				continue
			}

			writef(tabWriter, "%v\t%v\t%v\t%v\n", pod.Name, "", pod.Status.Phase, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printPodTemplates(w io.Writer, podTemplates *v1.PodTemplateList, pattern string) {
	if len(podTemplates.Items) > 0 {
		writef(w, "\nPodTemplates\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\n", "NAME")
		for _, podTemplate := range podTemplates.Items {
			if !matchesPattern(podTemplate.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\n", podTemplate.Name)
		}
		flushTabWriter(tabWriter)
	}
}
func printConfigMaps(w io.Writer, cms *v1.ConfigMapList, pattern string) {
	if len(cms.Items) > 0 {
		writef(w, "\nConfigMaps\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\n", "NAME", "DATA", "AGE")
		for _, configMap := range cms.Items {
			if !matchesPattern(configMap.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\n", configMap.Name, len(configMap.Data), "")
		}
		flushTabWriter(tabWriter)
	}
}
func printEndpoints(w io.Writer, endPoints *v1.EndpointsList, pattern string) {
	if len(endPoints.Items) > 0 {
		writef(w, "\nEndpoints\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\n", "NAME", "ENDPOINTS", "AGE")
		for _, endpoint := range endPoints.Items {
			if !matchesPattern(endpoint.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\n", endpoint.Name, "", "")
		}
		flushTabWriter(tabWriter)
	}
}
func printEvents(w io.Writer, events *v1.EventList, pattern string) {
	if len(events.Items) > 0 {
		writef(w, "\nEvents\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE")
		for _, event := range events.Items {
			if !matchesPattern(event.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", event.Namespace, "", event.Type, "", event.InvolvedObject.Kind+"/"+event.InvolvedObject.Name, event.Message)
		}
		flushTabWriter(tabWriter)
	}
}
func printLimitRanges(w io.Writer, limitRanges *v1.LimitRangeList, pattern string) {
	if len(limitRanges.Items) > 0 {
		writef(w, "\nLimitRanges\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
		for _, limitRange := range limitRanges.Items {
			if !matchesPattern(limitRange.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\n", limitRange.Name, limitRange.CreationTimestamp)
		}
		flushTabWriter(tabWriter)
	}
}
func printNamespaces(w io.Writer, namespaces *v1.NamespaceList, pattern string) {
	if len(namespaces.Items) > 0 {
		writef(w, "\nNamespaces\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\n", "NAME", "STATUS", "AGE")
		for _, namespace := range namespaces.Items {
			if !matchesPattern(namespace.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\n", namespace.Name, namespace.Status, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printPVs(w io.Writer, pvs *v1.PersistentVolumeList, pattern string) {
	if len(pvs.Items) > 0 {
		writef(w, "\nPersistentVolumes\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGECLASS", "REASON", "AGE")

		for _, pv := range pvs.Items {
			if !matchesPattern(pv.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pv.Name, "", pv.Spec.AccessModes, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.Spec.ClaimRef.Namespace+"/"+pv.Spec.ClaimRef.Name, pv.Spec.StorageClassName, pv.Status.Reason, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printPVCs(w io.Writer, pvcs *v1.PersistentVolumeClaimList, pattern string) {
	if len(pvcs.Items) > 0 {
		writef(w, "\nPersistentVolumeClaims\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS MODES", "STORAGECLASS", "AGE")
		for _, pvc := range pvcs.Items {
			if !matchesPattern(pvc.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pvc.Name, pvc.Status, "", pvc.Status.Capacity.Cpu(), pvc.Spec.AccessModes, pvc.Spec.StorageClassName, "")
		}
		flushTabWriter(tabWriter)
	}
}

func printResourceQuotas(w io.Writer, resQuotas *v1.ResourceQuotaList, pattern string) {
	if len(resQuotas.Items) > 0 {
		writef(w, "\nResourceQuotas\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
		for _, resQ := range resQuotas.Items {
			if !matchesPattern(resQ.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\n", resQ.Name, resQ.CreationTimestamp)
		}
		flushTabWriter(tabWriter)
	}
}
func printSecrets(w io.Writer, secrets *v1.SecretList, pattern string) {
	if len(secrets.Items) > 0 {
		writef(w, "\nSecrets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "TYPE", "DATA", "AGE")
		for _, secret := range secrets.Items {
			if !matchesPattern(secret.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\n", secret.Name, secret.Type, len(secret.Data), "")
		}
		flushTabWriter(tabWriter)
	}
}
func printServices(w io.Writer, services *v1.ServiceList, pattern string) {
	if len(services.Items) > 0 {
		writef(w, "\nServices\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE")

		for _, service := range services.Items {
			if !matchesPattern(service.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", service.Name, service.Spec.Type, service.Spec.ClusterIP, service.Spec.ExternalIPs, service.Spec.Ports, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printServiceAccounts(w io.Writer, serviceAccs *v1.ServiceAccountList, pattern string) {
	if len(serviceAccs.Items) > 0 {
		writef(w, "\nServiceAccounts\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\n", "NAME", "SECRETS", "AGE")
		for _, serviceAcc := range serviceAccs.Items {
			if !matchesPattern(serviceAcc.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\n", serviceAcc.Name, len(serviceAcc.Secrets), "")
		}
		flushTabWriter(tabWriter)
	}
}
func printDaemonSets(w io.Writer, daemonsets *appsv1.DaemonSetList, pattern string) {
	if len(daemonsets.Items) > 0 {
		writef(w, "\nDaemonSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE")
		for _, ds := range daemonsets.Items {
			if !matchesPattern(ds.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printDeployments(w io.Writer, deployments *appsv1.DeploymentList, pattern string) {
	if len(deployments.Items) > 0 {
		writef(w, "\nDeployments\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
		for _, deployment := range deployments.Items {
			if !matchesPattern(deployment.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
		}
		flushTabWriter(tabWriter)
	}
}
func printReplicaSets(w io.Writer, rsets *appsv1.ReplicaSetList, pattern string) {
	if len(rsets.Items) > 0 {
		writef(w, "\nReplicaSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "DESIRED", "CURRENT", "READY", "AGE")
		for _, rs := range rsets.Items {
			if !matchesPattern(rs.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", rs.Name, "", "", "", "")
		}
		flushTabWriter(tabWriter)
	}
}
func printStateFulSets(w io.Writer, ssets *appsv1.StatefulSetList, pattern string) {
	if len(ssets.Items) > 0 {
		writef(w, "\nStatefulSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		writef(tabWriter, "%v\t%v\t%v\n", "NAME", "READY", "AGE")
		for _, sset := range ssets.Items {
			if !matchesPattern(sset.Name, pattern) {
				continue
			}
			writef(tabWriter, "%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
		}
		flushTabWriter(tabWriter)
	}
}

func printUnstructuredList(w io.Writer, resources *unstructured.UnstructuredList, pattern string) {
	if len(resources.Items) == 0 {
		return
	}

	kind := resources.GetKind()
	if kind == "" {
		kind = "Resources"
	}

	writef(w, "\n%s\n--------------\n", kind)
	tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	writef(tabWriter, "%v\t%v\n", "NAMESPACE", "NAME")
	for _, resource := range resources.Items {
		if !matchesPattern(resource.GetName(), pattern) {
			continue
		}
		writef(tabWriter, "%v\t%v\n", resource.GetNamespace(), resource.GetName())
	}
	flushTabWriter(tabWriter)
}

func Printer(w io.Writer, resource interface{}, pattern string) {
	switch typedResource := resource.(type) {
	case *v1.PodList:
		printPodDetails(w, typedResource, pattern)
	case *v1.ConfigMapList:
		printConfigMaps(w, typedResource, pattern)
	case *v1.EndpointsList:
		printEndpoints(w, typedResource, pattern)
	case *v1.EventList:
		printEvents(w, typedResource, pattern)
	case *v1.LimitRangeList:
		printLimitRanges(w, typedResource, pattern)
	case *v1.NamespaceList:
		printNamespaces(w, typedResource, pattern)
	case *v1.PersistentVolumeList:
		printPVs(w, typedResource, pattern)
	case *v1.PersistentVolumeClaimList:
		printPVCs(w, typedResource, pattern)
	case *v1.PodTemplateList:
		printPodTemplates(w, typedResource, pattern)
	case *v1.ResourceQuotaList:
		printResourceQuotas(w, typedResource, pattern)
	case *v1.SecretList:
		printSecrets(w, typedResource, pattern)
	case *v1.ServiceList:
		printServices(w, typedResource, pattern)
	case *v1.ServiceAccountList:
		printServiceAccounts(w, typedResource, pattern)

		// these will be from the appsV1
	case *appsv1.DaemonSetList:
		printDaemonSets(w, typedResource, pattern)
	case *appsv1.DeploymentList:
		printDeployments(w, typedResource, pattern)
	case *appsv1.ReplicaSetList:
		printReplicaSets(w, typedResource, pattern)
	case *appsv1.StatefulSetList:
		printStateFulSets(w, typedResource, pattern)
	case *unstructured.UnstructuredList:
		printUnstructuredList(w, typedResource, pattern)
	}
}
