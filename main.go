package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"

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
	var aiRefinementFlag bool
	flag.BoolVar(&aiRefinementFlag, "ai", false, "Use AI to refine the search")
	var aiExamplesFile string
	flag.StringVar(&aiExamplesFile, "ai-examples", "examples.json", "The file containing examples to use for AI refinement")
	var helpFlag bool
	flag.BoolVar(&helpFlag, "help", false, "Print help")
	var formatType string
	flag.StringVar(&formatType, "format", "text", "The format to output the results in (text or json)")
	flag.Parse()

	if helpFlag {
		flag.PrintDefaults()
		return
	}
	if manifestFile == "" {
		fmt.Println("Please provide a manifest file to scan")
		return
	}
	if !slices.Contains(src.ValidOutputFormats, formatType) {
		fmt.Printf("Invalid output format, must be %v\n", src.ValidOutputFormats)
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
	var examples src.Examples
	if aiRefinementFlag {
		examples, err = src.LoadExamples("examples.json")
		if err != nil {
			fmt.Println(err)
			return
		}
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
	matchedDependencies := make(map[string][]string)
	for dependency, version := range depList.Dependencies {
		readme, err := dbConfig.FetchDependencyReadme(manifest, dependency, version)
		if err != nil {
			fmt.Println(err)
		}
		matchedWords, err := src.ReadmeMatch(readme, keywords, examples, aiRefinementFlag)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(matchedWords) != 0 {
			matchedDependencies[dependency] = matchedWords
		}
	}
	output := src.Output(matchedDependencies, formatType)
	fmt.Println(output)
}
