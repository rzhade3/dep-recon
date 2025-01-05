package manifest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPullRubygemsReadme(t *testing.T) {
	expected_readme := "This is the readme for the package"
	// Interpolate into a JSON body for the mock response
	mockResponse := fmt.Sprintf(`{"info": "%s"}`, expected_readme)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request
		if r.URL.Path != "/api/v1/gems/dependency.json" {
			t.Errorf("got %q, want %q", r.URL.Path, "/api/v1/gems/dependency.json")
		}
		// Send the response
		w.Write([]byte(mockResponse))
	}))
	// Make a request to the mock server
	ruby := Ruby{
		RegistryURL:        server.URL,
		DependencyFilePath: "Gemfile",
	}
	readme, err := ruby.PullDependencyReadme("dependency", "1.0.0")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if readme != expected_readme {
		t.Errorf("Expected %s, got %s", expected_readme, readme)
	}
}

func TestPullRubygemsReadmeNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send the response
		w.WriteHeader(http.StatusNotFound)
	}))
	// Make a request to the mock server
	ruby := Ruby{
		RegistryURL:        server.URL,
		DependencyFilePath: "Gemfile",
	}
	readme, err := ruby.PullDependencyReadme("dependency", "1.0.0")
	if err == nil {
		t.Error("expected an error, got nil")
	}
	if readme != "" {
		t.Errorf("expected empty string, got %q", readme)
	}
}

func TestParseGemfile(t *testing.T) {
	ruby := Ruby{
		RegistryURL:        "https://rubygems.org",
		DependencyFilePath: "../../fixtures/Gemfile",
	}
	depList, err := ruby.ListDependencies()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(depList.Dependencies) != 4 {
		t.Errorf("Expected 4, got %d", len(depList.Dependencies))
	}
	expected := map[string]string{
		"rails":       ">= 5.2.0",
		"sqlite3":     "",
		"puma":        "~> 3.11",
		"web-console": ">= 3.3.0",
	}
	for key, value := range expected {
		if depList.Dependencies[key] != value {
			t.Errorf("Expected %s, got %s", value, depList.Dependencies[key])
		}
	}
}
