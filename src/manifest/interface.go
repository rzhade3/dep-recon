package manifest

type DependencyList struct {
	DevDependencies map[string]string `json:"devDependencies"`
	Dependencies    map[string]string `json:"dependencies"`
}

type Language interface {
	// Pulls a README from the package manager's registry
	PullDependencyReadme(dependency, version string) (string, error)
	// Attempt to read from the cache
	GetEcosystem() string
	// ListDependencies lists the dependencies for a given language
	ListDependencies() (DependencyList, error)
	// GetDependencyFilePath returns the filepath to the dependency file
	GetDependencyFilePath() string
}
