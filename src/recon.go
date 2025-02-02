package src

import (
	"fmt"
	"path/filepath"

	"github.com/rzhade3/dep-recon/src/manifest"
)

var packageManifests = map[string][]string{
	"javascript": {"package.json"},
	"ruby":       {"Gemfile"},
	"golang":     {"go.mod"},
	"rust":       {"Cargo.toml"},
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
