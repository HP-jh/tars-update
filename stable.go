package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// StableVersions is the on-disk fallback table.
type StableVersions struct {
	mu     sync.RWMutex
	path   string
	data   map[string]string
}

// loadStableVersions reads the stable versions file.
func loadStableVersions(path string) map[string]string {
	m := make(map[string]string)
	b, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	_ = json.Unmarshal(b, &m)
	return m
}

// saveStableVersions writes the stable versions file atomically.
func saveStableVersions(path string, m map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
