%COMMENT%

package %PACKAGE_MODEL%

// KApply is a KObject representing a KApply item in K
type KApply struct {
	Label KLabel
	List  []KReference
}

// CastKApply returns true if argument is a KApply item.
// Also returns argument, for convenience.
func (ms *ModelState) CastKApply(ref KReference) (KReference, bool) {
	_, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		return nil, false
	}
	return ref, true
}

// CheckKApply returns true if argument is a KApply with given label and arity.
// Also returns argument, for convenience.
func (ms *ModelState) CheckKApply(ref KReference, expectedLabel KLabel, expectedArity int) (KReference, bool) {
	castObj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		return NullReference, false
	}
	if castObj.Label != expectedLabel {
		return NullReference, false
	}
	if len(castObj.List) != expectedArity {
		return NullReference, false
	}
	return ref, true
}

// KApplyLabel returns the label of a KApply item.
func (ms *ModelState) KApplyLabel(ref KReference) KLabel {
	castObj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		panic("KApplyLabel called for reference to item other than KApply")
	}
	return castObj.Label
}

// KApplyArity returns the arity of a KApply item (nr. of arguments)
func (ms *ModelState) KApplyArity(ref KReference) int {
	castObj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		panic("KApplyArity called for reference to item other than KApply")
	}
	return len(castObj.List)
}

// KApplyArg returns the nth argument in a KApply
func (ms *ModelState) KApplyArg(ref KReference, argIndex int) KReference {
	castObj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		panic("KApplyArity called for reference to item other than KApply")
	}
	return castObj.List[argIndex]
}

// GetKApplyObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetKApplyObject(ref KReference) (*KApply, bool) {
	castObj, typeOk := ref.(*KApply)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// NewKApply creates a new object.
func (ms *ModelState) NewKApply(label KLabel, arguments ...KReference) KReference {
	return &KApply{Label: label, List: arguments}
}

// NewKApplyConstant creates a new object.
func NewKApplyConstant(label KLabel, arguments ...KReference) KReference {
	return &KApply{Label: label, List: arguments}
}
