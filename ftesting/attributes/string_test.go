package attributes

import (
	"reflect"
	"testing"

	ctesting "github.com/laiambryant/gotestutils/ctesting"
)

// Test StringAttributes (suite)
func TestStringAttributes(t *testing.T) {
	var suite []ctesting.CharacterizationTest[bool]

	// GetAttributes tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}
		got := attr.GetAttributes()
		expected := StringAttributes{MinLen: 1, MaxLen: 5, Prefix: "pre", Suffix: "suf", Contains: "mid", UniqueChars: true}
		return reflect.DeepEqual(got, expected), nil
	}))

	// GetReflectType tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{}
		got := attr.GetReflectType()
		expected := reflect.TypeOf("")
		return got == expected, nil
	}))

	// GetDefaultImplementation tests
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{}
		got := attr.GetDefaultImplementation()
		return got != nil && reflect.TypeOf(got) == reflect.TypeOf(attr), nil
	}))

	// GetRandomValue tests - basic functionality
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 1, MaxLen: 10}
		got := attr.GetRandomValue()
		return got != nil, nil
	}))

	// Edge case: Max length zero
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 0, MaxLen: 0}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	// Edge case: Min length negative
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: -5, MaxLen: 10}
		result := attr.GetRandomValue()
		return result != nil, nil
	}))

	// Edge case: Min greater than max
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 10, MaxLen: 5}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) > 10 {
			return false, nil
		}
		return true, nil
	}))

	// With prefix and suffix
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 5, MaxLen: 10, Prefix: "pre_", Suffix: "_suf"}
		result := attr.GetRandomValue()
		if str, ok := result.(string); !ok || len(str) < 8 {
			return false, nil
		}
		return true, nil
	}))

	// Default max length test
	suite = append(suite, ctesting.NewCharacterizationTest(true, nil, func() (bool, error) {
		attr := StringAttributes{MinLen: 0, MaxLen: 0}
		result := attr.GetRandomValue()
		if result == nil {
			return false, nil
		}
		str := result.(string)
		return len(str) <= 10, nil
	}))

	// Custom allowed runes test
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

	results, _ := ctesting.VerifyCharacterizationTestsAndResults(t, suite, true)
	for i, passed := range results {
		if !passed {
			t.Fatalf("StringAttributes test %d failed", i+1)
		}
	}
}

func TestStringAttributes_MaxLenZero(t *testing.T) {
	attr := StringAttributes{MinLen: 0, MaxLen: 0}
	result := attr.GetRandomValue()
	if str, ok := result.(string); !ok || len(str) > 10 {
		t.Errorf("Expected empty or short string, got %v", result)
	}
}

func TestStringAttributes_MinLenNegative(t *testing.T) {
	attr := StringAttributes{MinLen: -5, MaxLen: 10}
	result := attr.GetRandomValue()
	if result == nil {
		t.Error("Expected string result, got nil")
	}
}

func TestStringAttributes_MinGreaterThanMax(t *testing.T) {
	attr := StringAttributes{MinLen: 10, MaxLen: 5}
	result := attr.GetRandomValue()
	if str, ok := result.(string); !ok {
		t.Errorf("Expected string result, got %T", result)
	} else if len(str) > 10 {
		t.Errorf("Expected length <= 10, got %d", len(str))
	}
}

func TestStringAttributes_WithPrefixSuffix(t *testing.T) {
	attr := StringAttributes{MinLen: 5, MaxLen: 10, Prefix: "pre_", Suffix: "_suf"}
	result := attr.GetRandomValue()
	if str, ok := result.(string); !ok {
		t.Errorf("Expected string result, got %T", result)
	} else {
		if len(str) < 8 {
			t.Errorf("Expected string with prefix and suffix, got %s", str)
		}
	}
}

func TestStringAttributes_DefaultMaxLen(t *testing.T) {
	attrs := StringAttributes{
		MinLen: 0,
		MaxLen: 0,
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil string")
	}

	str := result.(string)
	if len(str) > 10 {
		t.Errorf("Expected max length 10 (default), got %d", len(str))
	}
}

func TestStringAttributes_CustomAllowedRunes(t *testing.T) {
	attrs := StringAttributes{
		MinLen:       5,
		MaxLen:       5,
		AllowedRunes: []rune{'a', 'b', 'c'},
	}

	result := attrs.GetRandomValue()
	if result == nil {
		t.Fatal("Expected non-nil string")
	}

	str := result.(string)
	for _, r := range str {
		if r != 'a' && r != 'b' && r != 'c' {
			t.Errorf("Found disallowed rune %c in string %s", r, str)
		}
	}
}
