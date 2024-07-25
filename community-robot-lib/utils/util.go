package utils

import (
	"io/ioutil"
	"os"
	"regexp"

	"sigs.k8s.io/yaml"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,6})$`)

func IsValidEmail(email string) bool {
	return emailRe.MatchString(email)
}

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := []byte(os.ExpandEnv(string(b)))

	return yaml.Unmarshal(content, cfg)
}
