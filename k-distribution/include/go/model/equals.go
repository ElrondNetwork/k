%COMMENT%

package %PACKAGE_MODEL%

import (
	"bytes"
)

// Equals performs a deep comparison, recursively.
func (ms *ModelState) Equals(ref1 KReference, ref2 KReference) bool {
	if ref1 == ref2 {
		// identical references means the same object
		return true
	}
	if ref1.refType != ref2.refType {
		// different types cannot be equal
		return false
	}

	switch ref1.refType {
	case boolRef:
		return false // if they were equal, ref1 == ref2 condition would already have returned true
	case bottomRef:
		panic("there shouldn't be different references of type bottomRef, only one")
	case emptyKseqRef:
		panic("there shouldn't be different references of type emptyKseqRef, only one")
	case nonEmptyKseqRef:
		return ms.ksequenceEquals(ref1, ref2)
	default:
		// object types
		obj1 := ms.getReferencedObject(ref1)
		obj2 := ms.getReferencedObject(ref2)
		return obj1.equals(ms, obj2)
	}
}

func (k *KApply) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*KApply)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Label != other.Label {
		return false
	}
	if len(k.List) != len(other.List) {
		return false
	}
	for i := 0; i < len(k.List); i++ {
		if !ms.Equals(k.List[i], other.List[i]) {
			return false
		}
	}
	return true
}

func (k *InjectedKLabel) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*InjectedKLabel)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Label != other.Label {
		return false
	}
	return true
}

func (k *KToken) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*KToken)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Value == other.Value
}

func (k *KVariable) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*KVariable)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Name != other.Name {
		return false
	}
	return true
}

func (k *Map) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*Map)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for key, val := range k.Data {
		otherVal, found := other.Data[key]
		if !found {
			return false
		}
		if !ms.Equals(val, otherVal) {
			return false
		}
	}
	return true
}

func (k *List) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*List)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Sort != other.Sort {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for i := 0; i < len(k.Data); i++ {
		if !ms.Equals(k.Data[i], other.Data[i]) {
			return false
		}
	}
	return true
}

func (k *Set) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*Set)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for key := range k.Data {
		_, found := other.Data[key]
		if !found {
			return false
		}
	}
	return true
}

func (k *Array) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*Array)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Data.Equals(other.Data)
}

func (k *BigInt) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*BigInt)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	return k.Value.Cmp(other.Value) == 0
}

func (k *MInt) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*MInt)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	return k.Value == other.Value
}

func (k *Float) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*Float)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	return k.Value == other.Value
}

func (k *String) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*String)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	return k.Value == other.Value
}

// Equals ... Pointer comparison only for StringBuffer
func (k *StringBuffer) equals(ms *ModelState, arg KObject) bool {
	return k == arg
}

func (k *Bytes) equals(ms *ModelState, arg KObject) bool {
	other, typeOk := arg.(*Bytes)
	if !typeOk {
		panic("equals between different types should have been handled during reference Equals")
	}
	return bytes.Equal(k.Value, other.Value)
}

func (ms *ModelState) ksequenceEquals(ref1 KReference, ref2 KReference) bool {
	s1 := ms.KSequenceToSlice(ref1)
	s2 := ms.KSequenceToSlice(ref2)

	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if !ms.Equals(s1[i], s2[i]) {
			return false
		}
	}

	return true
}
