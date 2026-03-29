package cmd

import (
	"testing"
	"time"

	"github.com/arush-sal/ksearch/pkg/cache"
	"github.com/arush-sal/ksearch/pkg/util"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func TestRunUsesCacheBeforeDiscovery(t *testing.T) {
	t.Cleanup(func() {
		currentContextNameFn = currentContextName
		readCacheFn = cache.Read
		writeCachedSectionsFn = writeCachedSections
		getConfigOrDieFn = func() *rest.Config { return config.GetConfigOrDie() }
		newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface { return kubernetes.NewForConfigOrDie(cfg) }
		discoverResourcesFn = util.Discover
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
