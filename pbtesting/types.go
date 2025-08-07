package pbtesting

type InputType string

const (
	InputTypeInt        InputType = "int"
	InputTypeString     InputType = "string"
	InputTypeBool       InputType = "bool"
	InputTypeFloat      InputType = "float"
	InputTypeAny        InputType = "any"
	InputTypeComparable InputType = "comparable"
	InputTypeError      InputType = "error"
	InputTypeStruct     InputType = "struct"
)

func getTypeForInputType(inputType InputType) string {
	switch inputType {
	case InputTypeInt:
		return "int"
	case InputTypeString:
		return "string"
	case InputTypeBool:
		return "bool"
	case InputTypeFloat:
		return "float64"
	case InputTypeAny:
		return "any"
	case InputTypeComparable:
		return "comparable"
	case InputTypeError:
		return "error"
	case InputTypeStruct:
		return "struct{}"
	default:
		return ""
	}
}
