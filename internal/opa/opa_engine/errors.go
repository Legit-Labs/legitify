package opa_engine

import "fmt"

type ErrPolicyLoad struct {
	loaderError error
}

func NewErrPolicyLoad(err error) error {
	return &ErrPolicyLoad{err}
}

func (e *ErrPolicyLoad) Error() string {
	return fmt.Sprintf("load: %v", e.loaderError)
}

type ErrNoPolicies struct {
	policyPaths []string
}

func NewErrNoPolicies(policyPaths []string) error {
	return &ErrNoPolicies{policyPaths}
}

func (e *ErrNoPolicies) Error() string {
	return fmt.Sprintf("no policy .rego files found in %v", e.policyPaths)
}
