package pagination

import (
	"github.com/Legit-Labs/legitify/internal/clients/pagination"
	"github.com/xanzy/go-gitlab"
)

type Optioner interface {
	Done(resp interface{}) bool
	Next(resp interface{}, opts interface{})
}

type GLOpts = *gitlab.ListOptions
type GLResp = *gitlab.Response
type glOptioner struct {
}

func (gh *glOptioner) Done(resp interface{}) bool {
	r := resp.(GLResp)
	return r.CurrentPage == r.TotalPages
}
func (gh *glOptioner) Next(resp interface{}, opts interface{}) {
	r := resp.(GLResp)
	o := opts.(GLOpts)
	o.Page = r.NextPage
}

func New[T any](fn interface{}, opts interface{}) *pagination.Basic[T, GLOpts, GLResp] {
	return pagination.New[T, GLOpts, GLResp](fn, opts, &glOptioner{})
}
func NewMapper[T any, U any](fn interface{}, opts interface{}, mapper func(T) []U) *pagination.MappedPager[T, U, GLOpts, GLResp] {
	return pagination.NewMapper[T, U, GLOpts, GLResp](fn, opts, mapper, &glOptioner{})
}
