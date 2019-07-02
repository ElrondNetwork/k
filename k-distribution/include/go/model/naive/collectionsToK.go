%COMMENT%

package %PACKAGE_MODEL%

import (
	"sort"
)

// CollectionsToK converts all collections to standard K items, like KApply, KToken, etc.
func (ms *ModelState) CollectionsToK(k KReference) KReference {
	return k.collectionsToK(ms)
}

func (k *Map) collectionsToK(ms *ModelState) K {
	if len(k.Data) == 0 {
		return &KApply{Label: UnitFor(k.Label), List: nil}
	}

	// sort entries
	orderedKVPairs := ms.MapOrderedKeyValuePairs(k)

	// process
	elemLabel := ElementFor(k.Label)
	var result K
	for i, pair := range orderedKVPairs {
		elemK := &KApply{Label: elemLabel, List: []KReference{pair.Key, pair.Value.collectionsToK(ms)}}
		if i == 0 {
			result = elemK
		} else {
			newResult := &KApply{Label: k.Label, List: []K{result, elemK}}
			result = newResult
		}
	}

	return result
}

func (k *List) collectionsToK(ms *ModelState) K {
	if len(k.Data) == 0 {
		return &KApply{Label: UnitFor(k.Label), List: nil}
	}

	elemLabel := ElementFor(k.Label)
	var result K
	for i, elem := range k.Data {
		elemK := &KApply{Label: elemLabel, List: []K{elem}}
		if i == 0 {
			result = elemK
		} else {
			newResult := &KApply{Label: k.Label, List: []K{result, elemK}}
			result = newResult
		}
	}

	return result
}

func (k *Set) collectionsToK(ms *ModelState) K {
	if len(k.Data) == 0 {
		return &KApply{Label: UnitFor(k.Label), List: nil}
	}

	// sort keys
	var keysAsString []string
	keyStrToK := make(map[string]K)
	for key := range k.Data {
		keyAsStr := key.String()
		keysAsString = append(keysAsString, keyAsStr)
		kkey, err := ms.ToKItem(key)
		if err != nil {
			panic(err)
		}
		keyStrToK[keyAsStr] = kkey
	}
	sort.Strings(keysAsString)

	// process
	elemLabel := ElementFor(k.Label)
	var result K
	for i, keyAsStr := range keysAsString {
		key := keyStrToK[keyAsStr]
		elemK := &KApply{Label: elemLabel, List: []K{key}}
		if i == 0 {
			result = elemK
		} else {
			newResult := &KApply{Label: k.Label, List: []K{result, elemK}}
			result = newResult
		}
	}

	return result
}

func (k *Array) collectionsToK(ms *ModelState) K {
	result := &KApply{Label: UnitForArray(k.Sort), List: []K{
		ms.NewString("uid"),
		ms.FromUint64(k.Data.MaxSize),
	}}

	slice := k.Data.ToSlice()
	for i, item := range slice {
		newResult := &KApply{Label: ElementForArray(k.Sort), List: []K{
			result,
			ms.FromInt(i),
			item,
		}}
		result = newResult
	}

	return result

}

func (k *KApply) collectionsToK(ms *ModelState) K {
	newList := make([]K, len(k.List))
	for i, child := range k.List {
		newList[i] = child.collectionsToK(ms)
	}
	return &KApply{Label: k.Label, List: newList}
}

func (k *KSequence) collectionsToK(ms *ModelState) K {
	newKs := make([]K, len(k.Ks))
	for i, child := range k.Ks {
		newKs[i] = child.collectionsToK(ms)
	}
	return &KSequence{Ks: newKs}
}

func (k *InjectedKLabel) collectionsToK(ms *ModelState) K {
	return k
}

func (k *KToken) collectionsToK(ms *ModelState) K {
	return k
}

func (k *KVariable) collectionsToK(ms *ModelState) K {
	return k
}

func (k *bigInt) collectionsToK(ms *ModelState) K {
	return k
}

func (k *MInt) collectionsToK(ms *ModelState) K {
	return k
}

func (k *Float) collectionsToK(ms *ModelState) K {
	return k
}

func (k *String) collectionsToK(ms *ModelState) K {
	return k
}

func (k *StringBuffer) collectionsToK(ms *ModelState) K {
	return k
}

func (k *Bytes) collectionsToK(ms *ModelState) K {
	return k
}

func (k *Bool) collectionsToK(ms *ModelState) K {
	return k
}

func (k *Bottom) collectionsToK(ms *ModelState) K {
	return k
}
