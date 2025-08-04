package ptesting

import (
	s "github.com/laiambryant/gotestutils/ptesting/strategies"
	gtu "github.com/laiambryant/gotestutils/testing"
)

type PTest[retT comparable, argT any] struct {
	Func           gtu.PTestFunc[retT, argT]
	ValidationFunc func(...argT) bool
	Property       []s.Property
	iterations     uint
}
