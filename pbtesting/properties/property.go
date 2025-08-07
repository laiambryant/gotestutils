package properties

type Property interface {
	Verify(any) bool
}

type IntegerMax struct {
	max int
}

func (i IntegerMax) Verify(value int) bool {
	return i.max < value
}

type IntegerMin struct {
	min int
}

func (i IntegerMin) Verify(value int) bool {
	return value > i.min
}

type IntegerSigned struct {
	signed bool
}

func (i IntegerSigned) Verify(value int) bool {
	return (value > 0 && i.signed) || (value < 0)
}
