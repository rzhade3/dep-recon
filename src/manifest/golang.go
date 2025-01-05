package manifest

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Golang struct {
	// The URL of the package manager's registry
	RegistryURL string
	// The filepath to the dependency file
	DependencyFilePath string
}

var DefaultGolangRegistryURL = "https://pkg.go.dev"

func (g Golang) PullDependencyReadme(dependency, version string) (string, error) {
	// Golang doesn't have a great way to pull a README from the package manager's registry
	// Instead, we will return the <meta name="description>" from the pkg.go.dev homepage for
	// The dependency. This is not a great solution, but it's the best we can do for now
	goPkgUrl, err := url.JoinPath(g.RegistryURL, dependency)
	if err != nil {
		return "", err
	}
	// Fetch the README from the goPkgUrl
	res, err := http.Get(goPkgUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", err
	}
	// Parse the body as HTML
	// Find the <meta name="description>" tag
	// Return the content of the <meta name="description>" tag
	htmlBody, err := html.Parse(res.Body)
	if err != nil {
		return "", err
	}
	var crawler func(node *html.Node) string
	crawler = func(node *html.Node) string {
		if node.Type == html.ElementNode && node.Data == "meta" {
			for _, attr := range node.Attr {
				if attr.Key == "name" && attr.Val == "description" || attr.Val == "Description" {
					for _, attr := range node.Attr {
						if attr.Key == "content" {
							return attr.Val
						}
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if result := crawler(c); result != "" {
				return result
			}
		}
		return ""
	}
	description := crawler(htmlBody)
	return description, nil
}

func (g Golang) GetEcosystem() string {
	return "golang"
}

func (g Golang) ListDependencies() (DependencyList, error) {
	// Open the go.mod file
	file, err := os.Open(g.DependencyFilePath)
	if err != nil {
		return DependencyList{}, err
	}
	defer file.Close()
	mod_regex := regexp.MustCompile(`([\w.\-/]+)\s+(v[\d+.\-a-zA-Z]+)`)
	depList := DependencyList{
		DevDependencies: map[string]string{},
		Dependencies:    map[string]string{},
	}
	// Scan the file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "// indirect") {
			continue
		}
		match := mod_regex.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		packageName := match[1]
		packageVersion := match[2]
		depList.Dependencies[packageName] = packageVersion
	}
	return depList, nil
}

func (g Golang) GetDependencyFilePath() string {
	return g.DependencyFilePath
}
