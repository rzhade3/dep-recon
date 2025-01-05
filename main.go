package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rzhade3/dep-recon/src"
	"github.com/rzhade3/dep-recon/src/db"
)

func main() {
	// Parse the flags
	var cacheDir string
	flag.StringVar(&cacheDir, "cache", ".", "The directory to store the cache")
	var manifestFile string
	flag.StringVar(&manifestFile, "scan", "", "The manifest file to scan for dependencies (required)")
	var keywordFile string
	flag.StringVar(&keywordFile, "keywords", "keywords.json", "The file containing keywords to search for in the README")
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "Print help")
	flag.Parse()

	if helpFlag {
		flag.PrintDefaults()
		return
	}
	if manifestFile == "" {
		fmt.Println("Please provide a manifest file to scan")
		return
	}
	// Check if manifest file exists on filesystem
	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		fmt.Printf("Manifest file %s does not exist\n", manifestFile)
		return
	}
	manifest, err := src.ValidateManifestFilepath(manifestFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	keywords, err := src.LoadKeywords(keywordFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialize the cache
	cacheFilename := filepath.Join(cacheDir, "cache.db")
	dbConfig, err := db.InitializeCache(cacheFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	depList, err := manifest.ListDependencies()
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: Make readme dependent on version number
	for dependency, version := range depList.Dependencies {
		readme, err := dbConfig.FetchDependencyReadme(manifest, dependency, version)
		if err != nil {
			fmt.Println(err)
		}
		matched_words, err := src.ReadmeRecon(readme, keywords)
		if err != nil {
			fmt.Println(err)
		}
		if len(matched_words) > 0 {
			fmt.Printf("Dependency: %s matched with the following keywords %v\n", dependency, matched_words)
		}
	}
}
