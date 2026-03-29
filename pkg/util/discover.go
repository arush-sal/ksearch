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
	seen := make(map[string]bool)
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
			if !matchesKindsFilter(filter, resource.Kind, resource.Name) {
				continue
			}

			logicalKey := discoveryDedupKey(resource.Kind, resource.Name, gv.Group, resource.Namespaced)
			if seen[logicalKey] {
				continue
			}
			seen[logicalKey] = true

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

func canonicalResourceName(kind, resource string) string {
	switch strings.ToLower(strings.TrimSpace(resource)) {
	case "pods":
		return "pods"
	case "configmaps":
		return "configmaps"
	case "endpoints":
		return "endpoints"
	case "events":
		return "events"
	case "limitranges":
		return "limitranges"
	case "namespaces":
		return "namespaces"
	case "persistentvolumes":
		return "persistentvolumes"
	case "persistentvolumeclaims":
		return "persistentvolumeclaims"
	case "podtemplates":
		return "podtemplates"
	case "resourcequotas":
		return "resourcequotas"
	case "secrets":
		return "secrets"
	case "services":
		return "services"
	case "serviceaccounts":
		return "serviceaccounts"
	case "daemonsets":
		return "daemonsets"
	case "deployments":
		return "deployments"
	case "replicasets":
		return "replicasets"
	case "statefulsets":
		return "statefulsets"
	}

	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "pod", "pods":
		return "pods"
	case "configmap", "configmaps":
		return "configmaps"
	case "endpoint", "endpoints":
		return "endpoints"
	case "event", "events":
		return "events"
	case "limitrange", "limitranges":
		return "limitranges"
	case "namespace", "namespaces":
		return "namespaces"
	case "persistentvolume", "persistentvolumes":
		return "persistentvolumes"
	case "persistentvolumeclaim", "persistentvolumeclaims":
		return "persistentvolumeclaims"
	case "podtemplate", "podtemplates":
		return "podtemplates"
	case "resourcequota", "resourcequotas":
		return "resourcequotas"
	case "secret", "secrets":
		return "secrets"
	case "service", "services":
		return "services"
	case "serviceaccount", "serviceaccounts":
		return "serviceaccounts"
	case "daemonset", "daemonsets":
		return "daemonsets"
	case "deployment", "deployments":
		return "deployments"
	case "replicaset", "replicasets":
		return "replicasets"
	case "statefulset", "statefulsets":
		return "statefulsets"
	}

	return ""
}

func discoveryDedupKey(kind, resource, group string, namespaced bool) string {
	if canonical := canonicalResourceName(kind, resource); canonical != "" {
		return fmt.Sprintf("canonical/%s/%t", canonical, namespaced)
	}

	return fmt.Sprintf("%s/%s/%t", group, resource, namespaced)
}
