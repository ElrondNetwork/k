%COMMENT%

package %PACKAGE_MODEL%

// RecycleUnused sends to the recycle bin all objects left without references.
// This goes recursively through the whole sub-tree.
func (ms *ModelState) RecycleUnused(ref KReference) {
	if ref.constantObject {
		return
	}

	switch ref.refType {
	case boolRef:
	case bottomRef:
	case emptyKseqRef:
	case smallPositiveIntRef:
	case smallNegativeIntRef:
	case bigIntRef:
		obj, _ := ms.getBigIntObject(ref)
		if obj.reuseStatus == active && obj.referenceCount < 1 {
            // recycle
            obj.referenceCount = 0
            obj.reuseStatus = inRecycleBin
            ms.bigIntRecycleBin = append(ms.bigIntRecycleBin, ref)
		}
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		for _, child := range ks {
			ms.RecycleUnused(child)
		}
	default:
		// object types
		obj := ms.getReferencedObject(ref)
		obj.recycleUnused(ms)
	}
}

func (k *KApply) recycleUnused(ms *ModelState) {
	for _, child := range k.List {
		ms.RecycleUnused(child)
	}
}

func (k *InjectedKLabel) recycleUnused(ms *ModelState) {
}

func (k *KToken) recycleUnused(ms *ModelState) {
}

func (k *KVariable) recycleUnused(ms *ModelState) {
}

func (k *Map) recycleUnused(ms *ModelState) {
	for _, v := range k.Data {
		ms.RecycleUnused(v)
	}
}

func (k *List) recycleUnused(ms *ModelState) {
	for _, item := range k.Data {
		ms.RecycleUnused(item)
	}
}

func (k *Set) recycleUnused(ms *ModelState) {
}

func (k *Array) recycleUnused(ms *ModelState) {
	for i := 0; i < len(k.Data.data); i++ {
		if k.Data.data[i] != NullReference {
			ms.RecycleUnused(k.Data.data[i])
		}
	}
}

func (k *MInt) recycleUnused(ms *ModelState) {
}

func (k *Float) recycleUnused(ms *ModelState) {
}

func (k *String) recycleUnused(ms *ModelState) {
}

func (k *StringBuffer) recycleUnused(ms *ModelState) {
}

func (k *Bytes) recycleUnused(ms *ModelState) {
}
