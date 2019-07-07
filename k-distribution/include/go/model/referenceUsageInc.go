%COMMENT%

package %PACKAGE_MODEL%

// IncreaseUsage increments all reference counters in tree below given root.
// It goes recursively through the whole sub-tree.
func (ms *ModelState) IncreaseUsage(ref KReference) {
	if ref.constantObject {
		return
	}

	switch ref.refType {
	case boolRef:
	case bottomRef:
	case emptyKseqRef:
	case smallPositiveIntRef:
	case smallNegativeIntRef:
	case stringRef:
	case bytesRef:
	case ktokenRef:
	case bigIntRef:
		obj, _ := ms.getBigIntObject(ref)
		if obj.reuseStatus == active {
			if obj.referenceCount < 1 {
				obj.referenceCount = 1
			} else {
				obj.referenceCount++
			}
		}
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		for _, child := range ks {
			ms.IncreaseUsage(child)
		}
	case kapplyRef:
		for _, child := range ms.kapplyArgSlice(ref) {
			ms.IncreaseUsage(child)
		}
	default:
		// object types
		obj := ms.getReferencedObject(ref)
		obj.increaseUsage(ms)
	}
}

func (k *InjectedKLabel) increaseUsage(ms *ModelState) {
}

func (k *KVariable) increaseUsage(ms *ModelState) {
}

func (k *Map) increaseUsage(ms *ModelState) {
	for _, v := range k.Data {
		ms.IncreaseUsage(v)
	}
}

func (k *List) increaseUsage(ms *ModelState) {
	for _, item := range k.Data {
		ms.IncreaseUsage(item)
	}
}

func (k *Set) increaseUsage(ms *ModelState) {
}

func (k *Array) increaseUsage(ms *ModelState) {
	for i := 0; i < len(k.Data.data); i++ {
		if k.Data.data[i] != NullReference {
			ms.IncreaseUsage(k.Data.data[i])
		}
	}
}

func (k *MInt) increaseUsage(ms *ModelState) {
}

func (k *Float) increaseUsage(ms *ModelState) {
}

func (k *StringBuffer) increaseUsage(ms *ModelState) {
}
