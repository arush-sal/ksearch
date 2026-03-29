package cmd

import (
	"testing"
	"time"

	"github.com/arush-sal/ksearch/pkg/cache"
	"github.com/arush-sal/ksearch/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func TestRunUsesCacheBeforeDiscovery(t *testing.T) {
	t.Cleanup(func() {
		currentContextNameFn = currentContextName
		readCacheFn = cache.Read
		writeCachedSectionsFn = writeCachedSections
		writeCacheFn = cache.Write
		getConfigOrDieFn = func() *rest.Config { return config.GetConfigOrDie() }
		newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface { return kubernetes.NewForConfigOrDie(cfg) }
		discoverResourcesFn = util.Discover
		getterFn = util.Getter
	})

	currentContextNameFn = func() (string, error) {
		return "ctx", nil
	}
	readCacheFn = func(key string, ttl time.Duration) (*cache.CacheEntry, error) {
		return &cache.CacheEntry{
			Sections: []cache.SectionEntry{{Kind: "Pods", Output: ""}},
		}, nil
	}
	writeCachedSectionsFn = func(sections []cache.SectionEntry) error {
		return nil
	}
	getConfigOrDieFn = func() *rest.Config {
		t.Fatal("expected cache hit to return before config initialization")
		return nil
	}
	newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface {
		t.Fatal("expected cache hit to return before clientset creation")
		return nil
	}
	discoverResourcesFn = func(discoveryClient discovery.DiscoveryInterface, kinds string) ([]util.ResourceMeta, error) {
		t.Fatal("expected cache hit to return before discovery")
		return nil, nil
	}

	if err := runRoot(rootCmd, nil); err != nil {
		t.Fatalf("runRoot returned error: %v", err)
	}
}

func TestRunUsesFetchedKindsForCacheSections(t *testing.T) {
	t.Cleanup(func() {
		currentContextNameFn = currentContextName
		readCacheFn = cache.Read
		writeCachedSectionsFn = writeCachedSections
		writeCacheFn = cache.Write
		getConfigOrDieFn = func() *rest.Config { return config.GetConfigOrDie() }
		newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface { return kubernetes.NewForConfigOrDie(cfg) }
		discoverResourcesFn = util.Discover
		getterFn = util.Getter
	})

	currentContextNameFn = func() (string, error) {
		return "ctx", nil
	}
	readCacheFn = func(key string, ttl time.Duration) (*cache.CacheEntry, error) {
		return nil, nil
	}
	getConfigOrDieFn = func() *rest.Config {
		return &rest.Config{}
	}
	newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface {
		return fake.NewSimpleClientset()
	}
	discoverResourcesFn = func(discoveryClient discovery.DiscoveryInterface, kinds string) ([]util.ResourceMeta, error) {
		return []util.ResourceMeta{
			{Kind: "ConfigMap", Resource: "configmaps", Namespaced: true},
			{Kind: "Pod", Resource: "pods", Namespaced: true},
		}, nil
	}
	getterFn = func(namespace string, clientset kubernetes.Interface, cfg *rest.Config, resources []util.ResourceMeta, ch chan util.FetchResult) {
		defer close(ch)
		ch <- util.FetchResult{Kind: "ConfigMap", Resource: nil}
		ch <- util.FetchResult{Kind: "Pod", Resource: &v1.PodList{}}
	}

	var captured []cache.SectionEntry
	writeCacheFn = func(key string, meta cache.CacheMeta, sections []cache.SectionEntry) error {
		captured = append([]cache.SectionEntry(nil), sections...)
		return nil
	}

	if err := runRoot(rootCmd, nil); err != nil {
		t.Fatalf("runRoot returned error: %v", err)
	}

	if len(captured) != 2 {
		t.Fatalf("expected 2 cache sections, got %d", len(captured))
	}

	if captured[0].Kind != "ConfigMap" {
		t.Fatalf("expected first section kind ConfigMap, got %q", captured[0].Kind)
	}

	if captured[1].Kind != "Pod" {
		t.Fatalf("expected second section kind Pod, got %q", captured[1].Kind)
	}
}
