package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestGetter_CustomKinds(t *testing.T) {
	t.Parallel()

	clientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-config",
			Namespace: "default",
		},
	})

	results := make(chan interface{})

	go Getter("default", clientset, nil, []ResourceMeta{
		{Kind: "ConfigMap", Resource: "configmaps", Namespaced: true},
	}, results)

	var received []interface{}
	for item := range results {
		received = append(received, item)
	}

	if len(received) != 1 {
		t.Fatalf("expected exactly one result, got %d", len(received))
	}

	configMaps, ok := received[0].(*v1.ConfigMapList)
	if !ok {
		t.Fatalf("expected *v1.ConfigMapList, got %T", received[0])
	}

	if len(configMaps.Items) != 1 {
		t.Fatalf("expected one configmap in result, got %d", len(configMaps.Items))
	}
}

func TestGetter_UnknownKind(t *testing.T) {
	t.Parallel()

	clientset := fake.NewSimpleClientset()
	results := make(chan interface{})

	go Getter("default", clientset, nil, []ResourceMeta{
		{Kind: "NonExistentKind", Resource: "nonexistentkinds", Namespaced: true},
	}, results)

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel to close")
	}
}

func TestGetter_ChannelAlwaysClosed(t *testing.T) {
	t.Parallel()

	clientset := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},
	})
	results := make(chan interface{})

	go Getter("default", clientset, nil, []ResourceMeta{
		{Kind: "Pod", Resource: "pods", Namespaced: true},
	}, results)

	done := make(chan struct{})
	go func() {
		for range results {
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel to close")
	}
}

func TestGetter_CustomResource(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/apis/example.com/v1/namespaces/default/widgets" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"apiVersion":"example.com/v1","kind":"WidgetList","metadata":{},"items":[{"apiVersion":"example.com/v1","kind":"Widget","metadata":{"name":"demo","namespace":"default"}}]}`))
	}))
	defer server.Close()

	clientset := fake.NewSimpleClientset()
	results := make(chan interface{})
	cfg := &rest.Config{Host: server.URL}

	go Getter("default", clientset, cfg, []ResourceMeta{
		{Kind: "Widget", Resource: "widgets", APIGroup: "example.com", APIVersion: "v1", Namespaced: true},
	}, results)

	select {
	case item, ok := <-results:
		if !ok {
			t.Fatal("expected custom resource list before channel close")
		}

		list, ok := item.(*unstructured.UnstructuredList)
		if !ok {
			t.Fatalf("expected *unstructured.UnstructuredList, got %T", item)
		}

		if list.GetKind() != "Widget" {
			t.Fatalf("expected list kind Widget, got %q", list.GetKind())
		}

		if len(list.Items) != 1 || list.Items[0].GetName() != "demo" {
			t.Fatalf("unexpected custom resource items: %#v", list.Items)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for custom resource result")
	}

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("expected channel to be closed after custom resource result")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel to close")
	}
}
