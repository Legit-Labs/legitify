package pagination

type Optioner interface {
	Done(resp interface{}) bool
	Advance(resp interface{}, opts interface{})
	OptionsIndex(fnInputsCount int, isVariadic bool) int
}
