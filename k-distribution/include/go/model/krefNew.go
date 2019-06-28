%COMMENT%

package %PACKAGE_MODEL%

// NewKVariable creates a new object and returns the reference.
func (ms *ModelState) NewKVariable(name string) KReference {
	return ms.addObject(&KVariable{Name: name})
}

// NewInjectedKLabel creates a new InjectedKLabel object and returns the reference.
func (ms *ModelState) NewInjectedKLabel(label KLabel) KReference {
	return ms.addObject(&InjectedKLabel{Label: label})
}

// NewKToken creates a new object and returns the reference.
func (ms *ModelState) NewKToken(sort Sort, value string) KReference {
	return ms.addObject(&KToken{Sort: sort, Value: value})
}

// NewKTokenConstant creates a new KToken constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewKTokenConstant(sort Sort, value string) KReference {
	ref := constantsModel.NewKToken(sort, value)
	ref.constantObject = true
	return ref
}

// NewBytes creates a new K string object from a Go string
func (ms *ModelState) NewBytes(value []byte) KReference {
	return ms.addObject(&Bytes{Value: value})
}

// NewList creates a new object and returns the reference.
func (ms *ModelState) NewList(sort Sort, label KLabel, value []KReference) KReference {
	return ms.addObject(&List{Sort: sort, Label: label, Data: value})
}

// NewMap creates a new object and returns the reference.
func (ms *ModelState) NewMap(sort Sort, label KLabel, value map[KMapKey]KReference) KReference {
	return ms.addObject(&Map{Sort: sort, Label: label, Data: value})
}

// NewSet creates a new object and returns the reference.
func (ms *ModelState) NewSet(sort Sort, label KLabel, value map[KMapKey]bool) KReference {
	return ms.addObject(&Set{Sort: sort, Label: label, Data: value})
}

// NewArray creates a new object and returns the reference.
func (ms *ModelState) NewArray(sort Sort, value *DynamicArray) KReference {
	return ms.addObject(&Array{Sort: sort, Data: value})
}
