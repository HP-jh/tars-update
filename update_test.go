package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStableVersions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stable_versions.json")

	// Save
	m := map[string]string{
		"TarsWingGuard": "v0.1.0",
		"_last_sync":    "2026-10-17T21:30:00+08:00",
	}
	if err := saveStableVersions(path, m); err != nil {
		t.Fatalf("save: %v", err)
	}

	// Load
	loaded := loadStableVersions(path)
	if loaded["TarsWingGuard"] != "v0.1.0" {
		t.Errorf("loaded version = %q, want v0.1.0", loaded["TarsWingGuard"])
	}
}

func TestPlatformAssetName(t *testing.T) {
	p := Platform{OS: "linux", Arch: "amd64"}
	name := p.AssetName("twg", "v0.2.0")
	if name != "twg-v0.2.0-linux-amd64" {
		t.Errorf("asset name = %q, want twg-v0.2.0-linux-amd64", name)
	}

	p2 := Platform{OS: "windows", Arch: "amd64"}
	name2 := p2.AssetName("twg", "v0.2.0")
	if name2 != "twg-v0.2.0-windows-amd64.exe" {
		t.Errorf("windows asset name = %q", name2)
	}
}

func TestParseSHA256SUMS(t *testing.T) {
	input := `808586072656585bbb34e4f2a230de46ae9112ccbc4a3285582921403bc7b6b4  knowledge-v0.1.0-android-arm64
00f1ce0fff2f4a6d2ef9cc8145c09e4263489b9427c1b281edd61a1e9eab3bf0  *knowledge-v0.1.0-darwin-amd64`
	m := parseSHA256SUMS(input)
	if m["knowledge-v0.1.0-android-arm64"] == "" {
		t.Error("missing android-arm64 hash")
	}
	if m["knowledge-v0.1.0-darwin-amd64"] == "" {
		t.Error("missing darwin-amd64 hash (star prefix)")
	}
}

func TestFileHelpers(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	data := []byte("hello update")
	if err := writeFile(src, data); err != nil {
		t.Fatalf("writeFile: %v", err)
	}

	if !fileExists(src) {
		t.Fatal("fileExists should be true")
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile: %v", err)
	}

	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(b) != string(data) {
		t.Errorf("dst content = %q, want %q", b, data)
	}

	hash, err := fileSHA256(src)
	if err != nil {
		t.Fatalf("fileSHA256: %v", err)
	}
	if len(hash) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash))
	}
}
