%COMMENT%

package %PACKAGE_MODEL%

// KApply is a KObject representing a KApply item in K
type KApply struct {
	Label        KLabel
	List         []KReference
	nrReferences int
}

func (*KApply) referenceType() kreferenceType {
	return kapplyRef
}

// KApplyReferenceAndObject is used in the rules to hold both the reference and the object itself.
// We keep it as a value type, for the duration of the rule match. Should be short-lived.
type KApplyReferenceAndObject struct {
	ref KReference
	obj *KApply
}

// GetReference retrieves the references, but also increments the reference count.
func (ro KApplyReferenceAndObject) GetReference() KReference {
	ro.obj.nrReferences++
	return ro.ref
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

// ReleaseKApply signals that the state that the KApply was part of is no longer needed,
// we can decrement nrReferences.
func (ms *ModelState) ReleaseKApply(ro KApplyReferenceAndObject) {
	if ro.ref.constantObject {
	    return
	}
    if ro.obj.nrReferences > 0 {
        ro.obj.nrReferences--
    }
}

// MatchKApply matches a KApply, but does not release the reference.
func (ms *ModelState) MatchKApply(ref KReference) (KApplyReferenceAndObject, bool) {
	obj, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		return KApplyReferenceAndObject{ref: NullReference, obj: nil}, false
	}
	return KApplyReferenceAndObject{ref: ref, obj: obj}, true
}

// CastKApply returns true if argument is a KApply item.
// Also returns argument, for convenience.
func (ms *ModelState) CastKApply(ref KReference) (KReference, bool) {
	_, typeOk := ms.GetKApplyObject(ref)
	if !typeOk {
		return NullReference, false
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
	if ref.refType != kapplyRef {
		return nil, false
	}
	obj := ms.getReferencedObject(ref)
	castObj, typeOk := obj.(*KApply)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// KApply0Ref yields a reference to a KApply with 0 arguments.
func (ms *ModelState) KApply0Ref(label KLabel) KReference {
	return ms.addObject(&KApply{Label: label, List: nil, nrReferences: 1})
}

// NewKApply creates a new object and returns the reference.
func (ms *ModelState) NewKApply(label KLabel, arguments ...KReference) KReference {
	return ms.addObject(&KApply{Label: label, List: arguments, nrReferences: 1})
}

// ReuseKApply takes an existing KApply object that is no longer needed and replaces its arguments.
func (ms *ModelState) ReuseKApply(reuseRo KApplyReferenceAndObject, label KLabel, arguments ...KReference) KReference {
	if reuseRo.ref.refType != kapplyRef {
		panic("wrong reuse object type provided to ReuseKApply")
	}
	if reuseRo.obj.Label != label {
		panic("wrong reuse object provided to ReuseKApply, Label mismatch")
	}
	if len(arguments) != len(reuseRo.obj.List) {
		panic("wrong reuse object provided to ReuseKApply, arity mismatch")
	}
	if reuseRo.ref.constantObject || reuseRo.obj.nrReferences > 0 {
		// cannot reuse objects that are still in use elsewhere
		return ms.NewKApply(label, arguments...)
	}

	// the reuse itself
	reuseRo.obj.nrReferences = 1
	for i, arg := range arguments {
		reuseRo.obj.List[i] = arg
	}

	return reuseRo.ref

}

// NewKApplyConstant creates a new integer constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewKApplyConstant(label KLabel, arguments ...KReference) KReference {
	ref := constantsModel.NewKApply(label, arguments...)
	ref.constantObject = true
	return ref
}
