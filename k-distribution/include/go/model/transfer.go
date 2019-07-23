%COMMENT%

package %PACKAGE%

// transfer copies the structures underlying a reference from one data container to another.
// It is similar to deep copy.
// Will only transfer references that actually point to the "from" model.
func transfer(from, to *ModelData, ref KReference) KReference {
	refType, dataRef, value := parseKrefBasic(ref)
	if dataRef != from.selfRef {
		return ref
	}

	// collection types
	if isCollectionType(refType) {
		_, _, sortInt, labelInt, index := parseKrefCollection(ref)
		obj := from.getReferencedObject(index)
		copiedObj := obj.transfer(from, to)
		return to.addCollectionObject(Sort(sortInt), KLabel(labelInt), copiedObj)
	}

	switch refType {
	case nullRef:
		return ref
	case boolRef:
		return ref
	case bottomRef:
		return ref
	case emptyKseqRef:
		return ref
	case nonEmptyKseqRef:
		ks := from.ksequenceToSlice(ref)
		newKs := make([]KReference, len(ks))
		for i, child := range ks {
			newKs[i] = transfer(from, to, child)
		}
		return to.assembleKSequence(newKs...)
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
		_, _, label, arity, index := parseKrefKApply(ref)
		if arity == 0 {
			return to.newKApply(label)
		}
		argSlice := from.allKApplyArgs[index : index+arity]
		argCopy := make([]KReference, len(argSlice))
		for i, child := range argSlice {
			argCopy[i] = transfer(from, to, child)
		}
		return to.newKApply(label, argCopy...)
	case stringRef:
		_, _, startIndex, length := parseKrefBytes(ref)
		if length == 0 {
			return ref
		}
		return to.newBytes(stringRef, from.allBytes[startIndex:startIndex+length])
	case bytesRef:
		_, _, startIndex, length := parseKrefBytes(ref)
		if length == 0 {
			return ref
		}
		return to.newBytes(bytesRef, from.allBytes[startIndex:startIndex+length])
	case ktokenRef:
		_, _, sort, length, index := parseKrefKToken(ref)
		return to.newKToken(sort, from.allBytes[index:index+length])
	default:
		// object types
		obj := from.getReferencedObject(value)
		copiedObj := obj.transfer(from, to)
		return to.addObject(copiedObj)
	}
}

func (k *InjectedKLabel) transfer(from, to *ModelData) KObject {
	return k
}

func (k *KVariable) transfer(from, to *ModelData) KObject {
	return k
}

func (k *Map) transfer(from, to *ModelData) KObject {
	mapCopy := make(map[KMapKey]KReference)
	for key, val := range k.Data {
		mapCopy[key] = transfer(from, to, val)
	}
	return &Map{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  mapCopy,
	}
}

func (k *List) transfer(from, to *ModelData) KObject {
	listCopy := make([]KReference, len(k.Data))
	for i, elem := range k.Data {
		listCopy[i] = transfer(from, to, elem)
	}
	return &List{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  listCopy,
	}
}

func (k *Set) transfer(from, to *ModelData) KObject {
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

func (k *Array) transfer(from, to *ModelData) KObject {
	dataCopy := make([]KReference, len(k.Data.data))
	for i, elem := range k.Data.data {
		dataCopy[i] = transfer(from, to, elem)
	}

	return &Array{
		Sort: k.Sort,
		Data: &DynamicArray{
			MaxSize: k.Data.MaxSize,
			data:    dataCopy,
			Default: transfer(from, to, k.Data.Default),
			ms:      k.Data.ms,
		},
	}
}

func (k *MInt) transfer(from, to *ModelData) KObject {
	return k
}

func (k *Float) transfer(from, to *ModelData) KObject {
	return k
}

func (k *StringBuffer) transfer(from, to *ModelData) KObject {
	return k
}
