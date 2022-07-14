package githubcollected

type CollectedEntity interface {
	ViolationEntityType() string
	CanonicalLink() string
	Name() string
	ID() int64
}
