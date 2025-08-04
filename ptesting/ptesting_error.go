package ptesting

import (
	"fmt"

	s "github.com/laiambryant/gotestutils/ptesting/strategies"
)

type InvalidPropertyError struct {
	property s.Property
}

func (i InvalidPropertyError) Error() string {
	return fmt.Sprintf("invalid property: %v", i.property)
}
