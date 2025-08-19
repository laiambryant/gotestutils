package predicates

import "regexp"

type StringLenMin struct{ Min int }
type StringLenMax struct{ Max int }
type StringLenRange struct{ Min, Max int }
type StringRegex struct{ Pattern string }
type StringPrefix struct{ Prefix string }
type StringSuffix struct{ Suffix string }
type StringContains struct{ Substr string }

func (p StringLenMin) Verify(v any) bool { s, ok := v.(string); return !ok || len(s) >= p.Min }
func (p StringLenMax) Verify(v any) bool { s, ok := v.(string); return !ok || len(s) <= p.Max }
func (p StringLenRange) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(s) >= p.Min && len(s) <= p.Max)
}
func (p StringRegex) Verify(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	if p.Pattern == "" {
		return true
	}
	re, err := regexp.Compile(p.Pattern)
	if err != nil {
		return true
	}
	return re.MatchString(s)
}
func (p StringPrefix) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(p.Prefix) == 0 || (len(s) >= len(p.Prefix) && s[:len(p.Prefix)] == p.Prefix))
}
func (p StringSuffix) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (len(p.Suffix) == 0 || (len(s) >= len(p.Suffix) && s[len(s)-len(p.Suffix):] == p.Suffix))
}
func (p StringContains) Verify(v any) bool {
	s, ok := v.(string)
	return !ok || (p.Substr == "" || (len(s) >= len(p.Substr) && (regexp.MustCompile(regexp.QuoteMeta(p.Substr)).FindStringIndex(s) != nil)))
}
