package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/Songmu/flextime"
)

const (
	head = `# Top Go ORMs [![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/awesome-go-orms)](https://goreportcard.com/report/github.com/d-tsuji/awesome-go-orms) [![Actions Status](https://github.com/d-tsuji/awesome-go-orms/workflows/CI/badge.svg)](https://github.com/d-tsuji/awesome-go-orms/actions)
A list of popular github projects related to Go ORM(Object-Relational Mapping) (ranked by stars automatically)
Please update **list.txt** (via Pull Request)

| Project Name | Stars | Forks | Open Issues | Description | Last Update |
| ------------ | ----- | ----- | ----------- | ----------- | ----------- |
`
	tail         = "\n*Last Automatic Update: %v*"
	listFileName = "list.txt"
)

// Repo is the structure that represents the schema of the github api.
// See below: https://developer.github.com/v3/repos/
type Repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"html_url"`
	UpdatedAt   time.Time `json:"updated_at"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	OpenIssues  int       `json:"open_issues_count"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	f, err := os.Open(listFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	sc := bufio.NewScanner(f)

	var repos []Repo
	for sc.Scan() {
		repoName := sc.Text()
		if repoName != "" {
			log.Printf("URL: %s is not supported\n", repoName)
			continue
		}
		if strings.HasPrefix(repoName, "https://github.com/") {
			r, err := fetchRepo(repoName)
			if err != nil {
				return fmt.Errorf("fetch repo: %w", err)
			}
			repos = append(repos, *r)
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("scan %s: %w", listFileName, err)
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})

	readme, err := os.Create("README.md")
	if err != nil {
		return err
	}
	defer func() {
		_ = readme.Close()
	}()
	writeREADME(readme, repos)
	return nil
}

func fetchRepo(name string) (*Repo, error) {
	var r *Repo
	apiURL, err := getURL(name)
	if err != nil {
		return nil, fmt.Errorf("get URL for api call: %w", err)
	}
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("http get request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code: %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return r, nil
}

func getURL(repoURL string) (string, error) {
	repoURL = strings.TrimFunc(repoURL, func(r rune) bool {
		return unicode.IsSpace(r) || (r == rune('/'))
	})
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/repos%s", parsedURL.Path), nil
}

func writeREADME(w io.Writer, repos []Repo) {
	fmt.Fprint(w, head)
	for _, repo := range repos {
		fmt.Fprintf(w, "| [%s](%s) | %d | %d | %d | %s | %v |\n",
			repo.Name,
			repo.URL,
			repo.Stars,
			repo.Forks,
			repo.OpenIssues,
			repo.Description,
			repo.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintf(w, tail, flextime.Now().Format(time.RFC3339))
}
