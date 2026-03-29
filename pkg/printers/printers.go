package printers

import (
	"bytes"
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

func writeSectionIfMatched(w io.Writer, title, divider string, render func(*tabwriter.Writer) int) {
	var body bytes.Buffer
	tabWriter := tabwriter.NewWriter(&body, 0, 0, 1, ' ', 0)
	rows := render(tabWriter)
	if rows == 0 {
		return
	}

	flushTabWriter(tabWriter)
	writef(w, "\n%s\n%s\n", title, divider)
	_, _ = io.Copy(w, &body)
}

func printPodDetails(w io.Writer, pods *v1.PodList, pattern string) {
	if len(pods.Items) > 0 {
		writeSectionIfMatched(w, "Pods", "----", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "READY", "STATUS", "RESTARTS")
			for _, pod := range pods.Items {
				if !matchesPattern(pod.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\n", pod.Name, "", pod.Status.Phase, "")
			}
			return rows
		})
	}
}
func printPodTemplates(w io.Writer, podTemplates *v1.PodTemplateList, pattern string) {
	if len(podTemplates.Items) > 0 {
		writeSectionIfMatched(w, "PodTemplates", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\n", "NAME")
			for _, podTemplate := range podTemplates.Items {
				if !matchesPattern(podTemplate.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\n", podTemplate.Name)
			}
			return rows
		})
	}
}
func printConfigMaps(w io.Writer, cms *v1.ConfigMapList, pattern string) {
	if len(cms.Items) > 0 {
		writeSectionIfMatched(w, "ConfigMaps", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\n", "NAME", "DATA", "AGE")
			for _, configMap := range cms.Items {
				if !matchesPattern(configMap.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\n", configMap.Name, len(configMap.Data), "")
			}
			return rows
		})
	}
}
func printEndpoints(w io.Writer, endPoints *v1.EndpointsList, pattern string) {
	if len(endPoints.Items) > 0 {
		writeSectionIfMatched(w, "Endpoints", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\n", "NAME", "ENDPOINTS", "AGE")
			for _, endpoint := range endPoints.Items {
				if !matchesPattern(endpoint.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\n", endpoint.Name, "", "")
			}
			return rows
		})
	}
}
func printEvents(w io.Writer, events *v1.EventList, pattern string) {
	if len(events.Items) > 0 {
		writeSectionIfMatched(w, "Events", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE")
			for _, event := range events.Items {
				if !matchesPattern(event.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", event.Namespace, "", event.Type, "", event.InvolvedObject.Kind+"/"+event.InvolvedObject.Name, event.Message)
			}
			return rows
		})
	}
}
func printLimitRanges(w io.Writer, limitRanges *v1.LimitRangeList, pattern string) {
	if len(limitRanges.Items) > 0 {
		writeSectionIfMatched(w, "LimitRanges", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
			for _, limitRange := range limitRanges.Items {
				if !matchesPattern(limitRange.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\n", limitRange.Name, limitRange.CreationTimestamp)
			}
			return rows
		})
	}
}
func printNamespaces(w io.Writer, namespaces *v1.NamespaceList, pattern string) {
	if len(namespaces.Items) > 0 {
		writeSectionIfMatched(w, "Namespaces", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\n", "NAME", "STATUS", "AGE")
			for _, namespace := range namespaces.Items {
				if !matchesPattern(namespace.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\n", namespace.Name, namespace.Status, "")
			}
			return rows
		})
	}
}
func printPVs(w io.Writer, pvs *v1.PersistentVolumeList, pattern string) {
	if len(pvs.Items) > 0 {
		writeSectionIfMatched(w, "PersistentVolumes", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGECLASS", "REASON", "AGE")
			for _, pv := range pvs.Items {
				if !matchesPattern(pv.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pv.Name, "", pv.Spec.AccessModes, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.Spec.ClaimRef.Namespace+"/"+pv.Spec.ClaimRef.Name, pv.Spec.StorageClassName, pv.Status.Reason, "")
			}
			return rows
		})
	}
}
func printPVCs(w io.Writer, pvcs *v1.PersistentVolumeClaimList, pattern string) {
	if len(pvcs.Items) > 0 {
		writeSectionIfMatched(w, "PersistentVolumeClaims", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS MODES", "STORAGECLASS", "AGE")
			for _, pvc := range pvcs.Items {
				if !matchesPattern(pvc.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pvc.Name, pvc.Status, "", pvc.Status.Capacity.Cpu(), pvc.Spec.AccessModes, pvc.Spec.StorageClassName, "")
			}
			return rows
		})
	}
}

