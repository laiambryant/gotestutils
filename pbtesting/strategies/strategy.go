package property

type Strategy interface {
	Execute() any
}

type IntegerStrategy struct {
	min  int
	max  int
	size int
}

func (i IntegerStrategy) Execute() {
}
