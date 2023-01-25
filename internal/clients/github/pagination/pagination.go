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

func New[ApiRetT any](fn interface{}, opts interface{}) *pagination.Basic[ApiRetT, GHOpts, GHResp] {
	return pagination.New[ApiRetT, GHOpts, GHResp](fn, opts, &ghOptioner{})
}
func NewMapper[ApiRetT any, UserRetT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT) *pagination.MappedPager[ApiRetT, UserRetT, GHOpts, GHResp] {
	return pagination.NewMapper[ApiRetT, UserRetT, GHOpts, GHResp](fn, opts, mapper, &ghOptioner{})
}
