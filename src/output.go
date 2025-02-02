package src

import (
	"encoding/json"
	"fmt"
	"strings"
)

var ValidOutputFormats = []string{"text", "json"}

func Output(matchedDependencies map[string][]string, outputFormat string) string {
	switch outputFormat {
	case "text":
		return TextOutput(matchedDependencies)
	case "json":
		return JsonOutput(matchedDependencies)
	default:
		return "Invalid output format"
	}
}

func TextOutput(matchedDependencies map[string][]string) string {
	if len(matchedDependencies) == 0 {
		return "No dependencies matched"
	}
	outputString := ""
	for dependency, matchedWords := range matchedDependencies {
		if len(matchedWords) == 0 {
			continue
		}
		outputString += fmt.Sprintf("%s matched for %v\n", dependency, matchedWords)
	}
	return strings.TrimSuffix(outputString, "\n")
}

func JsonOutput(matchedDependencies map[string][]string) string {
	// JSON format the matched dependencies
	jsonOutput, err := json.Marshal(matchedDependencies)
	if err != nil {
		return "Something went wrong"
	}
	return string(jsonOutput)
}
