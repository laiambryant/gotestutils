package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

// Test StringAttributes (suite)
func TestStringAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}
		got := attr.GetAttributes()
		expected := StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}
		return reflect.DeepEqual(got, expected), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf("")
		return got == expected, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 1, MaxLen: 10}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 0, MaxLen: 0}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: -5, MaxLen: 10}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 10, MaxLen: 5}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 5, MaxLen: 10, Prefix: "pre_", Suffix: "_suf"}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) < 8 {
			return false, nil
		}
		return true, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 0, MaxLen: 0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		str := result.(string)
		return len(str) <= 10, nil
	}))

	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{
			MinLen:       5,
			MaxLen:       5,
			AllowedRunes: []rune{'a', 'b', 'c'},
		}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		str := result.(string)
		for _, r := range str {
			if r != 'a' && r != 'b' && r != 'c' {
				return false, nil
			}
		}
		return true, nil
	}))

	// Additional tests from individual test functions
	// TestStringAttributes_MaxLenZero
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 0, MaxLen: 0}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	// TestStringAttributes_MinLenNegative
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: -5, MaxLen: 10}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	// TestStringAttributes_MinGreaterThanMax
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 10, MaxLen: 5}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok {
			return false, nil
		} else if len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	// TestStringAttributes_WithPrefixSuffix
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 5, MaxLen: 10, Prefix: "pre_", Suffix: "_suf"}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok {
			return false, nil
		} else {
			return len(str) >= 8, nil
		}
	}))

	// TestStringAttributes_DefaultMaxLen
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attrs := StringAttributes{
			MinLen: 0,
			MaxLen: 0,
		}
		result := attrs.GetRandomValue()
		if result == nil {
			return false, nil
		}
		str := result.(string)
		return len(str) <= 10, nil
	}))

	// TestStringAttributes_CustomAllowedRunes (already covered by existing test)

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("StringAttributes test %d failed", i+1)
		}
	}
}
