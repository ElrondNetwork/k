%COMMENT%

package %PACKAGE_MODEL%

// CollectionsToK converts all collections to standard K items, like KApply, KToken, etc.
func (ms *ModelState) CollectionsToK(ref KReference) KReference {
	if ref.refType == nonEmptyKseqRef {
		ks := ms.KSequenceToSlice(ref)
		newKs := make([]KReference, len(ks))
		for i, child := range ks {
			obj := ms.getObject(child)
			newKs[i] = obj.collectionsToK(ms)
		}
		return ms.NewKSequence(newKs)
	} else if ref.refType == mapRef ||
		ref.refType == listRef ||
		ref.refType == setRef ||
		ref.refType == arrayRef ||
		ref.refType == kapplyRef {

		// object types
		obj := ms.getObject(ref)
		return obj.collectionsToK(ms)
	}

	// no processing required for the others
	return ref
}

func (k *Map) collectionsToK(ms *ModelState) KReference {
	if len(k.Data) == 0 {
		return ms.NewKApply(UnitFor(k.Label))
	}

	// sort entries
	orderedKVPairs := ms.MapOrderedKeyValuePairs(k)

	// process
	elemLabel := ElementFor(k.Label)
	var result KReference
	for i, pair := range orderedKVPairs {
		elemK := ms.NewKApply(elemLabel, pair.Key, ms.CollectionsToK(pair.Value))
		if i == 0 {
			result = elemK
		} else {
			newResult := ms.NewKApply(k.Label, result, elemK)
			result = newResult
		}
	}

	return result
}

func (k *List) collectionsToK(ms *ModelState) KReference {
	if len(k.Data) == 0 {
		return ms.KApply0Ref(UnitFor(k.Label))
	}

	elemLabel := ElementFor(k.Label)
	var result KReference
	for i, elem := range k.Data {
		elemK := ms.NewKApply(elemLabel, elem)
		if i == 0 {
			result = elemK
		} else {
			newResult := ms.NewKApply(k.Label, result, elemK)
			result = newResult
		}
	}

	return result
}

func (k *Set) collectionsToK(ms *ModelState) KReference {
	if len(k.Data) == 0 {
		return ms.KApply0Ref(UnitFor(k.Label))
	}

	// sort keys
	sortedKeys := ms.SetOrderedElements(k)

	// process
	elemLabel := ElementFor(k.Label)
	var result KReference
	for i, key := range sortedKeys {
		elemK := ms.NewKApply(elemLabel, key)
		if i == 0 {
			result = elemK
		} else {
			newResult := ms.NewKApply(k.Label, result, elemK)
			result = newResult
		}
	}

	return result
}

func (k *Array) collectionsToK(ms *ModelState) KReference {
	result := ms.NewKApply(UnitForArray(k.Sort),
		ms.NewString("uid"),
		ms.FromUint64(k.Data.MaxSize))

	slice := k.Data.ToSlice()
	for i, item := range slice {
		newResult := ms.NewKApply(ElementForArray(k.Sort),
			result,
			ms.FromInt(i),
			item)
		result = newResult
	}

	return result

}

func (k *KApply) collectionsToK(ms *ModelState) KReference {
	newList := make([]KReference, len(k.List))
	for i, child := range k.List {
		newList[i] = ms.CollectionsToK(child)
	}
	return ms.NewKApply(k.Label, newList...)
}

func (k *InjectedKLabel) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *KToken) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *KVariable) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *BigInt) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *MInt) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *Float) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *String) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *StringBuffer) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}

func (k *Bytes) collectionsToK(ms *ModelState) KReference {
	panic("collectionsToK shouldn't be called for this object type")
}
