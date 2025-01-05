package src

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rzhade3/dep-recon/src/manifest"
)

var packageManifests = map[string][]string{
	"javascript": {"package.json"},
	"ruby":       {"Gemfile"},
	"golang":     {"go.mod"},
	"rust":       {"Cargo.toml"},
}

// FindPackageManifests searches for package manifests in the given directory and returns a list of manifest objects
// in the directory and its subdirectories
// Returns an error if the directory cannot be read
func FindPackageManifests(scanDir string) ([]manifest.Language, error) {
	// Golang doesn't support the ** glob pattern, so we need to write a custom walk function
	languages := []manifest.Language{}
	err := filepath.Walk(scanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for lang, patterns := range packageManifests {
			for _, pattern := range patterns {
				match, err := filepath.Match(pattern, info.Name())
				if err != nil {
					return err
				}
				if match {
					switch lang {
					case "golang":
						languages = append(languages, manifest.Golang{
							RegistryURL:        manifest.DefaultGolangRegistryURL,
							DependencyFilePath: path,
						})
					case "ruby":
						languages = append(languages, manifest.Ruby{
							RegistryURL:        manifest.DefaultRubyRegistryURL,
							DependencyFilePath: path,
						})
					case "javascript":
						languages = append(languages, manifest.Javascript{
							RegistryURL:        manifest.DefaultJavascriptRegistryURL,
							DependencyFilePath: path,
						})
					case "rust":
						languages = append(languages, manifest.Rust{
							RegistryURL:        manifest.DefaultRustRegistryURL,
							DependencyFilePath: path,
						})
					}
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return languages, nil
}

// ValidateManifestFilepath checks if the manifest file matches a known package manager
// Returns the manifest object if it matches, otherwise returns an error
func ValidateManifestFilepath(manifestFilepath string) (manifest.Language, error) {
	// See if the manifest matches any of the known package managers
	for lang, patterns := range packageManifests {
		for _, pattern := range patterns {
			match, err := filepath.Match(pattern, filepath.Base(manifestFilepath))
			if err != nil {
				return nil, err
			}
			if match {
				switch lang {
				case "golang":
					return manifest.Golang{
						RegistryURL:        manifest.DefaultGolangRegistryURL,
						DependencyFilePath: manifestFilepath,
					}, nil
				case "ruby":
					return manifest.Ruby{
						RegistryURL:        manifest.DefaultRubyRegistryURL,
						DependencyFilePath: manifestFilepath,
					}, nil
				case "javascript":
					return manifest.Javascript{
						RegistryURL:        manifest.DefaultJavascriptRegistryURL,
						DependencyFilePath: manifestFilepath,
					}, nil
				case "rust":
					return manifest.Rust{
						RegistryURL:        manifest.DefaultRustRegistryURL,
						DependencyFilePath: manifestFilepath,
					}, nil
				}
				return nil, fmt.Errorf("manifest file %s does not match any known package manager", manifestFilepath)
			}
		}
	}
	return nil, fmt.Errorf("manifest file %s does not match any known package manager", manifestFilepath)
}

// Parses a content and sees if any of the keywords match in it
func ReadmeRecon(content string, keywords Keywords) ([]string, error) {
	matched_concepts := []string{}
	for concept, search_words := range keywords {
		for _, search_word := range search_words {
			// Check if the content contains the keyword
			if strings.Contains(content, search_word) {
				matched_concepts = append(matched_concepts, concept)
				break
			}
		}
	}
	return matched_concepts, nil
}

type Keywords map[string][]string

func LoadKeywords(keywordFile string) (Keywords, error) {
	// Read the file
	// Parse the JSON
	// Return the keywords
	file, err := os.ReadFile(keywordFile)
	if err != nil {
		return nil, err
	}
	var keywords Keywords
	err = json.Unmarshal(file, &keywords)
	if err != nil {
		return nil, err
	}
	return keywords, nil
}
