%COMMENT%

package %PACKAGE_MODEL%

// KApply is a KObject representing a KApply item in K
type KApply struct {
	Label KLabel
	List  []KReference
}

type KApplyReferenceAndObject struct {
    obj *KApply
}

// GetReference retrieves the references, but also increments the reference count.
func (ro KApplyReferenceAndObject) GetReference() KReference {
	return ro.obj
}

// Arg ...
func (ro KApplyReferenceAndObject) Arg(i int) KReference {
	return ro.obj.List[i]
}

// Label ...
func (ro KApplyReferenceAndObject) Label() KLabel {
	return ro.obj.Label
}

// Arity ...
func (ro KApplyReferenceAndObject) Arity() int {
	return len(ro.obj.List)
}

// ReleaseKApply matches a KApply, then signals that the reference is no longer needed.
func (ms *ModelState) ReleaseKApply(ref KReference) (KApplyReferenceAndObject, bool) {
	ro, typeOk := ms.MatchKApply(ref)
	if !typeOk {
		return ro, false
	}
	//ro.obj.nrReferences--
	return ro, true
}

// MatchKApply matches a KApply, but does not release the reference.
func (ms *ModelState) MatchKApply(ref KReference) (KApplyReferenceAndObject, bool) {
	obj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		return KApplyReferenceAndObject{obj: nil}, false
	}
	return KApplyReferenceAndObject{obj: obj}, true
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

// KApply0Ref yields a reference to a KApply with 0 arguments.
func (ms *ModelState) KApply0Ref(label KLabel) KReference {
	return &KApply{Label: label, List: nil}
}

// NewKApply creates a new object.
func (ms *ModelState) NewKApply(label KLabel, arguments ...KReference) KReference {
	return &KApply{Label: label, List: arguments}
}

// ReuseKApply takes an existing KApply object that is no longer needed and replaces its arguments.
func (ms *ModelState) ReuseKApply(rreuseRo KApplyReferenceAndObject, label KLabel, arguments ...KReference) KReference {
	/*reuseKapp, typeOk := reuseObj.(*KApply)
	if !typeOk {
		panic("wrong object type for reference")
	}
	if reuseKapp.Label != label {
	    panic("wrong reuse object provided to ReuseKApply, Label mismatch")
	}
	if len(arguments) != len(reuseKapp.List) {
	     panic("wrong reuse object provided to ReuseKApply, arity mismatch")
	}
	for i, arg := range arguments {
	    reuseKapp.List[i] = arg
	}
	return reuseObj*/
	return ms.NewKApply(label, arguments...)
}

// NewKApplyConstant creates a new object.
func NewKApplyConstant(label KLabel, arguments ...KReference) KReference {
	return &KApply{Label: label, List: arguments}
}
