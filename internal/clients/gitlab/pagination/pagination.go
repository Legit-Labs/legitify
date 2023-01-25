package pagination

import (
	"github.com/Legit-Labs/legitify/internal/clients/pagination"
	"github.com/xanzy/go-gitlab"
)

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

func (gl *glOptioner) OptionsIndex(fnInputsCount int, isVariadic bool) int {
	index := fnInputsCount - 1
	if isVariadic {
		index--
	}
	return index
}

func New[ApiRetT any](fn interface{}, opts interface{}) *pagination.Basic[ApiRetT, GLOpts, GLResp] {
	return pagination.New[ApiRetT, GLOpts, GLResp](fn, opts, &glOptioner{})
}
func NewMapper[ApiRetT any, UserRetT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT) *pagination.MappedPager[ApiRetT, UserRetT, GLOpts, GLResp] {
	return pagination.NewMapper[ApiRetT, UserRetT, GLOpts, GLResp](fn, opts, mapper, &glOptioner{})
}
