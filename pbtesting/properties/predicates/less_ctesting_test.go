package predicates

import (
	"testing"
	"github.com/laiambryant/gotestutils/ctesting"
)

func TestLessCharacterization(t *testing.T) {
	suite := []ctesting.CharacterizationTest[bool]{
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(1, 2), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(2, 1), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(2, 2), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(int8(1), int8(2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(int8(1), int16(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(uint8(1), uint8(2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(uint8(1), uint16(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(-1, 0), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(uint(1), int(2)), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less("a", "b"), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less("b", "a"), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less("a", "a"), nil }),
		ctesting.NewCharacterizationTest(true, nil, func() (bool, error) { return less(float32(1.1), float32(1.2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(float32(1.1), float64(1.2)), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less([]int{1}, []int{1}), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(nil, nil), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(nil, 1), nil }),
		ctesting.NewCharacterizationTest(false, nil, func() (bool, error) { return less(1, nil), nil }),
	}
	ctesting.VerifyCharacterizationTestsAndResults(t, suite, false)
}

