%COMMENT%

package %PACKAGE_MODEL%

// BoolTrue ... K boolean value with value true
var BoolTrue = &Bool{Value: true}

// BoolFalse ... K boolean value with value false
var BoolFalse = &Bool{Value: false}

// ToBool ... Convert Go bool to K Bool
func ToBool(b bool) *Bool {
	if b {
		return BoolTrue
	}
	return BoolFalse
}

// IsTrue ... Checks if argument is identical to the K Bool with the value true
func IsTrue(c K) bool {
	if b, typeOk := c.(*Bool); typeOk {
		return b.Value
	}
	return false
}
