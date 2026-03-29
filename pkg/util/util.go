package util

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type FetchResult struct {
	Kind     string
	Resource interface{}
}

// Getter the
// This should be go routine ready. Such that getter can be called via goroutines and over a channel the value can be passed to a switch type through with the respective printer can be called.
func Getter(namespace string, clientset kubernetes.Interface, restConfig *rest.Config, resources []ResourceMeta, c chan FetchResult) {
	defer close(c)
	ctx := context.Background()

	for _, meta := range resources {
		var (
			err  error
			list interface{}
		)

		canonical := canonicalResourceName(meta.Kind, meta.Resource)
		if canonical != "" && endpointMatchesTyped(meta, canonical) {
			switch canonical {
			case "pods":
				list, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "configmaps":
				list, err = clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "endpoints":
				list, err = clientset.CoreV1().Endpoints(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "events":
				list, err = clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "limitranges":
				list, err = clientset.CoreV1().LimitRanges(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "namespaces":
				list, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "persistentvolumes":
				list, err = clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "persistentvolumeclaims":
				list, err = clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "podtemplates":
				list, err = clientset.CoreV1().PodTemplates(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "resourcequotas":
				list, err = clientset.CoreV1().ResourceQuotas(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "secrets":
				list, err = clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "services":
				list, err = clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "serviceaccounts":
				list, err = clientset.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "daemonsets":
				list, err = clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "deployments":
				list, err = clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "replicasets":
				list, err = clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			case "statefulsets":
				list, err = clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
				handleError(err, meta.Kind)
			}
		} else {
			var dynamicList interface{}
			unstructuredList, listErr := listUnstructuredResource(ctx, restConfig, namespace, meta)
			if unstructuredList != nil {
				dynamicList = unstructuredList
			}
			list, err = dynamicList, listErr
			handleError(err, meta.Kind)
		}

		c <- FetchResult{Kind: meta.Kind, Resource: list}
	}
}

func handleError(err error, r string) {
	if err != nil {
		log.Errorf("There was an error getting the %s from clientset", r)
	}
}
