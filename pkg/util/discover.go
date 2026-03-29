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
			if len(filter) > 0 && !filter[strings.ToLower(resource.Kind)] {
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
