package src

import (
	"sort"
	"testing"
)

func TestReadmeKeywordMatch(t *testing.T) {
	concepts := Concepts{
		"javascript": []string{"npm", "node", "javascript"},
		"ruby":       []string{"gem", "ruby"},
		"golang":     []string{"go", "golang"},
	}
	matched, err := ReadmeKeywordMatch("This contains npm and gem words", concepts)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(matched) != 2 {
		t.Errorf("Expected 2, got %d", len(matched))
	}
	sort.Strings(matched)
	if matched[0] != "javascript" || matched[1] != "ruby" {
		t.Errorf("Expected [javascript ruby], got %v", matched)
	}

	matched, err = ReadmeKeywordMatch("This contains no keywords", concepts)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(matched) != 0 {
		t.Errorf("Expected 0, got %d", len(matched))
	}
}

// func TestReadmeAiMatch(t *testing.T) {
// 	concepts := Concepts{
// 		"javascript": []string{"npm", "node", "javascript"},
// 		"ruby":       []string{"gem", "ruby"},
// 		"golang":     []string{"go", "golang"},
// 	}
// 	examples := Examples{
// 		"gem is a package manager for ruby":       []string{"ruby"},
// 		"go is a programming language":            []string{"golang"},
// 		"npm is a package manager for javascript": []string{"javascript"},
// 	}
// 	_, err := ReadmeAiMatch("This contains open and read words", concepts, examples)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %s", err)
// 	}
// }

func TestGeneratePrompt(t *testing.T) {
	concepts := Concepts{
		"javascript": []string{"npm", "node", "javascript"},
		"ruby":       []string{"gem", "ruby"},
		"golang":     []string{"go", "golang"},
	}
	examples := Examples{
		"npm is a package manager for javascript": []string{"javascript"},
	}
	content := "This contains open and read words"
	prompt := generatePrompt(concepts, examples, content)
	expectedPrompt := `You are a software developer and you are reading a README file. You are interested if the README references any of the following topics:
javascript, ruby, golang
Classify the README file based on the topics it references.
Desired format: <comma_separated_list_of_keywords>
Text: npm is a package manager for javascript
Keywords: javascript
Text: This contains open and read words
Keywords:`
	if prompt == expectedPrompt {
		t.Errorf("Expected %s, got %s", expectedPrompt, prompt)
	}
}
