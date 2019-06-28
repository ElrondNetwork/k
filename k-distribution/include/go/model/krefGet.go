%COMMENT%

package %PACKAGE_MODEL%

// GetKTokenObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetKTokenObject(ref KReference) (*KToken, bool) {
	if ref.refType != ktokenRef {
		return nil, false
	}
	ms.getObject(ref)
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*KToken)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// GetFloatObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetFloatObject(ref KReference) (*Float, bool) {
	if ref.refType != floatRef {
		return nil, false
	}
	ms.getObject(ref)
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*Float)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// IsFloat returns true if reference points to a float
func IsFloat(ref KReference) bool {
	return ref.refType == floatRef
}

// IsMInt returns true if reference points to a string buffer
func IsMInt(ref KReference) bool {
	return ref.refType == mintRef
}

// IsBottom returns true if reference points to bottom
func IsBottom(ref KReference) bool {
	return ref.refType == bottomRef
}


