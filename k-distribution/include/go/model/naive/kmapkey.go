%COMMENT%

package %PACKAGE_MODEL%

import (
	"errors"
	"fmt"
	"strconv"
)

// KMapKey is a compact representation of a K item to be used as key in a map.
type KMapKey interface {
	toKItem() (KReference, error)
	String() string
}

// MapKey converts a K item to a map key, if possible
func (*ModelState) MapKey(k K) (KMapKey, bool) {
	return mapKey(k)
}

// ToKItem converts a map key back to a regular K item
func (*ModelState) ToKItem(mapKey KMapKey) (KReference, error) {
	return mapKey.toKItem()
}

// kmapKeyBasic ... representation of basic types: Int, String, Bool
type kmapKeyBasic struct {
	typeName string
	value    string
}

type kmapKeyKApply0 struct {
	label KLabel
}

type kmapKeyKApply1 struct {
	label KLabel
	arg1  KMapKey
}

type kmapBottom struct {
}

func mapKey(k K) (KMapKey, bool) {
	uak, implementsInterface := k.(usableAsKey)
	if !implementsInterface {
		return kmapBottom{}, false
	}
	return uak.convertToMapKey()
}

// KUsableAsKey ... A K Item that can be used as key in a map
type usableAsKey interface {
	convertToMapKey() (KMapKey, bool)
}

func (k *KToken) convertToMapKey() (KMapKey, bool) {
	return *k, true
}

func (k *KApply) convertToMapKey() (KMapKey, bool) {
	switch len(k.List) {
	case 0:
		return kmapKeyKApply0{label: k.Label}, true
	case 1:
		argAsKey, argOk := mapKey(k.List[0])
		if !argOk {
			return kmapBottom{}, false
		}
		return kmapKeyKApply1{label: k.Label, arg1: argAsKey}, true
	default:
		return kmapBottom{}, false
	}
}

func (k *BigInt) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "Int", value: k.Value.String()}, true
}

func (k *Bool) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "Bool", value: fmt.Sprintf("%t", k.Value)}, true
}

func (k *String) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "String", value: k.Value}, true
}

func (k *Bottom) convertToMapKey() (KMapKey, bool) {
	return kmapBottom{}, true
}

// String provides a string representation of the key
func (k KToken) String() string {
	return fmt.Sprintf("KToken(%s)_%s", k.Sort.Name(), k.Value)
}

// ToKItem converts a map key back to a regular K item
func (k KToken) ToKItem() (K, error) {
	return &k, nil
}

// String provides a string representation of the key
func (mapKey kmapKeyBasic) String() string {
	return fmt.Sprintf("%s_%s", mapKey.typeName, mapKey.value)
}

// ToKItem converts a map key back to a regular K item
func (mapKey kmapKeyBasic) toKItem() (K, error) {
	switch mapKey.typeName {
	case "Int":
		b, err := parseBigInt(mapKey.value)
		if err != nil {
			return IntZero, &parseIntError{parseVal: mapKey.value}
		}
		return &BigInt{Value: b}, nil
	case "Bool":
		b, err := strconv.ParseBool(mapKey.value)
		if err != nil {
			return NoResult, err
		}
		return ToKBool(b), nil
	case "String":
		return &String{Value: mapKey.value}, nil
	default:
		return NoResult, errors.New("unable to convert KMapKey to K. Unknown type")
	}

}

// ToKItem converts a map key back to a regular K item
func (k KToken) toKItem() (K, error) {
	return &k, nil
}

// String provides a string representation of the key
func (mapKey kmapKeyKApply0) String() string {
	return fmt.Sprintf("KApply(%s)", mapKey.label.Name())
}

// ToKItem converts a map key back to a regular K item
func (mapKey kmapKeyKApply0) toKItem() (K, error) {
	return &KApply{Label: mapKey.label, List: nil}, nil
}

// String provides a string representation of the key
func (mapKey kmapKeyKApply1) String() string {
	return fmt.Sprintf("KApply(%s)_%s", mapKey.label.Name(), mapKey.arg1.String())
}

// ToKItem converts a map key back to a regular K item
func (mapKey kmapKeyKApply1) toKItem() (K, error) {
	argKItem, err := mapKey.arg1.toKItem()
	if err != nil {
		return NoResult, err
	}
	return &KApply{Label: mapKey.label, List: []K{argKItem}}, nil
}

// String provides a string representation of the key
func (mapKey kmapBottom) String() string {
	return "Bottom"
}

// ToKItem converts a map key back to a regular K item
func (mapKey kmapBottom) toKItem() (K, error) {
	return InternedBottom, nil
}
