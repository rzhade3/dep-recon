package manifest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPullRustDependencyReadme(t *testing.T) {
	expectedReadme := "This is the readme for the package"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/crates/glob/0.3.1/readme" {
			t.Errorf("got %q, want %q", r.URL.Path, "/api/v1/crates/glob/0.3.1/readme")
		}
		w.Write([]byte(expectedReadme))
	}))
	rust := Rust{
		RegistryURL:        server.URL,
		DependencyFilePath: "../../fixtures/Cargo.toml",
	}
	readme, err := rust.PullDependencyReadme("glob", "0.3.1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if readme != expectedReadme {
		t.Errorf("Expected %s, got %s", expectedReadme, readme)
	}
}

func TestListRustDependencies(t *testing.T) {
	rust := Rust{
		RegistryURL:        DefaultRustRegistryURL,
		DependencyFilePath: "../../fixtures/Cargo.toml",
	}
	depList, err := rust.ListDependencies()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(depList.Dependencies) != 1 {
		t.Errorf("Expected 1, got %d", len(depList.Dependencies))
	}
	expectedDependencies := map[string]string{
		"glob": "0.3.1",
	}
	for key, value := range expectedDependencies {
		if depList.Dependencies[key] != value {
			t.Errorf("Expected %s, got %s", value, depList.Dependencies[key])
		}
	}
}
