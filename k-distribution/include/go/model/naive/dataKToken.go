%COMMENT%

package %PACKAGE_MODEL%

// KToken is a KObject representing a KToken item in K
type KToken struct {
	Value string
	Sort  Sort
}

// GetKTokenObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetKTokenObject(ref KReference) (*KToken, bool) {
	castObj, typeOk := ref.(*KToken)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// NewKToken creates a new object and returns the reference.
func (ms *ModelState) NewKToken(sort Sort, value string) KReference {
	return &KToken{Sort: sort, Value: value}
}

// NewKTokenConstant creates a new KToken constant.
func NewKTokenConstant(sort Sort, value string) KReference {
	return &KToken{Sort: sort, Value: value}
}
