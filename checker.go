package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// ReleaseInfo holds parsed GitHub release data.
type ReleaseInfo struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	HTMLURL     string  `json:"html_url"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
	Prerelease  bool    `json:"prerelease"`
	Draft       bool    `json:"draft"`
}

// Asset represents a release asset.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
}

// Checker queries GitHub for available updates.
type Checker struct {
	Owner      string
	Repo       string
	HTTPClient *http.Client
}

// NewChecker creates a checker for the given GitHub repo.
func NewChecker(owner, repo string) *Checker {
	return &Checker{
		Owner: owner,
		Repo:  repo,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LatestRelease fetches the latest release from GitHub.
func (c *Checker) LatestRelease() (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", c.Owner, c.Repo)
	return c.fetchRelease(url)
}

// ReleaseByTag fetches a specific release by tag.
func (c *Checker) ReleaseByTag(tag string) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", c.Owner, c.Repo, tag)
	return c.fetchRelease(url)
}

// ListReleases fetches the last N releases.
func (c *Checker) ListReleases(perPage int) ([]ReleaseInfo, error) {
	if perPage <= 0 {
		perPage = 10
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d", c.Owner, c.Repo, perPage)
	resp, err := c.httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status %d", resp.StatusCode)
	}
	var releases []ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decode releases: %w", err)
	}
	return releases, nil
}

// FindAssetForPlatform selects the asset matching current OS/Arch.
func (r *ReleaseInfo) FindAssetForPlatform() (Asset, bool) {
	osArch := platformSuffix()
	for _, a := range r.Assets {
		if strings.Contains(a.Name, osArch) {
			return a, true
		}
	}
	return Asset{}, false
}

// FindAsset selects an asset by substring match.
func (r *ReleaseInfo) FindAsset(substr string) (Asset, bool) {
	for _, a := range r.Assets {
		if strings.Contains(a.Name, substr) {
			return a, true
		}
	}
	return Asset{}, false
}

func (c *Checker) fetchRelease(url string) (*ReleaseInfo, error) {
	resp, err := c.httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("release not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status %d", resp.StatusCode)
	}
	var ri ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&ri); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}
	return &ri, nil
}

func (c *Checker) httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "tars-update/1.0")
	return c.HTTPClient.Do(req)
}

func platformSuffix() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goos == "windows" {
		return goos + "-" + goarch + ".exe"
	}
	return goos + "-" + goarch
}
