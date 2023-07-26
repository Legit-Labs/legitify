package pagination

import (
	"reflect"

	"github.com/Legit-Labs/legitify/internal/clients/pagination"
	"github.com/google/go-github/v53/github"
)

type GHResp = *github.Response
type ghOptioner struct {
}

func (gh *ghOptioner) Done(resp interface{}) bool {
	r := resp.(GHResp)
	return r.NextPage == 0
}
func (gh *ghOptioner) Advance(resp interface{}, opts interface{}) {
	r := resp.(GHResp)
	p := reflect.ValueOf(opts).Elem()
	p.FieldByName("Page").SetInt(int64(r.NextPage))
}
func (gh *ghOptioner) OptionsIndex(fnInputsCount int, isVariadic bool) int {
	return fnInputsCount - 1
}

func New[ApiRetT any](fn interface{}, opts interface{}) *pagination.Basic[ApiRetT, GHResp] {
	return pagination.New[ApiRetT, GHResp](fn, opts, &ghOptioner{})
}
func NewMapper[ApiRetT any, UserRetT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT) *pagination.MappedPager[ApiRetT, UserRetT, GHResp] {
	return pagination.NewMapper[ApiRetT, UserRetT, GHResp](fn, opts, mapper, &ghOptioner{})
}
