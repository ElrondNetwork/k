%COMMENT%

package %PACKAGE_MODEL%

import (
	"math/big"
)

// DeepCopy yields a fresh copy of the K item given as argument.
func (ms *ModelState) DeepCopy(ref KReference) KReference {
	return ref.deepCopy()
}

func (k *KApply) deepCopy() K {
	listCopy := make([]K, len(k.List))
	for i, child := range k.List {
		listCopy[i] = child.deepCopy()
	}
	return &KApply{Label: k.Label, List: listCopy}
}

func (k *InjectedKLabel) deepCopy() K {
	return &InjectedKLabel{Label: k.Label}
}

func (k *KToken) deepCopy() K {
	return &KToken{Sort: k.Sort, Value: k.Value}
}

func (k *KVariable) deepCopy() K {
	return &KVariable{Name: k.Name}
}

func (k *Map) deepCopy() K {
	mapCopy := make(map[KMapKey]K)
	for key, val := range k.Data {
		mapCopy[key] = val.deepCopy()
	}
	return &Map{Data: mapCopy}
}

func (k *List) deepCopy() K {
	listCopy := make([]K, len(k.Data))
	for i, elem := range k.Data {
		listCopy[i] = elem.deepCopy()
	}
	return &List{Sort: k.Sort, Label: k.Label, Data: listCopy}
}

func (k *Set) deepCopy() K {
	mapCopy := make(map[KMapKey]bool)
	for key := range k.Data {
		mapCopy[key] = true
	}
	return &Set{Data: mapCopy}
}

func (k *Array) deepCopy() K {
	return k // TODO: not implemented
}

func (k *BigInt) deepCopy() K {
	intCopy := new(big.Int)
	intCopy.Set(k.Value)
	return &BigInt{Value: intCopy}
}

func (k *MInt) deepCopy() K {
	return k // not implemented
}

func (k *Float) deepCopy() K {
	return k // not implemented
}

func (k *String) deepCopy() K {
	return &String{Value: k.Value}
}

func (k *StringBuffer) deepCopy() K {
	return k // no deep copy needed here
}

func (k *Bytes) deepCopy() K {
	bytesCopy := make([]byte, len(k.Value))
	copy(bytesCopy, k.Value)
	return &Bytes{Value: bytesCopy}
}

func (k *Bool) deepCopy() K {
	return &Bool{Value: k.Value}
}

func (k *Bottom) deepCopy() K {
	return &Bottom{}
}

func (k *KSequence) deepCopy() K {
	ksCopy := make([]K, len(k.Ks))
	for i, elem := range k.Ks {
		ksCopy[i] = elem.deepCopy()
	}
	return &KSequence{Ks: ksCopy}
}
