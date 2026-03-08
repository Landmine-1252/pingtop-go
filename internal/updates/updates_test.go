package updates

import (
	"testing"
	"time"

	"pingtop/internal/pingtop"
)

func TestVersionHelpers(t *testing.T) {
	if _, _, _, ok := parseVersionTag("v" + pingtop.Version); !ok {
		t.Fatalf("expected version tag to parse: %s", pingtop.Version)
	}
	if normalizeRepoURL("https://github.com/Landmine-1252/pingtop.git") != "https://github.com/Landmine-1252/pingtop" {
		t.Fatal("failed to normalize https repo url")
	}
	if normalizeRepoURL("git@github.com:Landmine-1252/pingtop.git") != "https://github.com/Landmine-1252/pingtop" {
		t.Fatal("failed to normalize ssh repo url")
	}
	if !isNewerVersion("v0.1.0", "v0.2.0") {
		t.Fatal("expected newer version comparison to succeed")
	}
	if isNewerVersion("v0.2.0", "v0.1.0") {
		t.Fatal("expected older version comparison to fail")
	}
}

func TestBuildReleaseAPIURLRequiresGitHubRepo(t *testing.T) {
	url, err := buildReleaseAPIURL("https://github.com/Landmine-1252/pingtop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "https://api.github.com/repos/Landmine-1252/pingtop/releases/latest"
	if url != expected {
		t.Fatalf("unexpected api url: %q", url)
	}
	if _, err := buildReleaseAPIURL("https://example.com/notgithub/repo"); err == nil {
		t.Fatal("expected non-github repo to fail")
	}
}

func TestUpdateManagerMarksAvailableVersion(t *testing.T) {
	manager := NewUpdateManager(
		"v0.1.0",
		"https://github.com/Landmine-1252/pingtop",
		true,
		func(repoURL string, timeout time.Duration) (string, string, error) {
			return "v0.2.0", "https://github.com/Landmine-1252/pingtop/releases/tag/v0.2.0", nil
		},
	)
	manager.run()
	status := manager.Snapshot()
	if status.State != "available" || status.LatestVersion != "v0.2.0" {
		t.Fatalf("unexpected update status: %#v", status)
	}
}
