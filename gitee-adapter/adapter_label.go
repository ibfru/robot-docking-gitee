package sdkadapter

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/opensourceways/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"
)

func formatErr(err error, doWhat string) error {
	if err == nil {
		return err
	}

	var msg []byte
	var v gitee.GenericSwaggerError
	if errors.As(err, &v) {
		msg = v.Body()
	}

	return fmt.Errorf("failed to %s, err: %s, msg: %q", doWhat, err.Error(), msg)
}

func (c *ClientTarget) GetPRLabels(pr *PRParameter) (*sets.String, error) {
	lc := sets.NewString()

	p := int32(1)
	opt := gitee.GetV5ReposOwnerRepoPullsNumberLabelsOpts{}
	for {
		opt.Page = optional.NewInt32(p)
		number, _ := strconv.ParseInt(pr.Number, 10, 32)
		ls, _, err := c.ac.PullRequestsApi.GetV5ReposOwnerRepoPullsNumberLabels(
			context.Background(), pr.Org, pr.Repo, int32(number), &opt)
		if err != nil {
			return nil, formatErr(err, "list labels of pr")
		}

		j := len(ls)
		if j == 0 {
			break
		}

		for i := 0; i < j; i++ {
			lc.Insert(ls[i].Name)
		}

		p++
	}

	return &lc, nil
}

func (c *ClientTarget) AddPRLabels(pr *PRParameter) error {
	opt := gitee.PullRequestLabelPostParam{Body: pr.Labels}
	number, _ := strconv.ParseInt(pr.Number, 10, 32)
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberLabels(
		context.Background(), pr.Org, pr.Repo, int32(number), opt)
	return formatErr(err, "add multi label for pr")
}

func (c *ClientTarget) DeletePRLabels(pr *PRParameter) error {
	if pr.Labels == nil || len(pr.Labels) == 0 {
		return fmt.Errorf("can not found laebl to remove")
	}

	// gitee's bug, it can't deal with the label which includes '/'
	label := strings.Replace(strings.Join(pr.Labels, ","), "/", "%2F", -1)

	number, _ := strconv.ParseInt(pr.Number, 10, 32)
	v, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsLabel(
		context.Background(), pr.Org, pr.Repo, int32(number), label, nil)

	if err == nil || (v != nil && v.StatusCode == 404) {
		return nil
	}
	return formatErr(err, "remove label of pr")
}

func (c *ClientTarget) GetRepoLabels(lp *LabelParameter) (*sets.String, error) {
	lc := sets.NewString()

	ls, _, err := c.ac.LabelsApi.GetV5ReposOwnerRepoLabels(context.Background(), lp.Org, lp.Repo, nil)
	if j := len(ls); j != 0 {
		for i := 0; i < j; i++ {
			lc.Insert(ls[i].Name)
		}
	}

	return &lc, formatErr(err, "get repo labels")
}

func (c *ClientTarget) AddRepoLabels(lp *LabelParameter) error {
	if lp.Color == "" {
		v := rand.New(rand.NewSource(time.Now().Unix()))
		lp.Color = fmt.Sprintf("%02x%02x%02x", v.Intn(255), v.Intn(255), v.Intn(255))
	}
	param := gitee.LabelPostParam{
		Name:  lp.Name,
		Color: lp.Color,
	}

	_, _, err := c.ac.LabelsApi.PostV5ReposOwnerRepoLabels(context.Background(), lp.Org, lp.Repo, param)

	return formatErr(err, "create a repo label")
}

func (c *ClientTarget) AddIssueLabels(iss *IssueParameter) error {
	opt := gitee.PullRequestLabelPostParam{Body: iss.Labels}
	number, _ := strconv.ParseInt(iss.Number, 10, 32)
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberLabels(
		context.Background(), iss.Org, iss.Repo, int32(number), opt)
	return formatErr(err, "add multi label for pr")
}
