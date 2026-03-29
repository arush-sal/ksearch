package util

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

type ResourceMeta struct {
	Kind, Resource, APIGroup, APIVersion string
	Namespaced                           bool
}

func Discover(dc discovery.DiscoveryInterface, kinds string) ([]ResourceMeta, error) {
	_, lists, err := dc.ServerGroupsAndResources()
	if err != nil && lists == nil {
		return nil, err
	}
	if err != nil && lists != nil {
		log.Warnf("partial discovery failure: %v", err)
	}

	filter := parseKindsFilter(kinds)
	resources := make([]ResourceMeta, 0)
	for _, list := range lists {
		if list == nil {
			continue
		}

		gv, parseErr := schema.ParseGroupVersion(list.GroupVersion)
		if parseErr != nil {
			return nil, fmt.Errorf("parse group version %q: %w", list.GroupVersion, parseErr)
		}

		for _, resource := range list.APIResources {
			if !hasVerb(resource.Verbs, "list") {
				continue
			}
			if _, ok := canonicalResourceName(resource.Kind, resource.Name); !ok {
				continue
			}
			if !matchesKindsFilter(filter, resource.Kind, resource.Name) {
				continue
			}

			resources = append(resources, ResourceMeta{
				Kind:       resource.Kind,
				Resource:   resource.Name,
				APIGroup:   gv.Group,
				APIVersion: gv.Version,
				Namespaced: resource.Namespaced,
			})
		}
	}

	return resources, nil
}

func parseKindsFilter(kinds string) map[string]bool {
	filter := make(map[string]bool)
	if strings.TrimSpace(kinds) == "" {
		return filter
	}

	for _, kind := range strings.Split(kinds, ",") {
		kind = strings.ToLower(strings.TrimSpace(kind))
		if kind == "" {
			continue
		}
		filter[kind] = true
	}

	return filter
}

func hasVerb(verbs []string, target string) bool {
	for _, verb := range verbs {
		if strings.EqualFold(verb, target) {
			return true
		}
	}

	return false
}

func matchesKindsFilter(filter map[string]bool, kind, resource string) bool {
	if len(filter) == 0 {
		return true
	}

	return filter[strings.ToLower(kind)] || filter[strings.ToLower(resource)]
}

func canonicalResourceName(kind, resource string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(resource)) {
	case "pods":
		return "pods", true
	case "configmaps":
		return "configmaps", true
	case "endpoints":
		return "endpoints", true
	case "events":
		return "events", true
	case "limitranges":
		return "limitranges", true
	case "namespaces":
		return "namespaces", true
	case "persistentvolumes":
		return "persistentvolumes", true
	case "persistentvolumeclaims":
		return "persistentvolumeclaims", true
	case "podtemplates":
		return "podtemplates", true
	case "resourcequotas":
		return "resourcequotas", true
	case "secrets":
		return "secrets", true
	case "services":
		return "services", true
	case "serviceaccounts":
		return "serviceaccounts", true
	case "daemonsets":
		return "daemonsets", true
	case "deployments":
		return "deployments", true
	case "replicasets":
		return "replicasets", true
	case "statefulsets":
		return "statefulsets", true
	}

	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "pod", "pods":
		return "pods", true
	case "configmap", "configmaps":
		return "configmaps", true
	case "endpoint", "endpoints":
		return "endpoints", true
	case "event", "events":
		return "events", true
	case "limitrange", "limitranges":
		return "limitranges", true
	case "namespace", "namespaces":
		return "namespaces", true
	case "persistentvolume", "persistentvolumes":
		return "persistentvolumes", true
	case "persistentvolumeclaim", "persistentvolumeclaims":
		return "persistentvolumeclaims", true
	case "podtemplate", "podtemplates":
		return "podtemplates", true
	case "resourcequota", "resourcequotas":
		return "resourcequotas", true
	case "secret", "secrets":
		return "secrets", true
	case "service", "services":
		return "services", true
	case "serviceaccount", "serviceaccounts":
		return "serviceaccounts", true
	case "daemonset", "daemonsets":
		return "daemonsets", true
	case "deployment", "deployments":
		return "deployments", true
	case "replicaset", "replicasets":
		return "replicasets", true
	case "statefulset", "statefulsets":
		return "statefulsets", true
	}

	return "", false
}
