package manifest

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type Javascript struct {
	// The URL of the package manager's registry
	RegistryURL string
	// The filepath to the dependency file
	DependencyFilePath string
}

var DefaultJavascriptRegistryURL = "https://registry.npmjs.org"

type PackageJson struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type NpmResponse struct {
	Readme string `json:"readme"`
}

func (j Javascript) PullDependencyReadme(dependency, version string) (string, error) {
	npmUrl, err := url.JoinPath(j.RegistryURL, dependency)
	if err != nil {
		return "", err
	}
	res, err := http.Get(npmUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", err
	}
	var data NpmResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Readme, nil
}

func (j Javascript) GetEcosystem() string {
	return "npm"
}

func (j Javascript) readDependencyFile() ([]byte, error) {
	content, err := os.ReadFile(j.DependencyFilePath)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func (j Javascript) ListDependencies() (DependencyList, error) {
	// Parse the package.json
	var packageDependencies PackageJson
	var depList DependencyList
	content, err := j.readDependencyFile()
	if err != nil {
		return depList, err
	}
	err = json.Unmarshal(content, &packageDependencies)
	if err != nil {
		return depList, err
	}
	depList.Dependencies = packageDependencies.Dependencies
	depList.DevDependencies = packageDependencies.DevDependencies
	return depList, nil
}

func (j Javascript) GetDependencyFilePath() string {
	return j.DependencyFilePath
}
