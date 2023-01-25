package pagination

import (
	"github.com/Legit-Labs/legitify/internal/clients/pagination"
	"github.com/google/go-github/v49/github"
)

type Optioner interface {
	Done(resp interface{}) bool
	Advance(resp interface{}, opts interface{})
}

type GHOpts = *github.ListOptions
type GHResp = *github.Response
type ghOptioner struct {
}

func (gh *ghOptioner) Done(resp interface{}) bool {
	r := resp.(GHResp)
	return r.NextPage == 0
}
func (gh *ghOptioner) Advance(resp interface{}, opts interface{}) {
	r := resp.(GHResp)
	o := opts.(GHOpts)
	o.Page = r.NextPage
}

func New[T any](fn interface{}, opts interface{}) *pagination.Basic[T, GHOpts, GHResp] {
	return pagination.New[T, GHOpts, GHResp](fn, opts, &ghOptioner{})
}
func NewMapper[T any, U any](fn interface{}, opts interface{}, mapper func(T) []U) *pagination.MappedPager[T, U, GHOpts, GHResp] {
	return pagination.NewMapper[T, U, GHOpts, GHResp](fn, opts, mapper, &ghOptioner{})
}
