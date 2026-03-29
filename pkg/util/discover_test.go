package util

import (
	"errors"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	discoveryfake "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestDiscover_AllWhenEmpty(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 3 {
		t.Fatalf("expected 3 listable resources, got %d", len(resources))
	}
}

func TestDiscover_FilterByKinds(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "ConfigMaps")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	if resources[0].Kind != "ConfigMap" {
		t.Fatalf("expected ConfigMap, got %q", resources[0].Kind)
	}
}

func TestDiscover_SkipsNonListable(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	for _, resource := range resources {
		if resource.Kind == "Lease" {
			t.Fatalf("expected non-listable resource to be skipped: %#v", resource)
		}
	}
}

func TestDiscover_SkipsUnsupportedListableKinds(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	for _, resource := range resources {
		if resource.Kind == "Widget" {
			t.Fatalf("expected unsupported listable resource to be skipped: %#v", resource)
		}
	}
}

func TestDiscover_CaseInsensitive(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "configmaps")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	if resources[0].Kind != "ConfigMap" {
		t.Fatalf("expected ConfigMap, got %q", resources[0].Kind)
	}
}

func TestDiscover_MultipleKinds(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "Pods,ConfigMaps")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	if resources[0].Kind != "Pod" {
		t.Fatalf("expected first resource to be Pod, got %q", resources[0].Kind)
	}

	if resources[1].Kind != "ConfigMap" {
		t.Fatalf("expected second resource to be ConfigMap, got %q", resources[1].Kind)
	}
}

func TestDiscover_AcceptsResourceNameForms(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "pods,secrets")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(resources))
	}

	if resources[0].Kind != "Pod" {
		t.Fatalf("expected first resource to be Pod, got %q", resources[0].Kind)
	}

	if resources[1].Kind != "Secret" {
		t.Fatalf("expected second resource to be Secret, got %q", resources[1].Kind)
	}
}

func TestDiscover_PartialFailureContinues(t *testing.T) {
	t.Parallel()

	dc := newFakeDiscovery()
	dc.PrependReactor("get", "resource", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("partial discovery failure")
	})

	resources, err := Discover(dc, "")
	if err != nil {
		t.Fatalf("expected partial failure to continue, got error: %v", err)
	}

	if len(resources) != 3 {
		t.Fatalf("expected partial discovery to return listable resources, got %d", len(resources))
	}
}

func newFakeDiscovery() *discoveryfake.FakeDiscovery {
	return &discoveryfake.FakeDiscovery{
		Fake: &k8stesting.Fake{
			Resources: []*metav1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []metav1.APIResource{
						{Kind: "Pod", Name: "pods", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
						{Kind: "ConfigMap", Name: "configmaps", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
						{Kind: "Secret", Name: "secrets", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
						{Kind: "Widget", Name: "widgets", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
						{Kind: "Lease", Name: "leases", Namespaced: true, Verbs: metav1.Verbs{"get"}},
					},
				},
			},
		},
	}
}
