package sdkadapter

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/oauth2"

	"github.com/opensourceways/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"
)

type ClientTarget struct {
	ac *gitee.APIClient
}

var ct *ClientTarget
var onceClient sync.Once
var token string

func initialClient() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(context.Background(), ts)
	fmt.Println("fdauhsgdikahsfglkajhsgdlkajhsgdads")

	cfg := &gitee.Configuration{
		BasePath:      "https://gitee.com/api",
		DefaultHeader: make(map[string]string),
		UserAgent:     "robot",
		HTTPClient:    tc,
	}
	ct = &ClientTarget{
		ac: gitee.NewAPIClient(cfg),
	}
}

func GetClientInstance(t []byte) *ClientTarget {
	token = string(t)
	onceClient.Do(initialClient)
	return ct
}

type LabelParameter struct {
	Org    string
	Repo   string
	Name   string
	Color  string
	Extras any
}

type LabelClient interface {
	GetRepoLabels(lp *LabelParameter) (*sets.String, error)
	AddRepoLabels(lp *LabelParameter) error

	GetPRLabels(pr *PRParameter) (*sets.String, error)
	AddPRLabels(pr *PRParameter) error
	DeletePRLabels(pr *PRParameter) error

	//GetIssueLabels(iss *IssueParameter) (sets.String, error)
	AddIssueLabels(iss *IssueParameter) error
	//RemoveIssueLabels(iss *PRParameter) error
}

type PRParameter struct {
	Org       string `json:"org" binding:"orgValid"`
	Repo      string
	Number    string
	Labels    []string
	Comment   string
	CommentID string
	Reviewers []string
	Payload   any
	Extras    any
}

type PRClient interface {
	AddPRComment(pr *PRParameter) error
	DeletePRComment(pr *PRParameter) error

	AssignPR(pr *PRParameter) error
}

type IssueParameter struct {
	Org       string
	Repo      string
	Number    string
	Labels    []string
	Comment   string
	CommentID string
	Reviewers []string
	Payload   any
	Extras    any
}

type IssueClient interface {
	AddIssueComment(iss *IssueParameter) error
	DeleteIssueComment(iss *IssueParameter) error
}

type ContentInfo struct {
	Type        *string `json:"type"`
	Size        float32 `json:"size"`
	Name        *string `json:"name"`
	Path        *string `json:"path"`
	Sha         *string `json:"sha"`
	Url         *string `json:"url,omitempty"`
	HtmlUrl     *string `json:"html_url,omitempty"`
	DownloadUrl *string `json:"download_url,omitempty"`
	Content     *string `json:"content,omitempty"`
}

type RepoClient interface {
	GetRepoContentsByPath(org, repo, path string) ([]*ContentInfo, error)
	ListCollaborator(org, repo string) ([]string, error)
}
