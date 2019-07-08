%COMMENT%

package %PACKAGE%

// Bottom is a K item that contains no data
type Bottom struct {
}

// InternedBottom is usually used as a dummy object
var InternedBottom = &Bottom{}

// IsBottom returns true if reference points to bottom
func IsBottom(ref KReference) bool {
	_, typeOk := ref.(*Bottom)
	return typeOk
}

// Float is a KObject representing a float in K
type Float struct {
	Value float32
}

// GetFloatObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetFloatObject(ref KReference) (*Float, bool) {
	castObj, typeOk := ref.(*Float)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// IsFloat returns true if reference points to a float
func IsFloat(ref KReference) bool {
	_, typeOk := ref.(*Float)
	return typeOk
}

// MInt is a KObject representing a machine integer in K
type MInt struct {
	Value int32
}

// IsMInt returns true if reference points to a string buffer
func IsMInt(ref KReference) bool {
	_, typeOk := ref.(*MInt)
	return typeOk
}

// InjectedKLabel is a KObject representing an InjectedKLabel item in K
type InjectedKLabel struct {
	Label KLabel
}

// NewInjectedKLabel creates a new InjectedKLabel object and returns the reference.
func (ms *ModelState) NewInjectedKLabel(label KLabel) KReference {
	return &InjectedKLabel{Label: label}
}
