package pagination

type Optioner interface {
	Done(resp interface{}) bool
	Next(resp interface{}, opts interface{})
}
