package util

import (
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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

	go Getter("default", clientset, []ResourceMeta{
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

	go Getter("default", clientset, []ResourceMeta{
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

	go Getter("default", clientset, []ResourceMeta{
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
