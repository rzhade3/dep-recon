package manifest

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
)

type Rust struct {
	// The URL of the package manager's registry
	RegistryURL string
	// The filepath to the dependency file
	DependencyFilePath string
}

var DefaultRustRegistryURL = "https://crates.io"

func (r Rust) PullDependencyReadme(dependency, version string) (string, error) {
	// The api_path looks something like "/api/v1/crates/{crate_name}/{version}/readme"
	crates_url, err := url.JoinPath(r.RegistryURL, "api/v1/crates", dependency, version, "readme")
	if err != nil {
		return "", err
	}
	// Fetch the README from the crates_url
	res, err := http.Get(crates_url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", err
	}
	// Return the response body as a string
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (r Rust) GetEcosystem() string {
	return "cargo"
}

type TomlData struct {
	Dependencies []string
}

func (r Rust) readDependencyFile() (string, error) {
	content, err := os.ReadFile(r.DependencyFilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (r Rust) ListDependencies() (DependencyList, error) {
	var depList DependencyList
	content, err := r.readDependencyFile()
	if err != nil {
		return depList, err
	}
	// Parse the TOML file
	_, err = toml.Decode(content, &depList)
	if err != nil {
		return depList, err
	}
	return depList, nil
}

func (r Rust) GetDependencyFilePath() string {
	return r.DependencyFilePath
}
