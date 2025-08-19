package predicates

import "testing"

func TestStringProperties(t *testing.T) {
    assertProp(t, StringLenMin{Min: 3}, "ab", false)
    assertProp(t, StringLenMin{Min: 3}, "abc", true)
    assertProp(t, StringLenMax{Max: 3}, "abcd", false)
    assertProp(t, StringLenMax{Max: 3}, "abc", true)
    assertProp(t, StringLenRange{Min: 2, Max: 3}, "a", false)
    assertProp(t, StringLenRange{Min: 2, Max: 3}, "ab", true)
    assertProp(t, StringRegex{Pattern: "^a.+z$"}, "abz", true)
    assertProp(t, StringRegex{Pattern: "^a.+z$"}, "ax", false)
    assertProp(t, StringRegex{Pattern: "((a){10000})"}, 1, false)
    assertProp(t, StringPrefix{Prefix: "pre"}, "prefix", true)
    assertProp(t, StringPrefix{Prefix: "pre"}, "xprefix", false)
    assertProp(t, StringSuffix{Suffix: "suf"}, "endsuf", true)
    assertProp(t, StringSuffix{Suffix: "suf"}, "sufend", false)
    assertProp(t, StringContains{Substr: "mid"}, "amidb", true)
    assertProp(t, StringContains{Substr: "mid"}, "none", false)
}
