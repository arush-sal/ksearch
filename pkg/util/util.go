package util

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Getter the
// This should be go routine ready. Such that getter can be called via goroutines and over a channel the value can be passed to a switch type through with the respective printer can be called.
func Getter(namespace string, clientset kubernetes.Interface, resources []ResourceMeta, c chan interface{}) {
	defer close(c)
	ctx := context.Background()
	var err error
	var list interface{}

	for _, meta := range resources {
		switch meta.Kind {
		case "Pod", "Pods":
			list, err = clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "ConfigMap", "ConfigMaps":
			list, err = clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Endpoint", "Endpoints":
			list, err = clientset.CoreV1().Endpoints(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Event", "Events":
			list, err = clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "LimitRange", "LimitRanges":
			list, err = clientset.CoreV1().LimitRanges(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Namespace", "Namespaces":
			list, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "PersistentVolume", "PersistentVolumes":
			list, err = clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "PersistentVolumeClaim", "PersistentVolumeClaims":
			list, err = clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "PodTemplate", "PodTemplates":
			list, err = clientset.CoreV1().PodTemplates(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "ResourceQuota", "ResourceQuotas":
			list, err = clientset.CoreV1().ResourceQuotas(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Secret", "Secrets":
			list, err = clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Service", "Services":
			list, err = clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "ServiceAccount", "ServiceAccounts":
			list, err = clientset.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)

		// these will be from the AppsV1
		case "DaemonSet", "DaemonSets":
			list, err = clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "Deployment", "Deployments":
			list, err = clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "ReplicaSet", "ReplicaSets":
			list, err = clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		case "StatefulSet", "StatefulSets":
			list, err = clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
			handleError(err, meta.Kind)
		default:
			log.Debugf("kind %q not handled, skipping", meta.Kind)
			continue
		}

		if list != nil {
			c <- list
		}
	}
}

func handleError(err error, r string) {
	if err != nil {
		log.Errorf("There was an error getting the %s from clientset", r)
	}
}
