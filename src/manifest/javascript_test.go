package manifest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPullNpmReadme(t *testing.T) {
	mockResponse := `{"readme": "This is the readme for the package"}`
	expectedReadme := "This is the readme for the package"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request
		if r.URL.Path != "/express" {
			t.Errorf("got %q, want %q", r.URL.Path, "/express")
		}
		// Send the response
		w.Write([]byte(mockResponse))
	}))
	javascript := Javascript{
		RegistryURL:        server.URL,
		DependencyFilePath: "../../fixtures/package.json",
	}
	readme, err := javascript.PullDependencyReadme("express", "4.18.2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if readme != expectedReadme {
		t.Errorf("Expected %s, got %s", expectedReadme, readme)
	}
}

func TestParsePackageJson(t *testing.T) {
	javascript := Javascript{
		RegistryURL:        DefaultJavascriptRegistryURL,
		DependencyFilePath: "../../fixtures/package.json",
	}
	depList, err := javascript.ListDependencies()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if len(depList.Dependencies) != 1 {
		t.Errorf("Expected 1, got %d", len(depList.Dependencies))
	}
	expectedDependencies := map[string]string{
		"express": "^4.18.2",
	}
	for key, value := range expectedDependencies {
		if depList.Dependencies[key] != value {
			t.Errorf("Expected %s, got %s", value, depList.Dependencies[key])
		}
	}
	if len(depList.DevDependencies) != 2 {
		t.Errorf("Expected 1, got %d", len(depList.DevDependencies))
	}
	expectedDevDependencies := map[string]string{
		"eslint":  "^8.47.0",
		"nodemon": "^3.0.1",
	}
	for key, value := range expectedDevDependencies {
		if depList.DevDependencies[key] != value {
			t.Errorf("Expected %s, got %s", value, depList.DevDependencies[key])
		}
	}
}
