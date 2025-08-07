package pbtesting

import (
	"fmt"

	p "github.com/laiambryant/gotestutils/pbtesting/properties"
)

type InvalidPropertyError struct {
	property p.Property
}

func (i InvalidPropertyError) Error() string {
	return fmt.Sprintf("invalid property: %v", i.property)
}
