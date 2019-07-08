%COMMENT%

package %PACKAGE%

// Map is a KObject representing a map in K
type Map struct {
	Sort  Sort
	Label KLabel
	Data  map[KMapKey]KReference
}

// Set is a KObject representing a set in K
type Set struct {
	Sort  Sort
	Label KLabel
	Data  map[KMapKey]bool
}

// List is a KObject representing a list in K
type List struct {
	Sort  Sort
	Label KLabel
	Data  []KReference
}

// Array is a KObject holding an array that can grow
type Array struct {
	Sort Sort
	Data *DynamicArray
}

// IsMap returns true if reference points to a map with given sort
func (ms *ModelState) IsMap(ref KReference, expectedSort Sort) bool {
	obj, typeOk := ms.GetMapObject(ref)
	if !typeOk {
		return false
	}
	return obj.Sort == expectedSort
}

// IsSet returns true if reference points to a set with given sort
func (ms *ModelState) IsSet(ref KReference, expectedSort Sort) bool {
	obj, typeOk := ms.GetSetObject(ref)
	if !typeOk {
		return false
	}
	return obj.Sort == expectedSort
}

// IsList returns true if reference points to a list with given sort
func (ms *ModelState) IsList(ref KReference, expectedSort Sort) bool {
	obj, typeOk := ms.GetListObject(ref)
	if !typeOk {
		return false
	}
	return obj.Sort == expectedSort
}

// IsArray returns true if reference points to an array with given sort
func (ms *ModelState) IsArray(ref KReference, expectedSort Sort) bool {
	obj, typeOk := ms.GetArrayObject(ref)
	if !typeOk {
		return false
	}
	return obj.Sort == expectedSort
}

// GetListObject yields the cast object for a List reference, if possible.
func (ms *ModelState) GetListObject(ref KReference) (*List, bool) {
	castObj, typeOk := ref.(*List)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// IsEmptyList returns true only if argument references an empty list, with given sort and label.
func (ms *ModelState) IsEmptyList(ref KReference, expectedSort Sort, expectedLabel KLabel) bool {
	castObj, typeOk := ms.GetListObject(ref)
	if !typeOk {
		return false
	}
	if castObj.Sort != expectedSort {
		return false
	}
	if castObj.Label != expectedLabel {
		return false
	}
	return len(castObj.Data) == 0
}

// ListSplitHeadTail returns true only if argument references an empty list.
// Returns nothing if it is not a list, it is empty, or if sort or label do not match.
func (ms *ModelState) ListSplitHeadTail(ref KReference, expectedSort Sort, expectedLabel KLabel) (ok bool, head KReference, tail KReference) {
	castObj, typeOk := ms.GetListObject(ref)
	if !typeOk {
		return false, NullReference, NullReference
	}
	if castObj.Sort != expectedSort {
		return false, NullReference, NullReference
	}
	if castObj.Label != expectedLabel {
		return !ok, NullReference, NullReference
	}
	if len(castObj.Data) == 0 {
		return false, NullReference, NullReference
	}
	tailRef := &List{Sort: castObj.Sort, Label: castObj.Label, Data: castObj.Data[1:]}
	return true, castObj.Data[0], tailRef
}

// GetMapObject yields the cast object for a Map reference, if possible.
func (ms *ModelState) GetMapObject(k KReference) (*Map, bool) {
	castObj, typeOk := k.(*Map)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// GetSetObject yields the cast object for a Set reference, if possible.
func (ms *ModelState) GetSetObject(k KReference) (*Set, bool) {
	castObj, typeOk := k.(*Set)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// GetArrayObject yields the cast object for an Array reference, if possible.
func (ms *ModelState) GetArrayObject(k KReference) (*Array, bool) {
	castObj, typeOk := k.(*Array)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// NewList creates a new object and returns the reference.
func (ms *ModelState) NewList(sort Sort, label KLabel, value []KReference) KReference {
	return &List{Sort: sort, Label: label, Data: value}
}

// NewMap creates a new object and returns the reference.
func (ms *ModelState) NewMap(sort Sort, label KLabel, value map[KMapKey]KReference) KReference {
	return &Map{Sort: sort, Label: label, Data: value}
}

// NewSet creates a new object and returns the reference.
func (ms *ModelState) NewSet(sort Sort, label KLabel, value map[KMapKey]bool) KReference {
	return &Set{Sort: sort, Label: label, Data: value}
}

// NewArray creates a new object and returns the reference.
func (ms *ModelState) NewArray(sort Sort, value *DynamicArray) KReference {
	return &Array{Sort: sort, Data: value}
}
