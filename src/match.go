package src

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func ReadmeMatch(content string, concepts Concepts, examples Examples, aiMatch bool) ([]string, error) {
	matchedConcepts, err := ReadmeKeywordMatch(content, concepts)
	if err != nil {
		return nil, err
	}
	if len(matchedConcepts) == 0 {
		return matchedConcepts, nil
	}
	if !aiMatch {
		return matchedConcepts, nil
	}
	aiMatchedConcepts, err := ReadmeAiMatch(content, concepts, examples)
	if err != nil {
		return nil, err
	}
	return aiMatchedConcepts, nil
}

// Parses a content and sees if any of the keywords match in it
func ReadmeKeywordMatch(content string, concepts Concepts) ([]string, error) {
	matchedConcepts := []string{}
	for concept, keywords := range concepts {
		for _, searchWord := range keywords {
			// Check if the content contains the keyword
			if strings.Contains(content, searchWord) {
				matchedConcepts = append(matchedConcepts, concept)
				break
			}
		}
	}
	return matchedConcepts, nil
}

// Parses a content and sees if any of the keywords match in it using AI
func ReadmeAiMatch(content string, concepts Concepts, examples Examples) ([]string, error) {
	prompt := generatePrompt(concepts, examples, content)
	// Run the AI model
	aiClient, err := NewOpenAiClient()
	if err != nil {
		return nil, err
	}
	completion, err := aiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return nil, err
	}

	completionContent := completion.Choices[0].Message.Content
	return strings.Split(completionContent, ","), nil
}

// Generates a system prompt for the AI Model
func generatePrompt(concepts Concepts, examples Examples, content string) string {
	// This method may be hard to follow
	// For an example of a generated prompt, take a look at match_test.go
	aiPrompt := `You are a software developer and you are reading a README file. You are interested if the README references any of the following topics:
%s
Classify the README file based on the topics it references.
Desired format: comma_separated_list_of_keywords`
	// Generate a list of keywords from the keywords map
	conceptList := []string{}
	for concept := range concepts {
		conceptList = append(conceptList, concept)
	}
	// Format the keywords as a comma-separated list
	formattedKeywords := strings.Join(conceptList, ", ")
	aiPrompt = fmt.Sprintf(aiPrompt, formattedKeywords)
	// Few shot examples
	examplesPrompt := `
Text: %s
Keywords: %s`
	for example, matchedConcepts := range examples {
		formattedMatchedConcepts := strings.Join(matchedConcepts, ", ")
		aiPrompt += fmt.Sprintf(examplesPrompt, example, formattedMatchedConcepts)
	}
	aiPrompt += "\nText: " + content
	aiPrompt += "\nKeywords:"
	return aiPrompt
}

type Concepts map[string][]string

// Loads the keywords from a JSON file
func LoadKeywords(keywordFile string) (Concepts, error) {
	// Read the file
	// Parse the JSON
	// Return the keywords
	file, err := os.ReadFile(keywordFile)
	if err != nil {
		return nil, err
	}
	var concepts Concepts
	err = json.Unmarshal(file, &concepts)
	if err != nil {
		return nil, err
	}
	return concepts, nil
}

type Examples map[string][]string

// Loads the few shot examples from a JSON file
func LoadExamples(exampleFile string) (Examples, error) {
	// Read the file
	file, err := os.ReadFile(exampleFile)
	if err != nil {
		return nil, err
	}
	var examples Examples
	err = json.Unmarshal(file, &examples)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func NewOpenAiClient() (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return client, nil
}
