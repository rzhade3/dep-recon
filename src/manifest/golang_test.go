package manifest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPullGolangDependencyReadme(t *testing.T) {
	htmlResponse, err := os.ReadFile("../../fixtures/golang_readme.html")
	expectedReadme := "Package gabs implements a wrapper around creating and parsing unknown or dynamic map structures resulting from JSON parsing."
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/github.com/Jeffail/gabs/v2" {
			t.Errorf("got %q, want %q", r.URL.Path, "/github.com/Jeffail/gabs/v2")
		}
		w.Write(htmlResponse)
	}))
	golang := Golang{
		RegistryURL:        server.URL,
		DependencyFilePath: "../../fixtures/go.mod",
	}
	readme, err := golang.PullDependencyReadme("github.com/Jeffail/gabs/v2", "v2.7.0")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if readme != expectedReadme {
		t.Errorf("Expected %s, got %s", expectedReadme, readme)
	}
}

func TestParseGoMod(t *testing.T) {
	golang := Golang{
		RegistryURL:        "https://pkg.go.dev",
		DependencyFilePath: "../../fixtures/go.mod",
	}
	depList, err := golang.ListDependencies()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if len(depList.Dependencies) != 2 {
		t.Errorf("Expected 2, got %d", len(depList.Dependencies))
	}
	expectedDependencies := map[string]string{
		"github.com/Jeffail/gabs/v2": "v2.7.0",
		"modernc.org/sqlite":         "v1.33.1",
	}
	for key, value := range expectedDependencies {
		if depList.Dependencies[key] != value {
			t.Errorf("Expected %s, got %s", value, depList.Dependencies[key])
		}
	}
}
