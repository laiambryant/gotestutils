package pbtesting

import (
	"fmt"

	p "github.com/laiambryant/gotestutils/pbtesting/properties"
)

type InvalidPropertyError struct {
	predicate p.Predicate
}

func (i InvalidPropertyError) Error() string {
	return fmt.Sprintf("invalid property: %v", i.predicate)
}

type FunctionNotProvidedError struct{}

func (fnp FunctionNotProvidedError) Error() string {
	return "a function must be provided for the property-based test suite to work"
}
