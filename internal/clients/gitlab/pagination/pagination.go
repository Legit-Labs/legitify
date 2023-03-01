package pagination

import (
	"reflect"

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
	p := reflect.ValueOf(opts).Elem()
	p.FieldByName("Page").SetInt(int64(r.NextPage))
}

func (gl *glOptioner) OptionsIndex(fnInputsCount int, isVariadic bool) int {
	index := fnInputsCount - 1
	if isVariadic {
		index--
	}
	return index
}

func New[ApiRetT any](fn interface{}, opts interface{}) *pagination.Basic[ApiRetT, GLResp] {
	return pagination.New[ApiRetT, GLResp](fn, opts, &glOptioner{})
}
func NewMapper[ApiRetT any, UserRetT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT) *pagination.MappedPager[ApiRetT, UserRetT, GLResp] {
	return pagination.NewMapper[ApiRetT, UserRetT, GLResp](fn, opts, mapper, &glOptioner{})
}
