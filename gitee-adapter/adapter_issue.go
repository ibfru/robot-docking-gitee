package sdkadapter

import (
	"context"
	"github.com/opensourceways/go-gitee/gitee"
	"strconv"
)

func (c *ClientTarget) AddIssueComment(iss *IssueParameter) error {
	opt := gitee.PullRequestCommentPostParam{Body: iss.Comment}
	number, _ := strconv.ParseInt(iss.Number, 10, 32)
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(
		context.Background(), iss.Org, iss.Repo, int32(number), opt)
	return formatErr(err, "create comment of pr")
}

func (c *ClientTarget) DeleteIssueComment(iss *IssueParameter) error {
	var id int32
	if len(iss.CommentID) > 0 {
		i, e := strconv.Atoi(iss.CommentID)
		if e != nil {
			id = int32(i)
		}
	}

	_, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsCommentsId(
		context.Background(), iss.Org, iss.Repo, id, nil)
	return formatErr(err, "delete comment of pr")
}
