%COMMENT%

package %PACKAGE_MODEL%

// KReference is a reference to a K item.
// In this implementation, it is just an alias to the regular pointer to the item.
type KReference K

// CastKApply yields the cast object for a KSequence reference, if possible.
func (ms *ModelState) CastKSequence(ref KReference) (KSequence, bool) {
	kapp, isKapp := ref.(KSequence)
	if !isKapp {
		return EmptyKSequence, false
	}
	return kapp, true
}

// CastKApply yields the cast object for a KApply reference, if possible.
func (ms *ModelState) CastKApply(ref KReference) (*KApply, bool) {
	kapp, isKapp := ref.(*KApply)
	if !isKapp {
		return nil, false
	}
	return kapp, true
}

// CastKApplyAndCheck yields the cast object for a KApply reference, if possible,
// and if thr KApply has the right label and arity.
func (ms *ModelState) CastKApplyAndCheck(ref K, lbl KLabel, arity int) (*KApply, bool) {
	kapp, isKapp := ms.CastKApply(ref)
	if !isKapp {
		return nil, false
	}
	if kapp.Label != lbl {
		return nil, false
	}
	if len(kapp.List) != arity {
		return nil, false
	}
	return kapp, true
}

// CastKApply yields the cast object for a InjectedKLabel reference, if possible.
func (ms *ModelState) CastInjectedKLabel(ref KReference) (*InjectedKLabel, bool) {
	cast, typeOk := ref.(*InjectedKLabel)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a KToken reference, if possible.
func (ms *ModelState) CastKToken(ref KReference) (*KToken, bool) {
	cast, typeOk := ref.(*KToken)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a KVariable reference, if possible.
func (ms *ModelState) CastKVariable(ref KReference) (*KVariable, bool) {
	cast, typeOk := ref.(*KVariable)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Map reference, if possible.
func (ms *ModelState) CastMap(ref KReference) (*Map, bool) {
	cast, typeOk := ref.(*Map)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Set reference, if possible.
func (ms *ModelState) CastSet(ref KReference) (*Set, bool) {
	cast, typeOk := ref.(*Set)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a List reference, if possible.
func (ms *ModelState) CastList(ref KReference) (*List, bool) {
	cast, typeOk := ref.(*List)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Array reference, if possible.
func (ms *ModelState) CastArray(ref KReference) (*Array, bool) {
	cast, typeOk := ref.(*Array)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Int reference, if possible.
func (ms *ModelState) CastInt(ref KReference) (*Int, bool) {
	cast, typeOk := ref.(*Int)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a MInt reference, if possible.
func (ms *ModelState) CastMInt(ref KReference) (*MInt, bool) {
	cast, typeOk := ref.(*MInt)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Float reference, if possible.
func (ms *ModelState) CastFloat(ref KReference) (*Float, bool) {
	cast, typeOk := ref.(*Float)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a String reference, if possible.
func (ms *ModelState) CastString(ref KReference) (*String, bool) {
	cast, typeOk := ref.(*String)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a StringBuffer reference, if possible.
func (ms *ModelState) CastStringBuffer(ref KReference) (*StringBuffer, bool) {
	cast, typeOk := ref.(*StringBuffer)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastKApply yields the cast object for a Bytes reference, if possible.
func (ms *ModelState) CastBytes(ref KReference) (*Bytes, bool) {
	cast, typeOk := ref.(*Bytes)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// CastBottom yields the cast object for a Bottom reference, if possible.
func (ms *ModelState) CastBottom(ref KReference) (*Bottom, bool) {
	cast, typeOk := ref.(*Bottom)
	if !typeOk {
		return nil, false
	}
	return cast, true
}

// KApply0Ref yields a reference to a KApply with 0 arguments.
func (ms *ModelState) KApply0Ref(label KLabel) KReference {
	return &KApply{Label: label, List: nil}
}

// NewKApply creates a new object and returns the reference.
func (ms *ModelState) NewKApply(label KLabel, arguments ...KReference) KReference {
	return &KApply{Label: label, List: arguments}
}

// NewInjectedKLabel creates a new InjectedKLabel object and returns the reference.
func (ms *ModelState) NewInjectedKLabel(label KLabel) KReference {
	return &InjectedKLabel{Label: label}
}
