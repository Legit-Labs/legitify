package pagination

import (
	"github.com/Legit-Labs/legitify/internal/clients/pagination"
	"github.com/xanzy/go-gitlab"
)

type Optioner interface {
	Done(resp interface{}) bool
	Advance(resp interface{}, opts interface{})
}

type GLOpts = *gitlab.ListOptions
type GLResp = *gitlab.Response
type glOptioner struct {
}

func (gl *glOptioner) Done(resp interface{}) bool {
	r := resp.(GLResp)
	return r.CurrentPage == r.TotalPages
}
func (gl *glOptioner) Advance(resp interface{}, opts interface{}) {
	r := resp.(GLResp)
	o := opts.(GLOpts)
	o.Page = r.NextPage
}

func New[ApiRetT any](fn interface{}, opts interface{}) *pagination.Basic[ApiRetT, GLOpts, GLResp] {
	return pagination.New[ApiRetT, GLOpts, GLResp](fn, opts, &glOptioner{})
}
func NewMapper[ApiRetT any, UserRetT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT) *pagination.MappedPager[ApiRetT, UserRetT, GLOpts, GLResp] {
	return pagination.NewMapper[ApiRetT, UserRetT, GLOpts, GLResp](fn, opts, mapper, &glOptioner{})
}
