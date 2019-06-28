%COMMENT%

package %PACKAGE_MODEL%

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
	if ref.refType != listRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*List)
	if !typeOk {
		panic("wrong object type for reference")
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
    tailRef := ms.addObject(&List{Sort: castObj.Sort, Label: castObj.Label, Data: castObj.Data[1:]})
	return true, castObj.Data[0], tailRef
}

// GetMapObject yields the cast object for a Map reference, if possible.
func (ms *ModelState) GetMapObject(ref KReference) (*Map, bool) {
	if ref.refType != mapRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*Map)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// GetSetObject yields the cast object for a Set reference, if possible.
func (ms *ModelState) GetSetObject(ref KReference) (*Set, bool) {
	if ref.refType != setRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*Set)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// GetArrayObject yields the cast object for an Array reference, if possible.
func (ms *ModelState) GetArrayObject(ref KReference) (*Array, bool) {
	if ref.refType != arrayRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*Array)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}