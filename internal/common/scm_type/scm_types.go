package scm_type

import "fmt"

type ScmType = string

const (
	GitHub ScmType = "github"
	GitLab ScmType = "gitlab"
)

var All = []ScmType{
	GitHub,
	GitLab,
}

func Validate(scmType ScmType) error {
	for _, e := range All {
		if e == scmType {
			return nil
		}
	}

	return fmt.Errorf("invalid scm type %s", scmType)
}
