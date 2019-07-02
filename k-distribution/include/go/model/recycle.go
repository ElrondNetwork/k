%COMMENT%

package %PACKAGE_MODEL%

// MarkInUse sets flags indicating the last step at which objects were used.
// It goes recursively through the whole tree.
func (ms *ModelState) MarkInUse(ref KReference, stepNr int) {
	switch ref.refType {
	case boolRef:
	case bottomRef:
	case emptyKseqRef:
	case smallPositiveIntRef:
	case smallNegativeIntRef:
	case bigIntRef:
		obj, _ := ms.getBigIntObject(ref)
		obj.lastInUse = stepNr
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		for _, child := range ks {
			ms.MarkInUse(child, stepNr)
		}
	default:
		// object types
		obj := ms.getReferencedObject(ref)
		obj.markInUse(ms, stepNr)
	}
}

func (k *KApply) markInUse(ms *ModelState, stepNr int) {
	for _, child := range k.List {
		ms.MarkInUse(child, stepNr)
	}
}

func (k *InjectedKLabel) markInUse(ms *ModelState, stepNr int) {
}

func (k *KToken) markInUse(ms *ModelState, stepNr int) {
}

func (k *KVariable) markInUse(ms *ModelState, stepNr int) {
}

func (k *Map) markInUse(ms *ModelState, stepNr int) {
	for _, v := range k.Data {
		ms.MarkInUse(v, stepNr)
	}
}

func (k *List) markInUse(ms *ModelState, stepNr int) {
	for _, item := range k.Data {
		ms.MarkInUse(item, stepNr)
	}
}

func (k *Set) markInUse(ms *ModelState, stepNr int) {
}

func (k *Array) markInUse(ms *ModelState, stepNr int) {
	for i := 0; i < len(k.Data.data); i++ {
		if k.Data.data[i] != NullReference {
			ms.MarkInUse(k.Data.data[i], stepNr)
		}
	}
}

func (k *MInt) markInUse(ms *ModelState, stepNr int) {
}

func (k *Float) markInUse(ms *ModelState, stepNr int) {
}

func (k *String) markInUse(ms *ModelState, stepNr int) {
}

func (k *StringBuffer) markInUse(ms *ModelState, stepNr int) {
}

func (k *Bytes) markInUse(ms *ModelState, stepNr int) {
}

// RecycleUnused will send to recycle bin objects that were not found to be used at given step.
// This goes recursively through the whole tree.
// Function MarkInUse must be called beforehand, for this to work.
func (ms *ModelState) RecycleUnused(ref KReference, stepNr int) {
	switch ref.refType {
	case boolRef:
	case bottomRef:
	case emptyKseqRef:
	case smallPositiveIntRef:
	case smallNegativeIntRef:
	case bigIntRef:
		obj, _ := ms.getBigIntObject(ref)
		if obj.lastInUse < stepNr {
			ms.bigIntRecycleBin = append(ms.bigIntRecycleBin, ref)
		}
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		for _, child := range ks {
			ms.RecycleUnused(child, stepNr)
		}
	default:
		// object types
		obj := ms.getReferencedObject(ref)
		obj.recycleUnused(ms, stepNr)
	}
}

func (k *KApply) recycleUnused(ms *ModelState, stepNr int) {
	for _, child := range k.List {
		ms.RecycleUnused(child, stepNr)
	}
}

func (k *InjectedKLabel) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *KToken) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *KVariable) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *Map) recycleUnused(ms *ModelState, stepNr int) {
	for _, v := range k.Data {
		ms.RecycleUnused(v, stepNr)
	}
}

func (k *List) recycleUnused(ms *ModelState, stepNr int) {
	for _, item := range k.Data {
		ms.RecycleUnused(item, stepNr)
	}
}

func (k *Set) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *Array) recycleUnused(ms *ModelState, stepNr int) {
	for i := 0; i < len(k.Data.data); i++ {
		if k.Data.data[i] != NullReference {
			ms.RecycleUnused(k.Data.data[i], stepNr)
		}
	}
}

func (k *MInt) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *Float) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *String) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *StringBuffer) recycleUnused(ms *ModelState, stepNr int) {
}

func (k *Bytes) recycleUnused(ms *ModelState, stepNr int) {
}
