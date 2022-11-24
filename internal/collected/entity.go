package collected

type Entity interface {
	ViolationEntityType() string
	CanonicalLink() string
	Name() string
	ID() int64
}
