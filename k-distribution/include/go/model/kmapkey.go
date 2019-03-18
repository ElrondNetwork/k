package %PACKAGE_MODEL%

import (
	"errors"
	"fmt"
	"strconv"
)

// KMapKey ... compact representation of a K item to be used as key in a map
type KMapKey interface {
	String() string
	ToKItem() (K, error)
}

// kmapKeyBasic ... representation of basic types: Int, String, Bool
type kmapKeyBasic struct {
	typeName string
	value    string
}

type kmapKeyToken struct {
	sortName string
	value    string
}

type kmapKeyKApply1 struct {
	labelName string
	arg1      KMapKey
}

type kmapNoResult struct {
}

// KUsableAsKey ... A K Item that can be used as key in a map
type usableAsKey interface {
	convertToMapKey() (KMapKey, bool)
}

// MapKey ... converts a K item to a map key, if possible
func MapKey(k K) (KMapKey, bool) {
	uak, implementsInterface := k.(usableAsKey)
	if !implementsInterface {
		return kmapNoResult{}, false
	}
	return uak.convertToMapKey()
}

func (k *KToken) convertToMapKey() (KMapKey, bool) {
	return kmapKeyToken{sortName: k.Sort.Name(), value: k.Value}, true
}

func (k *KApply) convertToMapKey() (KMapKey, bool) {
	if len(k.List) != 1 {
		return kmapNoResult{}, false
	}
	argAsKey, argOk := MapKey(k.List[0])
	if !argOk {
		return kmapNoResult{}, false
	}
	return kmapKeyKApply1{labelName: k.Label.Name(), arg1: argAsKey}, true
}

func (k *Int) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "Int", value: k.Value.String()}, true
}

func (k *Bool) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "Bool", value: fmt.Sprintf("%t", k.Value)}, true
}

func (k *String) convertToMapKey() (KMapKey, bool) {
	return kmapKeyBasic{typeName: "String", value: k.Value}, true
}

// String ... string representation of the key
func (mapKey kmapKeyBasic) String() string {
	return fmt.Sprintf("%s_%s", mapKey.typeName, mapKey.value)
}

// ToKItem ... convert a map key back to a regular K item
func (mapKey kmapKeyBasic) ToKItem() (K, error) {
	switch mapKey.typeName {
	case "Int":
		return ParseInt(mapKey.value)
	case "Bool":
		b, err := strconv.ParseBool(mapKey.value)
		if err != nil {
			return NoResult, err
		}
		return ToBool(b), nil
	case "String":
		return NewString(mapKey.value), nil
	default:
		return NoResult, errors.New("unable to convert KMapKey to K. Unknown type")
	}

}

// String ... string representation of the key
func (mapKey kmapKeyToken) String() string {
	return fmt.Sprintf("KToken(%s)_%s", mapKey.sortName, mapKey.value)
}

// ToKItem ... convert a map key back to a regular K item
func (mapKey kmapKeyToken) ToKItem() (K, error) {
	return &KToken{Sort: ParseSort(mapKey.sortName), Value: mapKey.value}, nil
}

// String ... string representation of the key
func (mapKey kmapKeyKApply1) String() string {
	return fmt.Sprintf("KApply(%s)_%s", mapKey.labelName, mapKey.arg1.String())
}

// ToKItem ... convert a map key back to a regular K item
func (mapKey kmapKeyKApply1) ToKItem() (K, error) {
	argKItem, err := mapKey.arg1.ToKItem()
	if err != nil {
		return NoResult, err
	}
	return &KApply{Label: ParseKLabel(mapKey.labelName), List: []K{argKItem}}, nil
}

// String ... string representation of the key
func (mapKey kmapNoResult) String() string {
	panic("Should never be called")
}

// ToKItem ... convert a map key back to a regular K item
func (mapKey kmapNoResult) ToKItem() (K, error) {
	panic("Should never be called")
}
