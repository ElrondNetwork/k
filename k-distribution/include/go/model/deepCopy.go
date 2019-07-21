%COMMENT%

package %PACKAGE%

// DeepCopy copies the structures underlying a reference from one model to the other.
// Arguments from and to can be the same model.
// Flag mainModelOnly, if set to true, will not deep copy constants, only references found in the main model.
func DeepCopy(from, to *ModelState, ref KReference, mainModelOnly bool) KReference {
	refType, constant, value := parseKrefBasic(ref)
	if mainModelOnly && constant {
		return ref
	}

	// collection types
	if isCollectionType(refType) {
		_, sortInt, labelInt, index := parseKrefCollection(ref)
		obj := from.getReferencedObject(index, false)
		copiedObj := obj.deepCopy(from, to, mainModelOnly)
		return to.addCollectionObject(Sort(sortInt), KLabel(labelInt), copiedObj)
	}

	switch refType {
	case boolRef:
		return ref
	case bottomRef:
		return ref
	case emptyKseqRef:
		return ref
	case nonEmptyKseqRef:
		ks := from.KSequenceToSlice(ref)
		newKs := make([]KReference, len(ks))
		for i, child := range ks {
			newKs[i] = DeepCopy(from, to, child, mainModelOnly)
		}
		return to.NewKSequence(newKs)
	case smallPositiveIntRef:
		return ref
	case smallNegativeIntRef:
		return ref
	case bigIntRef:
		obj, _ := from.getBigIntObject(ref)
		newRef, newObj := to.newBigIntObject()
		newObj.bigValue.Set(obj.bigValue)
		return newRef
	case kapplyRef:
		argSlice := from.kapplyArgSlice(ref)
		argCopy := make([]KReference, len(argSlice))
		for i, child := range argSlice {
			argCopy[i] = DeepCopy(from, to, child, mainModelOnly)
		}
		return to.NewKApply(from.KApplyLabel(ref), argCopy...)
	case stringRef:
		str, _ := from.GetString(ref)
		return to.NewString(str)
	case bytesRef:
		bytes, _ := from.GetBytes(ref)
		return to.NewBytes(bytes)
	case ktokenRef:
		ktoken, _ := from.GetKTokenObject(ref)
		return to.NewKToken(ktoken.Sort, ktoken.Value)
	default:
		// object types
		obj := from.getReferencedObject(value, constant)
		copiedObj := obj.deepCopy(from, to, mainModelOnly)
		if copiedObj == obj {
			// if no new instance was created,
			// it means that the object does not need to be deep copied
			return ref
		}
		return to.addObject(copiedObj)
	}
}

func (k *InjectedKLabel) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return &InjectedKLabel{Label: k.Label}
}

func (k *KVariable) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return &KVariable{Name: k.Name}
}

func (k *Map) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	mapCopy := make(map[KMapKey]KReference)
	for key, val := range k.Data {
		mapCopy[key] = DeepCopy(from, to, val, mainModelOnly)
	}
	return &Map{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  mapCopy,
	}
}

func (k *List) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	listCopy := make([]KReference, len(k.Data))
	for i, elem := range k.Data {
		listCopy[i] = DeepCopy(from, to, elem, mainModelOnly)
	}
	return &List{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  listCopy,
	}
}

func (k *Set) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	mapCopy := make(map[KMapKey]bool)
	for key := range k.Data {
		mapCopy[key] = true
	}
	return &Set{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  mapCopy,
	}
}

func (k *Array) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return k // TODO: not implemented
}

func (k *MInt) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return k // not implemented
}

func (k *Float) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return k // not implemented
}

func (k *StringBuffer) deepCopy(from, to *ModelState, mainModelOnly bool) KObject {
	return k // no deep copy needed here
}
