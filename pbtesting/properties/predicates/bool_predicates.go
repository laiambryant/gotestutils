package predicates

type BoolMustBeTrue struct{}
type BoolMustBeFalse struct{}

func (p BoolMustBeTrue) Verify(v any) bool  { b, ok := v.(bool); return !ok || b }
func (p BoolMustBeFalse) Verify(v any) bool { b, ok := v.(bool); return !ok || !b }