func printResourceQuotas(w io.Writer, resQuotas *v1.ResourceQuotaList, pattern string) {
	if len(resQuotas.Items) > 0 {
		writeSectionIfMatched(w, "ResourceQuotas", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\n", "NAME", "CREATED AT")
			for _, resQ := range resQuotas.Items {
				if !matchesPattern(resQ.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\n", resQ.Name, resQ.CreationTimestamp)
			}
			return rows
		})
	}
}
func printSecrets(w io.Writer, secrets *v1.SecretList, pattern string) {
	if len(secrets.Items) > 0 {
		writeSectionIfMatched(w, "Secrets", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\n", "NAME", "TYPE", "DATA", "AGE")
			for _, secret := range secrets.Items {
				if !matchesPattern(secret.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\n", secret.Name, secret.Type, len(secret.Data), "")
			}
			return rows
		})
	}
}
func printServices(w io.Writer, services *v1.ServiceList, pattern string) {
	if len(services.Items) > 0 {
		writeSectionIfMatched(w, "Services", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE")
			for _, service := range services.Items {
				if !matchesPattern(service.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\n", service.Name, service.Spec.Type, service.Spec.ClusterIP, service.Spec.ExternalIPs, service.Spec.Ports, "")
			}
			return rows
		})
	}
}
func printServiceAccounts(w io.Writer, serviceAccs *v1.ServiceAccountList, pattern string) {
	if len(serviceAccs.Items) > 0 {
		writeSectionIfMatched(w, "ServiceAccounts", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\n", "NAME", "SECRETS", "AGE")
			for _, serviceAcc := range serviceAccs.Items {
				if !matchesPattern(serviceAcc.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\n", serviceAcc.Name, len(serviceAcc.Secrets), "")
			}
			return rows
		})
	}
}
func printDaemonSets(w io.Writer, daemonsets *appsv1.DaemonSetList, pattern string) {
	if len(daemonsets.Items) > 0 {
		writeSectionIfMatched(w, "DaemonSets", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE")
			for _, ds := range daemonsets.Items {
				if !matchesPattern(ds.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
			}
			return rows
		})
	}
}
func printDeployments(w io.Writer, deployments *appsv1.DeploymentList, pattern string) {
	if len(deployments.Items) > 0 {
		writeSectionIfMatched(w, "Deployments", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
			for _, deployment := range deployments.Items {
				if !matchesPattern(deployment.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
			}
			return rows
		})
	}
}
func printReplicaSets(w io.Writer, rsets *appsv1.ReplicaSetList, pattern string) {
	if len(rsets.Items) > 0 {
		writeSectionIfMatched(w, "ReplicaSets", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", "NAME", "DESIRED", "CURRENT", "READY", "AGE")
			for _, rs := range rsets.Items {
				if !matchesPattern(rs.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\t%v\t%v\n", rs.Name, "", "", "", "")
			}
			return rows
		})
	}
}
func printStateFulSets(w io.Writer, ssets *appsv1.StatefulSetList, pattern string) {
	if len(ssets.Items) > 0 {
		writeSectionIfMatched(w, "StatefulSets", "--------------", func(tabWriter *tabwriter.Writer) int {
			rows := 0
			writef(tabWriter, "%v\t%v\t%v\n", "NAME", "READY", "AGE")
			for _, sset := range ssets.Items {
				if !matchesPattern(sset.Name, pattern) {
					continue
				}
				rows++
				writef(tabWriter, "%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
			}
			return rows
		})
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

	writeSectionIfMatched(w, kind, "--------------", func(tabWriter *tabwriter.Writer) int {
		rows := 0
		writef(tabWriter, "%v\t%v\n", "NAMESPACE", "NAME")
		for _, resource := range resources.Items {
			if !matchesPattern(resource.GetName(), pattern) {
				continue
			}
			rows++
			writef(tabWriter, "%v\t%v\n", resource.GetNamespace(), resource.GetName())
		}
		return rows
	})
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
