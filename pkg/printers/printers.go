package printers

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	v1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
)

func printPodDetails(pods *v1.PodList, resName string) {
	if len(pods.Items) > 0 {
		fmt.Printf("\nPods\n----\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "NAME", "READY", "STATUS", "RESTARTS")

		for _, pod := range pods.Items {
			if resName != "" {
				if strings.Contains(pod.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", pod.Name, "", pod.Status.Phase, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", pod.Name, "", pod.Status.Phase, "")
			}
		}
		w.Flush()
	}
}
func printPodTemplates(podTemplates *v1.PodTemplateList, resName string) {
	if len(podTemplates.Items) > 0 {
		fmt.Printf("\nPodTemplates\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\n", "NAME")
		for _, podTemplate := range podTemplates.Items {
			if resName != "" {
				if strings.Contains(podTemplate.Name, resName) {
					fmt.Fprintf(w, "%v\n", podTemplate.Name)
				}
			} else {
				fmt.Fprintf(w, "%v\n", podTemplate.Name)
			}
		}
		w.Flush()
	}
}
func printConfigMaps(cms *v1.ConfigMapList, resName string) {
	if len(cms.Items) > 0 {
		fmt.Printf("\nConfigMaps\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "NAME", "DATA", "AGE")
		for _, configMap := range cms.Items {
			if resName != "" {
				if strings.Contains(configMap.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\n", configMap.Name, len(configMap.Data), "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\n", configMap.Name, len(configMap.Data), "")
			}
		}
		w.Flush()
	}
}
func printEndpoints(endPoints *v1.EndpointsList, resName string) {
	if len(endPoints.Items) > 0 {
		fmt.Printf("\nEndpoints\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "NAME", "ENDPOINTS", "AGE")
		for _, endpoint := range endPoints.Items {
			if resName != "" {
				if strings.Contains(endpoint.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\n", endpoint.Name, "", "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\n", endpoint.Name, "", "")
			}
		}
		w.Flush()
	}
}
func printEvents(events *v1.EventList, resName string) {
	if len(events.Items) > 0 {
		fmt.Printf("\nEvents\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE")
		for _, event := range events.Items {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", event.Namespace, "", event.Type, "", event.InvolvedObject.Kind+"/"+event.InvolvedObject.Name, event.Message)
		}
		w.Flush()
	}
}
func printLimitRanges(limitRanges *v1.LimitRangeList, resName string) {
	if len(limitRanges.Items) > 0 {
		fmt.Printf("\nLimitRanges\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\n", "NAME", "CREATED AT")
		for _, limitRange := range limitRanges.Items {
			if resName != "" {
				if strings.Contains(limitRange.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\n", limitRange.Name, limitRange.CreationTimestamp)
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\n", limitRange.Name, limitRange.CreationTimestamp)
			}
		}
		w.Flush()
	}
}
func printNamespaces(namespaces *v1.NamespaceList, resName string) {
	if len(namespaces.Items) > 0 {
		fmt.Printf("\nNamespaces\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "NAME", "STATUS", "AGE")
		for _, namespace := range namespaces.Items {
			if resName != "" {
				if strings.Contains(namespace.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\n", namespace.Name, namespace.Status, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\n", namespace.Name, namespace.Status, "")
			}
		}
		w.Flush()
	}
}
func printPVs(pvs *v1.PersistentVolumeList, resName string) {
	if len(pvs.Items) > 0 {
		fmt.Printf("\nPersistentVolumes\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGECLASS", "REASON", "AGE")

		for _, pv := range pvs.Items {
			if resName != "" {
				if strings.Contains(pv.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pv.Name, "", pv.Spec.AccessModes, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.Spec.ClaimRef.Namespace+"/"+pv.Spec.ClaimRef.Name, pv.Spec.StorageClassName, pv.Status.Reason, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pv.Name, "", pv.Spec.AccessModes, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.Spec.ClaimRef.Namespace+"/"+pv.Spec.ClaimRef.Name, pv.Spec.StorageClassName, pv.Status.Reason, "")
			}
		}
		w.Flush()
	}
}
func printPVCs(pvcs *v1.PersistentVolumeClaimList, resName string) {
	if len(pvcs.Items) > 0 {
		fmt.Printf("\nPersistentVolumeClaims\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS MODES", "STORAGECLASS", "AGE")
		for _, pvc := range pvcs.Items {
			if resName != "" {
				if strings.Contains(pvc.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pvc.Name, pvc.Status, "", pvc.Status.Capacity.Cpu(), pvc.Spec.AccessModes, pvc.Spec.StorageClassName, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", pvc.Name, pvc.Status, "", pvc.Status.Capacity.Cpu(), pvc.Spec.AccessModes, pvc.Spec.StorageClassName, "")
			}

		}
		w.Flush()
	}
}

func printResourceQuotas(resQuotas *v1.ResourceQuotaList, resName string) {
	if len(resQuotas.Items) > 0 {
		fmt.Printf("\nResourceQuotas\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\n", "NAME", "CREATED AT")
		for _, resQ := range resQuotas.Items {
			if resName != "" {
				if strings.Contains(resQ.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\n", resQ.Name, resQ.CreationTimestamp)
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\n", resQ.Name, resQ.CreationTimestamp)
			}
		}
		w.Flush()
	}
}
func printSecrets(secrets *v1.SecretList, resName string) {
	if len(secrets.Items) > 0 {
		fmt.Printf("\nSecrets\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", "NAME", "TYPE", "DATA", "AGE")
		for _, secret := range secrets.Items {
			if resName != "" {
				if strings.Contains(secret.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", secret.Name, secret.Type, len(secret.Data), "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", secret.Name, secret.Type, len(secret.Data), "")
			}
		}
		w.Flush()
	}
}
func printServices(services *v1.ServiceList, resName string) {
	if len(services.Items) > 0 {
		fmt.Printf("\nServices\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE")

		for _, service := range services.Items {
			if resName != "" {
				if strings.Contains(service.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", service.Name, service.Spec.Type, service.Spec.ClusterIP, service.Spec.ExternalIPs, service.Spec.Ports, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", service.Name, service.Spec.Type, service.Spec.ClusterIP, service.Spec.ExternalIPs, service.Spec.Ports, "")
			}
		}
		w.Flush()
	}
}
func printServiceAccounts(serviceAccs *v1.ServiceAccountList, resName string) {
	if len(serviceAccs.Items) > 0 {
		fmt.Printf("\nServiceAccounts\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "NAME", "SECRETS", "AGE")
		for _, serviceAcc := range serviceAccs.Items {
			if resName != "" {
				if strings.Contains(serviceAcc.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\n", serviceAcc.Name, len(serviceAcc.Secrets), "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\n", serviceAcc.Name, len(serviceAcc.Secrets), "")
			}
		}
		w.Flush()
	}
}
func printDaemonSets(daemonsets *appsv1.DaemonSetList, resName string) {
	if len(daemonsets.Items) > 0 {
		fmt.Printf("\nDaemonSets\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", "NAMESPACE", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE")
		for _, ds := range daemonsets.Items {
			if resName != "" {
				if strings.Contains(ds.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", ds.Namespace, ds.Name, ds.Status.DesiredNumberScheduled, ds.Status.CurrentNumberScheduled, ds.Status.NumberReady, "", ds.Status.NumberAvailable, ds.Spec.Template.Spec.NodeSelector, "")
			}
		}
		w.Flush()
	}
}
func printDeployments(deployments *appsv1.DeploymentList, resName string) {
	if len(deployments.Items) > 0 {
		fmt.Printf("\nDeployments\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
		for _, deployment := range deployments.Items {
			if resName != "" {
				if strings.Contains(deployment.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", deployment.Name, deployment.Status.ReadyReplicas, "", deployment.Status.AvailableReplicas, "")
			}
		}
		w.Flush()
	}
}
func printReplicaSets(rsets *appsv1.ReplicaSetList, resName string) {
	if len(rsets.Items) > 0 {
		fmt.Printf("\nReplicaSets\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", "NAME", "DESIRED", "CURRENT", "READY", "AGE")
		for _, rs := range rsets.Items {
			if resName != "" {
				if strings.Contains(rs.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", rs.Name, "", "", "", "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", rs.Name, "", "", "", "")
			}
		}
		w.Flush()
	}
}
func printStateFulSets(ssets *appsv1.StatefulSetList, resName string) {
	if len(ssets.Items) > 0 {
		fmt.Printf("\nStatefulSets\n--------------\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%v\t%v\t%v\n", "NAME", "READY", "AGE")
		for _, sset := range ssets.Items {
			if resName != "" {
				if strings.Contains(sset.Name, resName) {
					fmt.Fprintf(w, "%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
				}
			} else {
				fmt.Fprintf(w, "%v\t%v\t%v\n", sset.Name, sset.Status.ReadyReplicas, "")
			}
		}
		w.Flush()
	}
}

func Printer(resource interface{}, resName string) {
	switch typedResource := resource.(type) {
	case *v1.PodList:
		pods := typedResource
		printPodDetails(pods, resName)
	case *v1.ConfigMapList:
		cms := typedResource
		printConfigMaps(cms, resName)
	case *v1.EndpointsList:
		endPoints := typedResource
		printEndpoints(endPoints, resName)
	case *v1.EventList:
		events := typedResource
		printEvents(events, resName)
	case *v1.LimitRangeList:
		limitRanges := typedResource
		printLimitRanges(limitRanges, resName)
	case *v1.NamespaceList:
		namespaces := typedResource
		printNamespaces(namespaces, resName)
	case *v1.PersistentVolumeList:
		pvs := typedResource
		printPVs(pvs, resName)
	case *v1.PersistentVolumeClaimList:
		pvcs := typedResource
		printPVCs(pvcs, resName)
	case *v1.PodTemplateList:
		podTemplates := typedResource
		printPodTemplates(podTemplates, resName)
	case *v1.ResourceQuotaList:
		resQuotas := typedResource
		printResourceQuotas(resQuotas, resName)
	case *v1.SecretList:
		secrets := typedResource
		printSecrets(secrets, resName)
	case *v1.ServiceList:
		services := typedResource
		printServices(services, resName)
	case *v1.ServiceAccountList:
		serviceAccs := typedResource
		printServiceAccounts(serviceAccs, resName)

		// these will be from the appsV1
	case *appsv1.DaemonSetList:
		daemonsets := typedResource
		printDaemonSets(daemonsets, resName)
	case *appsv1.DeploymentList:
		deployments := typedResource
		printDeployments(deployments, resName)
	case *appsv1.ReplicaSetList:
		rsets := typedResource
		printReplicaSets(rsets, resName)
	case *appsv1.StatefulSetList:
		ssets := typedResource
		printStateFulSets(ssets, resName)
	}
}
