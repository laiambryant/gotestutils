package properties

type Property interface {
	Verify(any) bool
}
