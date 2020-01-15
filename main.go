package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	head = `# Top Go ORMs
A list of popular github projects related to Go ORM(Object-Relational Mapping) (ranked by stars automatically)
Please update **list.txt** (via Pull Request)

| Project Name | Stars | Forks | Open Issues | Description | Last Update |
| ------------ | ----- | ----- | ----------- | ----------- | ----------- |
`
	tail = "\n*Last Automatic Update: %v*"

	gitHubUrl = "https://github.com/"
)

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
	content, err := ioutil.ReadFile("list.txt")
	if err != nil {
		return fmt.Errorf("file read: %s", err)
	}
	lines := strings.Split(string(content), "\n")

	var repos []Repo
	for _, line := range lines {
		if strings.HasPrefix(line, gitHubUrl) {
			var r Repo
			fn := func() error {
				apiUrl, err := getApiUrl(line)
				if err != nil {
					return fmt.Errorf("get apiurl: %s", err)
				}
				resp, err := http.Get(apiUrl)
				if err != nil {
					return fmt.Errorf("http get request: %s", err)
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("response code: %d", resp.StatusCode)
				}
				if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
					return fmt.Errorf("json decode: %s", err)
				}
				return nil
			}()
			if err := fn; err != nil {
				return err
			}
			repos = append(repos, r)
		} else if line != "" {
			log.Printf("url: %s is not supported\n", line)
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})
	if err := writeREADME(repos); err != nil {
		return fmt.Errorf("writeREADME: %s", err)
	}
	return nil
}

func getApiUrl(repoUrl string) (string, error) {
	repoUrl = strings.TrimFunc(repoUrl, func(r rune) bool {
		return unicode.IsSpace(r) || (r == rune('/'))
	})
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://api.github.com/%s", "repos" + parsedUrl.Path), nil
}

func writeREADME(repos []Repo) error {
	readme, err := os.OpenFile("README.md", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer readme.Close()
	readme.WriteString(head)
	for _, repo := range repos {
		readme.WriteString(
			fmt.Sprintf("| [%s](%s) | %d | %d | %d | %s | %v |\n",
				repo.Name,
				repo.URL,
				repo.Stars,
				repo.Forks,
				repo.OpenIssues,
				repo.Description,
				repo.UpdatedAt.Format("2006-01-02 15:04:05")))
	}
	readme.WriteString(fmt.Sprintf(tail, time.Now().Format(time.RFC3339)))

	return nil
}
