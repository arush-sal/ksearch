package cache

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SectionEntry struct {
	Kind   string `json:"kind"`
	Output string `json:"output"`
}

type CacheMeta struct {
	Context    string
	Namespace  string
	Kinds      string
	Pattern    string
	TTLSeconds int
}

type CacheEntry struct {
	WrittenAt  time.Time      `json:"written_at"`
	TTLSeconds int            `json:"ttl_seconds"`
	Context    string         `json:"context"`
	Namespace  string         `json:"namespace"`
	Kinds      string         `json:"kinds"`
	Pattern    string         `json:"pattern"`
	Sections   []SectionEntry `json:"sections"`
}

func KeyFor(context, namespace, kinds, pattern string) string {
	parts := []string{
		context,
		namespace,
		normalizedKinds(kinds),
		pattern,
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:])
}

func Read(key string, ttl time.Duration) (*CacheEntry, error) {
	data, err := os.ReadFile(cacheFile(key))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}

	if entry.WrittenAt.Add(ttl).Before(time.Now()) {
		return nil, nil
	}

	return &entry, nil
}

func Write(key string, meta CacheMeta, sections []SectionEntry) error {
	if err := os.MkdirAll(cacheDir(), 0o700); err != nil {
		return err
	}

	entry := CacheEntry{
		WrittenAt:  time.Now().UTC(),
		TTLSeconds: meta.TTLSeconds,
		Context:    meta.Context,
		Namespace:  meta.Namespace,
		Kinds:      meta.Kinds,
		Pattern:    meta.Pattern,
		Sections:   encodeSections(sections),
	}

	tempFile, err := os.CreateTemp(cacheDir(), "*.tmp")
	if err != nil {
		return err
	}

	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if err := tempFile.Chmod(0o600); err != nil {
		tempFile.Close()
		return err
	}

	if err := json.NewEncoder(tempFile).Encode(entry); err != nil {
		tempFile.Close()
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	finalPath := cacheFile(key)
	if err := os.Rename(tempName, finalPath); err != nil {
		return err
	}

	return os.Chmod(finalPath, 0o600)
}

func cacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".kube", "ksearch", "cache")
	}

	return filepath.Join(homeDir, ".kube", "ksearch", "cache")
}

func cacheFile(key string) string {
	return filepath.Join(cacheDir(), key+".json")
}

func normalizedKinds(kinds string) string {
	if strings.TrimSpace(kinds) == "" {
		return ""
	}

	parts := strings.Split(kinds, ",")
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		if part == "" {
			continue
		}
		normalized = append(normalized, part)
	}

	sort.Strings(normalized)
	return strings.Join(normalized, ",")
}

func encodeSections(sections []SectionEntry) []SectionEntry {
	encoded := make([]SectionEntry, 0, len(sections))
	for _, section := range sections {
		encoded = append(encoded, SectionEntry{
			Kind:   section.Kind,
			Output: base64.StdEncoding.EncodeToString([]byte(section.Output)),
		})
	}

	return encoded
}
