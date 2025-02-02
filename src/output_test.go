package src

import "testing"

func TestJsonOutput(t *testing.T) {
	matchedDependencies := map[string][]string{
		"express": {"npm", "node", "javascript"},
		"ruby":    {"gem", "ruby"},
	}
	expectedOutput := `{"express":["npm","node","javascript"],"ruby":["gem","ruby"]}`
	output := JsonOutput(matchedDependencies)
	if output != expectedOutput {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestTestOutput(t *testing.T) {
	matchedDependencies := map[string][]string{
		"express": {"npm", "node", "javascript"},
		"ruby":    {"gem", "ruby"},
	}
	expectedOutput := `express matched for [npm node javascript]
ruby matched for [gem ruby]`
	output := TextOutput(matchedDependencies)
	if output != expectedOutput {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}
