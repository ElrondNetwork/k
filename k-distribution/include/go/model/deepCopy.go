%COMMENT%

package %PACKAGE_MODEL%

import (
	"math/big"
)

// DeepCopy yields a fresh copy of the K item given as argument.
func (ms *ModelState) DeepCopy(ref KReference) KReference {
	switch ref.refType {
	case boolRef:
		return ref
	case bottomRef:
		return ref
	case emptyKseqRef:
		return ref
	case nonEmptyKseqRef:
		ks := ms.KSequenceToSlice(ref)
		newKs := make([]KReference, len(ks))
		for i, child := range ks {
			newKs[i] = ms.DeepCopy(child)
		}
		return ms.NewKSequence(newKs)
	case smallPositiveIntRef:
		return ref
	case smallNegativeIntRef:
		return ref
	case bigIntRef:
		obj, _ := ms.getBigIntObject(ref)
		intCopy := big.NewInt(0)
		intCopy.Set(obj.bigValue)
		return ms.addBigIntObject(intCopy)
	default:
		// object types
		obj := ms.getReferencedObject(ref)
		copiedObj := obj.deepCopy(ms)
		if copiedObj == obj {
			// if no new instance was created,
			// it means that the object does not need to be deep copied
			return ref
		}
		return ms.addObject(obj)
	}
}

func (k *KApply) deepCopy(ms *ModelState) KObject {
	listCopy := make([]KReference, len(k.List))
	for i, child := range k.List {
		listCopy[i] = ms.DeepCopy(child)
	}
	return &KApply{Label: k.Label, List: listCopy}
}

func (k *InjectedKLabel) deepCopy(ms *ModelState) KObject {
	return &InjectedKLabel{Label: k.Label}
}

func (k *KToken) deepCopy(ms *ModelState) KObject {
	return &KToken{Sort: k.Sort, Value: k.Value}
}

func (k *KVariable) deepCopy(ms *ModelState) KObject {
	return &KVariable{Name: k.Name}
}

func (k *Map) deepCopy(ms *ModelState) KObject {
	mapCopy := make(map[KMapKey]KReference)
	for key, val := range k.Data {
		mapCopy[key] = ms.DeepCopy(val)
	}
	return &Map{Data: mapCopy}
}

func (k *List) deepCopy(ms *ModelState) KObject {
	listCopy := make([]KReference, len(k.Data))
	for i, elem := range k.Data {
		listCopy[i] = ms.DeepCopy(elem)
	}
	return &List{
		Sort:  k.Sort,
		Label: k.Label,
		Data:  listCopy,
	}
}

func (k *Set) deepCopy(ms *ModelState) KObject {
	mapCopy := make(map[KMapKey]bool)
	for key := range k.Data {
		mapCopy[key] = true
	}
	return &Set{Data: mapCopy}
}

func (k *Array) deepCopy(ms *ModelState) KObject {
	return k // TODO: not implemented
}

func (k *MInt) deepCopy(ms *ModelState) KObject {
	return k // not implemented
}

func (k *Float) deepCopy(ms *ModelState) KObject {
	return k // not implemented
}

func (k *String) deepCopy(ms *ModelState) KObject {
	return &String{Value: k.Value}
}

func (k *StringBuffer) deepCopy(ms *ModelState) KObject {
	return k // no deep copy needed here
}

func (k *Bytes) deepCopy(ms *ModelState) KObject {
	bytesCopy := make([]byte, len(k.Value))
	copy(bytesCopy, k.Value)
	return &Bytes{Value: bytesCopy}
}
