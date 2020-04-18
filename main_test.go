package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func Test_writeREADME(t *testing.T) {
	flextime.Fix(time.Date(2020, time.April, 19, 8, 00, 00, 0, time.UTC))
	type args struct {
		repos []Repo
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				repos: []Repo{
					{
						Name:        "test",
						Description: "test",
						URL:         "dummy://test.test",
						UpdatedAt:   time.Date(2020, time.April, 10, 7, 45, 00, 0, time.UTC),
						Stars:       1,
						Forks:       2,
						OpenIssues:  3,
					},
				},
			},
			wantW: `# Top Go ORMs [![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/awesome-go-orms)](https://goreportcard.com/report/github.com/d-tsuji/awesome-go-orms) [![Actions Status](https://github.com/d-tsuji/awesome-go-orms/workflows/CI/badge.svg)](https://github.com/d-tsuji/awesome-go-orms/actions)
A list of popular github projects related to Go ORM(Object-Relational Mapping) (ranked by stars automatically)
Please update **list.txt** (via Pull Request)

| Project Name | Stars | Forks | Open Issues | Description | Last Update |
| ------------ | ----- | ----- | ----------- | ----------- | ----------- |
| [test](dummy://test.test) | 1 | 2 | 3 | test | 2020-04-10 07:45:00 |

*Last Automatic Update: 2020-04-19T08:00:00Z*`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := writeREADME(w, tt.args.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeREADME() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("writeREADME() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_getURL(t *testing.T) {
	type args struct {
		repoURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{"https://github.com/d-tsuji/awesome-go-orms"},
			want:    "https://api.github.com/repos/d-tsuji/awesome-go-orms",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getURL(tt.args.repoURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
