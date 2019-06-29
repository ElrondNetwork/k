%COMMENT%

package %PACKAGE_MODEL%

import (
	"bytes"
)

// Equals performs deep comparison
func (ms *ModelState) Equals(arg1 K, arg2 K) bool {
	return arg1.equals(arg2)
}

func (k *KApply) equals(arg K) bool {
	other, typeOk := arg.(*KApply)
	if !typeOk {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	if len(k.List) != len(other.List) {
		return false
	}
	for i := 0; i < len(k.List); i++ {
		if !k.List[i].equals(other.List[i]) {
			return false
		}
	}
	return true
}

func (k *InjectedKLabel) equals(arg K) bool {
	other, typeOk := arg.(*InjectedKLabel)
	if !typeOk {
		return false
	}
	if k.Label != other.Label {
		return false
	}
	return true
}

func (k *KToken) equals(arg K) bool {
	other, typeOk := arg.(*KToken)
	if !typeOk {
		return false
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Value == other.Value
}

func (k *KVariable) equals(arg K) bool {
	other, typeOk := arg.(*KVariable)
	if !typeOk {
		return false
	}
	if k.Name != other.Name {
		return false
	}
	return true
}

func (k *Map) equals(arg K) bool {
	other, typeOk := arg.(*Map)
	if !typeOk {
		return false
	}
	if len(k.Data) != len(other.Data) {
		return false
	}
	for key, val := range k.Data {
		otherVal, found := other.Data[key]
		if !found {
			return false
		}
		if !val.equals(otherVal) {
			return false
		}
	}
	return true
}

func (k *List) equals(arg K) bool {
	other, typeOk := arg.(*List)
	if !typeOk {
		return false
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
		if !k.Data[i].equals(other.Data[i]) {
			return false
		}
	}
	return true
}

func (k *Set) equals(arg K) bool {
	other, typeOk := arg.(*Set)
	if !typeOk {
		return false
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

func (k *Array) equals(arg K) bool {
	other, typeOk := arg.(*Array)
	if !typeOk {
		return false
	}
	if k.Sort != other.Sort {
		return false
	}
	return k.Data.Equals(other.Data)
}

func (k *BigInt) equals(arg K) bool {
	other, typeOk := arg.(*BigInt)
	if !typeOk {
		return false
	}
	return k.Value.Cmp(other.Value) == 0
}

func (k *MInt) equals(arg K) bool {
	other, typeOk := arg.(*MInt)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

func (k *Float) equals(arg K) bool {
	other, typeOk := arg.(*Float)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

func (k *String) equals(arg K) bool {
	other, typeOk := arg.(*String)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

// Equals ... Pointer comparison only for StringBuffer
func (k *StringBuffer) equals(arg K) bool {
	return k == arg
}

func (k *Bytes) equals(arg K) bool {
	other, typeOk := arg.(*Bytes)
	if !typeOk {
		return false
	}
	return bytes.Equal(k.Value, other.Value)
}

func (k *Bool) equals(arg K) bool {
	other, typeOk := arg.(*Bool)
	if !typeOk {
		return false
	}
	return k.Value == other.Value
}

func (k *Bottom) equals(arg K) bool {
	_, typeOk := arg.(*Bottom)
	if !typeOk {
		return false
	}
	return true
}

func (k *KSequence) equals(arg K) bool {
	other, typeOk := arg.(*KSequence)
	if !typeOk {
		return false
	}
	if len(k.Ks) != len(other.Ks) {
		return false
	}
	for i := 0; i < len(k.Ks); i++ {
		if !k.Ks[i].equals(other.Ks[i]) {
			return false
		}
	}
	return true
}
