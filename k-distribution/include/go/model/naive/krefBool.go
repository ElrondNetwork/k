%COMMENT%

package %PACKAGE_MODEL%

// BoolTrue ... K boolean value with value true
var BoolTrue = &Bool{Value: true}

// BoolFalse ... K boolean value with value false
var BoolFalse = &Bool{Value: false}

// CastKApply yields the cast object for a Bool reference, if possible.
func (ms *ModelState) CastToBool(ref KReference) (*Bool, bool) {
	cast, typeOk := ref.(*Bool)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// ToKBool ... Convert Go bool to K Bool
func ToKBool(b bool) KReference {
	if b {
		return BoolTrue
	}
	return BoolFalse
}

// IsTrue ... Checks if argument is identical to the K Bool with the value true
func IsTrue(c KReference) bool {
	if b, typeOk := c.(*Bool); typeOk {
		return b.Value
	}
	return false
}
