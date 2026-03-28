package printers

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	v1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
)

func matchesPattern(name, pattern string) bool {
	return pattern == "" || strings.Contains(name, pattern)
}

func printPodDetails(w io.Writer, pods *v1.PodList, pattern string) {
	if len(pods.Items) > 0 {
		fmt.Fprintf(w, "\nPods\n----\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "READY", "STATUS", "RESTARTS")

		for _, pod := range pods.Items {
			if !matchesPattern(pod.Name, pattern) {
				continue
			}

			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\n", pod.Name, "", pod.Status.Phase, "")
		}
		tabWriter.Flush()
	}
}
func printPodTemplates(w io.Writer, podTemplates *v1.PodTemplateList, pattern string) {
	if len(podTemplates.Items) > 0 {
		fmt.Fprintf(w, "\nPodTemplates\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\n", "NAME")
		for _, podTemplate := range podTemplates.Items {
			if !matchesPattern(podTemplate.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\n", podTemplate.Name)
		}
		tabWriter.Flush()
	}
}
func printConfigMaps(w io.Writer, cms *v1.ConfigMapList, pattern string) {
	if len(cms.Items) > 0 {
		fmt.Fprintf(w, "\nConfigMaps\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", "NAME", "DATA", "AGE")
		for _, configMap := range cms.Items {
			if !matchesPattern(configMap.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", configMap.Name, len(configMap.Data), "")
		}
		tabWriter.Flush()
	}
}
func printEndpoints(w io.Writer, endPoints *v1.EndpointsList, pattern string) {
	if len(endPoints.Items) > 0 {
		fmt.Fprintf(w, "\nEndpoints\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", "NAME", "ENDPOINTS", "AGE")
		for _, endpoint := range endPoints.Items {
			if !matchesPattern(endpoint.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", endpoint.Name, "", "")
		}
		tabWriter.Flush()
	}
}
func printEvents(w io.Writer, events *v1.EventList, pattern string) {
	if len(events.Items) > 0 {
		fmt.Fprintf(w, "\nEvents\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE")
		for _, event := range events.Items {
			if !matchesPattern(event.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", event.Namespace, "", event.Type, "", event.InvolvedObject.Kind+"/"+event.InvolvedObject.Name, event.Message)
		}
		tabWriter.Flush()
	}
}
func printLimitRanges(w io.Writer, limitRanges *v1.LimitRangeList, pattern string) {
	if len(limitRanges.Items) > 0 {
		fmt.Fprintf(w, "\nLimitRanges\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
		for _, limitRange := range limitRanges.Items {
			if !matchesPattern(limitRange.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\n", limitRange.Name, limitRange.CreationTimestamp)
		}
		tabWriter.Flush()
	}
}
func printNamespaces(w io.Writer, namespaces *v1.NamespaceList, pattern string) {
	if len(namespaces.Items) > 0 {
		fmt.Fprintf(w, "\nNamespaces\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", "NAME", "STATUS", "AGE")
		for _, namespace := range namespaces.Items {
			if !matchesPattern(namespace.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", namespace.Name, namespace.Status, "")
		}
		tabWriter.Flush()
	}
}
func printPVs(w io.Writer, pvs *v1.PersistentVolumeList, pattern string) {
	if len(pvs.Items) > 0 {
		fmt.Fprintf(w, "\nPersistentVolumes\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGECLASS", "REASON", "AGE")

		for _, pv := range pvs.Items {
			if !matchesPattern(pv.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pv.Name, "", pv.Spec.AccessModes, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.Spec.ClaimRef.Namespace+"/"+pv.Spec.ClaimRef.Name, pv.Spec.StorageClassName, pv.Status.Reason, "")
		}
		tabWriter.Flush()
	}
}
func printPVCs(w io.Writer, pvcs *v1.PersistentVolumeClaimList, pattern string) {
	if len(pvcs.Items) > 0 {
		fmt.Fprintf(w, "\nPersistentVolumeClaims\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS MODES", "STORAGECLASS", "AGE")
		for _, pvc := range pvcs.Items {
			if !matchesPattern(pvc.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pvc.Name, pvc.Status, "", pvc.Status.Capacity.Cpu(), pvc.Spec.AccessModes, pvc.Spec.StorageClassName, "")
		}
		tabWriter.Flush()
	}
}

func printResourceQuotas(w io.Writer, resQuotas *v1.ResourceQuotaList, pattern string) {
	if len(resQuotas.Items) > 0 {
		fmt.Fprintf(w, "\nResourceQuotas\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
		for _, resQ := range resQuotas.Items {
			if !matchesPattern(resQ.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\n", resQ.Name, resQ.CreationTimestamp)
		}
		tabWriter.Flush()
	}
}
func printSecrets(w io.Writer, secrets *v1.SecretList, pattern string) {
	if len(secrets.Items) > 0 {
		fmt.Fprintf(w, "\nSecrets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "TYPE", "DATA", "AGE")
		for _, secret := range secrets.Items {
			if !matchesPattern(secret.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\n", secret.Name, secret.Type, len(secret.Data), "")
		}
		tabWriter.Flush()
	}
}
func printServices(w io.Writer, services *v1.ServiceList, pattern string) {
	if len(services.Items) > 0 {
		fmt.Fprintf(w, "\nServices\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE")

		for _, service := range services.Items {
			if !matchesPattern(service.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", service.Name, service.Spec.Type, service.Spec.ClusterIP, service.Spec.ExternalIPs, service.Spec.Ports, "")
		}
		tabWriter.Flush()
	}
}
func printServiceAccounts(w io.Writer, serviceAccs *v1.ServiceAccountList, pattern string) {
	if len(serviceAccs.Items) > 0 {
		fmt.Fprintf(w, "\nServiceAccounts\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", "NAME", "SECRETS", "AGE")
		for _, serviceAcc := range serviceAccs.Items {
			if !matchesPattern(serviceAcc.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", serviceAcc.Name, len(serviceAcc.Secrets), "")
		}
		tabWriter.Flush()
	}
}
func printDaemonSets(w io.Writer, daemonsets *appsv1.DaemonSetList, pattern string) {
	if len(daemonsets.Items) > 0 {
		fmt.Fprintf(w, "\nDaemonSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE")
		for _, ds := range daemonsets.Items {
			if !matchesPattern(ds.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
		}
		tabWriter.Flush()
	}
}
func printDeployments(w io.Writer, deployments *appsv1.DeploymentList, pattern string) {
	if len(deployments.Items) > 0 {
		fmt.Fprintf(w, "\nDeployments\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
		for _, deployment := range deployments.Items {
			if !matchesPattern(deployment.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
		}
		tabWriter.Flush()
	}
}
func printReplicaSets(w io.Writer, rsets *appsv1.ReplicaSetList, pattern string) {
	if len(rsets.Items) > 0 {
		fmt.Fprintf(w, "\nReplicaSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "DESIRED", "CURRENT", "READY", "AGE")
		for _, rs := range rsets.Items {
			if !matchesPattern(rs.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\t%v\t%v\n", rs.Name, "", "", "", "")
		}
		tabWriter.Flush()
	}
}
func printStateFulSets(w io.Writer, ssets *appsv1.StatefulSetList, pattern string) {
	if len(ssets.Items) > 0 {
		fmt.Fprintf(w, "\nStatefulSets\n--------------\n")
		tabWriter := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", "NAME", "READY", "AGE")
		for _, sset := range ssets.Items {
			if !matchesPattern(sset.Name, pattern) {
				continue
			}
			fmt.Fprintf(tabWriter, "%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
		}
		tabWriter.Flush()
	}
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
	}
}
