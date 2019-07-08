%COMMENT%

package %PACKAGE%

// Bool represents a K boolean value
type Bool struct {
	Value bool
}

// BoolTrue represents a boolean value with value true
var BoolTrue = &Bool{Value: true}

// BoolFalse represents a boolean value with value false
var BoolFalse = &Bool{Value: false}

// CastToBool converts K Bool to Go bool, if possible.
func CastToBool(c KReference) (bool, bool) {
	b, typeOk := c.(*Bool)
	if !typeOk {
		return false, false
	}
	return b.Value, true
}

// ToKBool converts Go bool to K Bool.
func ToKBool(b bool) KReference {
	if b {
		return BoolTrue
	}
	return BoolFalse
}

// IsBool checks if the argument is a bool reference
func IsBool(c KReference) bool {
	_, typeOk := c.(*Bool)
	return typeOk
}

// IsTrue checks if argument is identical to the K Bool with the value true
func IsTrue(c KReference) bool {
	if b, typeOk := c.(*Bool); typeOk {
		return b.Value
	}
	return false
}
