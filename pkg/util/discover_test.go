package util

import (
	"errors"
	"testing"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/openapi"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
)

func TestDiscover_AllWhenEmpty(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 4 {
		t.Fatalf("expected 4 listable resources, got %d", len(resources))
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

func TestDiscover_ReturnsUnsupportedListableKinds(t *testing.T) {
	t.Parallel()

	resources, err := Discover(newFakeDiscovery(), "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	foundWidget := false
	for _, resource := range resources {
		if resource.Kind == "Widget" {
			foundWidget = true
		}
	}

	if !foundWidget {
		t.Fatal("expected listable widget resource to be discovered")
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

	if len(resources) != 4 {
		t.Fatalf("expected partial discovery to return listable resources, got %d", len(resources))
	}
}

func TestDiscover_DeduplicatesLogicalResourcesAcrossVersions(t *testing.T) {
	t.Parallel()

	dc := newPreferredDiscovery(
		[]*metav1.APIGroup{
			{
				Name: "example.com",
				Versions: []metav1.GroupVersionForDiscovery{
					{GroupVersion: "example.com/v1beta1", Version: "v1beta1"},
					{GroupVersion: "example.com/v1", Version: "v1"},
				},
				PreferredVersion: metav1.GroupVersionForDiscovery{GroupVersion: "example.com/v1", Version: "v1"},
			},
		},
		[]*metav1.APIResourceList{
			{
				GroupVersion: "example.com/v1beta1",
				APIResources: []metav1.APIResource{
					{Kind: "Widget", Name: "widgets", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
				},
			},
			{
				GroupVersion: "example.com/v1",
				APIResources: []metav1.APIResource{
					{Kind: "Widget", Name: "widgets", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
				},
			},
		},
	)

	resources, err := Discover(dc, "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 deduplicated resource, got %d", len(resources))
	}

	if resources[0].APIVersion != "v1" {
		t.Fatalf("expected preferred version to be kept, got %q", resources[0].APIVersion)
	}
}

func TestDiscover_DeduplicatesBuiltinsAcrossGroups(t *testing.T) {
	t.Parallel()

	dc := &discoveryfake.FakeDiscovery{
		Fake: &k8stesting.Fake{
			Resources: []*metav1.APIResourceList{
				{
					GroupVersion: "v1",
					APIResources: []metav1.APIResource{
						{Kind: "Event", Name: "events", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
					},
				},
				{
					GroupVersion: "events.k8s.io/v1",
					APIResources: []metav1.APIResource{
						{Kind: "Event", Name: "events", Namespaced: true, Verbs: metav1.Verbs{"get", "list"}},
					},
				},
			},
		},
	}

	resources, err := Discover(dc, "")
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 deduplicated event resource, got %d", len(resources))
	}

	if resources[0].APIGroup != "" || resources[0].APIVersion != "v1" {
		t.Fatalf("expected core/v1 event to be kept, got %q/%q", resources[0].APIGroup, resources[0].APIVersion)
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

type preferredDiscovery struct {
	*discoveryfake.FakeDiscovery
	groups []*metav1.APIGroup
}

func newPreferredDiscovery(groups []*metav1.APIGroup, resources []*metav1.APIResourceList) *preferredDiscovery {
	return &preferredDiscovery{
		FakeDiscovery: &discoveryfake.FakeDiscovery{
			Fake: &k8stesting.Fake{
				Resources: resources,
			},
		},
		groups: groups,
	}
}

func (d *preferredDiscovery) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return d.groups, d.Resources, nil
}

func (d *preferredDiscovery) ServerGroups() (*metav1.APIGroupList, error) {
	return &metav1.APIGroupList{}, nil
}

func (d *preferredDiscovery) ServerResourcesForGroupVersion(string) (*metav1.APIResourceList, error) {
	return nil, nil
}

func (d *preferredDiscovery) ServerResources() ([]*metav1.APIResourceList, error) {
	return d.Resources, nil
}

func (d *preferredDiscovery) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (d *preferredDiscovery) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (d *preferredDiscovery) ServerVersion() (*version.Info, error) {
	return &version.Info{}, nil
}

func (d *preferredDiscovery) OpenAPISchema() (*openapi_v2.Document, error) {
	return &openapi_v2.Document{}, nil
}

func (d *preferredDiscovery) OpenAPIV3() openapi.Client {
	return nil
}

func (d *preferredDiscovery) RESTClient() rest.Interface {
	return nil
}
