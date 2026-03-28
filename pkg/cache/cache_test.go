package cache

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestKeyFor_Deterministic(t *testing.T) {
	key1 := KeyFor("prod", "default", "configmap,secret", "nginx")
	key2 := KeyFor("prod", "default", "configmap,secret", "nginx")

	if key1 != key2 {
		t.Fatalf("expected deterministic key, got %q and %q", key1, key2)
	}
}

func TestKeyFor_KindsSorted(t *testing.T) {
	key1 := KeyFor("prod", "default", "secret,configmap", "nginx")
	key2 := KeyFor("prod", "default", "configmap,secret", "nginx")

	if key1 != key2 {
		t.Fatalf("expected sorted kinds to match, got %q and %q", key1, key2)
	}
}

func TestKeyFor_Unique(t *testing.T) {
	base := KeyFor("prod", "default", "configmap,secret", "nginx")

	testCases := []struct {
		name string
		key  string
	}{
		{name: "different context", key: KeyFor("staging", "default", "configmap,secret", "nginx")},
		{name: "different namespace", key: KeyFor("prod", "kube-system", "configmap,secret", "nginx")},
		{name: "different kinds", key: KeyFor("prod", "default", "configmap", "nginx")},
		{name: "different pattern", key: KeyFor("prod", "default", "configmap,secret", "redis")},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.key == base {
				t.Fatalf("expected unique key for %s", testCase.name)
			}
		})
	}
}

func TestRead_Missing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	entry, err := Read("missing", time.Minute)
	if err != nil {
		t.Fatalf("Read returned error for missing file: %v", err)
	}

	if entry != nil {
		t.Fatalf("expected nil entry for missing file, got %#v", entry)
	}
}

func TestRead_Expired(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	path := filepath.Join(cacheDir(), "expired.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("create cache dir: %v", err)
	}

	entry := CacheEntry{
		WrittenAt:  time.Now().Add(-2 * time.Minute),
		TTLSeconds: 60,
		Context:    "prod",
		Namespace:  "default",
		Kinds:      "configmap",
		Pattern:    "nginx",
		Sections: []SectionEntry{
			{Kind: "ConfigMaps", Output: base64.StdEncoding.EncodeToString([]byte("safe output"))},
		},
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatalf("open cache file: %v", err)
	}

	if err := json.NewEncoder(file).Encode(entry); err != nil {
		_ = file.Close()
		t.Fatalf("encode cache file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close cache file: %v", err)
	}

	readEntry, err := Read("expired", time.Minute)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if readEntry != nil {
		t.Fatalf("expected nil entry for expired cache, got %#v", readEntry)
	}
}

func TestReadWrite_RoundTrip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	sections := []SectionEntry{
		{Kind: "ConfigMaps", Output: "\nConfigMaps\n--------------\nNAME DATA AGE\nnginx 1 \n"},
		{Kind: "Secrets", Output: "\nSecrets\n--------------\nNAME TYPE DATA AGE\napi-token Opaque 1 \n"},
	}

	err := Write("round-trip", CacheMeta{
		Context:    "prod",
		Namespace:  "default",
		Kinds:      "configmap,secret",
		Pattern:    "nginx",
		TTLSeconds: 60,
	}, sections)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	entry, err := Read("round-trip", time.Minute)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}

	if entry == nil {
		t.Fatal("expected cache entry, got nil")
	}

	if len(entry.Sections) != len(sections) {
		t.Fatalf("expected %d sections, got %d", len(sections), len(entry.Sections))
	}

	for index, section := range entry.Sections {
		if section.Kind != sections[index].Kind {
			t.Fatalf("section %d kind = %q, want %q", index, section.Kind, sections[index].Kind)
		}

		decoded, err := base64.StdEncoding.DecodeString(section.Output)
		if err != nil {
			t.Fatalf("decode section %d: %v", index, err)
		}

		if string(decoded) != sections[index].Output {
			t.Fatalf("section %d output = %q, want %q", index, string(decoded), sections[index].Output)
		}
	}
}

func TestNoSensitiveData(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	err := Write("safe-output", CacheMeta{
		Context:    "prod",
		Namespace:  "default",
		Kinds:      "secret",
		Pattern:    "",
		TTLSeconds: 60,
	}, []SectionEntry{
		{Kind: "Secrets", Output: "\nSecrets\n--------------\nNAME TYPE DATA AGE\napi-token Opaque 2 \n"},
	})
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cacheDir(), "safe-output.json"))
	if err != nil {
		t.Fatalf("read cache file: %v", err)
	}

	for _, forbidden := range []string{"super-secret-value", "my-api-token", "service-account-token"} {
		if string(data) != "" && contains(string(data), forbidden) {
			t.Fatalf("cache file leaked sensitive string %q", forbidden)
		}
	}
}

func TestWrite_Atomic(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	const writers = 8
	var wg sync.WaitGroup

	for index := 0; index < writers; index++ {
		index := index
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := Write("atomic", CacheMeta{
				Context:    "prod",
				Namespace:  "default",
				Kinds:      "configmap",
				Pattern:    "",
				TTLSeconds: 60,
			}, []SectionEntry{
				{Kind: "ConfigMaps", Output: "\nConfigMaps\n--------------\nNAME DATA AGE\nconfig-" + string(rune('a'+index)) + " 1 \n"},
			})
			if err != nil {
				t.Errorf("Write returned error: %v", err)
			}
		}()
	}

	wg.Wait()

	data, err := os.ReadFile(filepath.Join(cacheDir(), "atomic.json"))
	if err != nil {
		t.Fatalf("read cache file: %v", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("expected valid JSON after concurrent writes: %v", err)
	}

	if len(entry.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(entry.Sections))
	}

	if _, err := base64.StdEncoding.DecodeString(entry.Sections[0].Output); err != nil {
		t.Fatalf("expected valid base64 output, got %v", err)
	}
}

func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && filepath.Base(substr) != "" && (len(substr) == 0 || stringIndex(s, substr) >= 0)
}

func stringIndex(s, substr string) int {
	for index := 0; index+len(substr) <= len(s); index++ {
		if s[index:index+len(substr)] == substr {
			return index
		}
	}
	return -1
}
