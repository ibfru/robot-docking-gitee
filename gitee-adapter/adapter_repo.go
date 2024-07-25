package sdkadapter

import (
	"context"
	"encoding/json"
)

func (c *ClientTarget) GetRepoContentsByPath(org, repo, path string) ([]*ContentInfo, error) {

	txt, _, err := c.ac.RepositoriesApi.GetV5ReposOwnerRepoContentsPath(
		context.Background(), org, repo, path, nil)
	if err == nil {
		if b, e := json.Marshal(txt); e == nil {
			ans := make([]*ContentInfo, 1)
			var cnt ContentInfo
			e = json.Unmarshal(b, &cnt)
			if e != nil {
				return nil, e
			}
			ans[0] = &cnt
			return ans, nil
		}
	}
	return nil, nil
}

func (c *ClientTarget) ListCollaborator(org, repo string) ([]string, error) {

	return nil, nil
}
