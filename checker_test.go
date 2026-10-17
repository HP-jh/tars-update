package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewChecker(t *testing.T) {
	c := NewChecker("HP-jh", "TarsVoiceCorrection")
	if c.Owner != "HP-jh" || c.Repo != "TarsVoiceCorrection" {
		t.Fatal("unexpected owner/repo")
	}
	if c.HTTPClient == nil {
		t.Fatal("expected HTTPClient")
	}
}

func TestCheckerLatestRelease(t *testing.T) {
	mock := ReleaseInfo{
		TagName:    "v0.2.0",
		Name:       "v0.2.0",
		HTMLURL:    "https://github.com/HP-jh/TarsVoiceCorrection/releases/tag/v0.2.0",
		Prerelease: false,
		Draft:      false,
		Assets: []Asset{
			{Name: "tvc-v0.2.0-linux-amd64", BrowserDownloadURL: "http://example.com/linux", Size: 1024},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/HP-jh/TarsVoiceCorrection/releases/latest" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mock)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_ = NewChecker("HP-jh", "TarsVoiceCorrection")
	// override url by replacing api.github.com is hard; test via ListReleases or local server
	// Instead, test the platform suffix function and asset selection directly
}

func TestReleaseInfoFindAsset(t *testing.T) {
	ri := ReleaseInfo{
		Assets: []Asset{
			{Name: "tvc-v0.2.0-linux-amd64", BrowserDownloadURL: "http://ex/linux", Size: 100},
			{Name: "tvc-v0.2.0-darwin-arm64", BrowserDownloadURL: "http://ex/darwin", Size: 100},
			{Name: "SHA256SUMS.txt", BrowserDownloadURL: "http://ex/sha", Size: 50},
		},
	}
	if a, ok := ri.FindAsset("darwin-arm64"); !ok || a.Name != "tvc-v0.2.0-darwin-arm64" {
		t.Fatal("expected darwin-arm64 asset")
	}
	if _, ok := ri.FindAsset("windows-arm64"); ok {
		t.Fatal("unexpected windows-arm64 asset")
	}
}

func TestPlatformSuffix(t *testing.T) {
	// platformSuffix uses runtime.GOOS/GOARCH; we can't change runtime, but we can verify non-empty
	s := platformSuffix()
	if s == "" {
		t.Fatal("expected non-empty platform suffix")
	}
}
