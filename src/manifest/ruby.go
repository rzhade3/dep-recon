package manifest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Ruby struct {
	// The URL of the package manager's registry
	RegistryURL string
	// The filepath to the dependency file
	DependencyFilePath string
}

var DefaultRubyRegistryURL = "https://rubygems.org"

type RubyGemsResponse struct {
	Info string `json:"info"`
}

func (r Ruby) PullDependencyReadme(dependency, version string) (string, error) {
	gemfilePath, err := url.JoinPath(r.RegistryURL, "api/v1/gems", dependency+".json")
	if err != nil {
		return "", err
	}
	res, err := http.Get(gemfilePath)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to pull readme for %s from RubyGems: %s", dependency, res.Status)
	}
	// Read the response body as a JSON
	var data RubyGemsResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	// Extract the readme from the JSON
	readme := data.Info
	// // Perhaps there is a homepage with more information. If so, let's try and parse that
	// homepage, ok := data["homepage_uri"].(string)
	// if ok {
	// 	res, err := http.Get(homepage)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	defer res.Body.Close()
	// 	if res.StatusCode != http.StatusOK {
	// 		return "", fmt.Errorf("failed to pull homepage for %s from RubyGems: %s", dependency, res.Status)
	// 	}
	// 	// Return the response body
	// 	readme = strings.NewDecoder(res.Body).Decode(&data)
	// }
	return readme, nil
}

func (r Ruby) GetEcosystem() string {
	return "rubygems"
}

func (r Ruby) readDependencyFile() (string, error) {
	content, err := os.ReadFile(r.DependencyFilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (r Ruby) ListDependencies() (DependencyList, error) {
	gem_regex := regexp.MustCompile(`gem ['"](?P<name>[^'"]+)['"](, ['"](?P<version>[^'"]+)['"])?`)
	content, err := r.readDependencyFile()
	if err != nil {
		return DependencyList{}, err
	}
	matches := gem_regex.FindAllStringSubmatch(content, -1)
	depList := DependencyList{
		DevDependencies: map[string]string{},
		Dependencies:    map[string]string{},
	}
	for _, match := range matches {
		packageName := match[1]
		packageVersion := strings.TrimSpace(match[3])
		packageVersion = strings.Trim(packageVersion, "'")
		packageVersion = strings.Trim(packageVersion, "\"")
		depList.Dependencies[packageName] = packageVersion
	}
	return depList, nil
}

func (r Ruby) GetDependencyFilePath() string {
	return r.DependencyFilePath
}
