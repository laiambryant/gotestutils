package strategy

type Property interface {
	Execute()
}

type IntegerProperty struct {
	numberOfIteration int
	min               int
	max               int
}

func (i IntegerProperty) Execute() {}
