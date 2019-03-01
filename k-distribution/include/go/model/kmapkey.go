package %PACKAGE_MODEL%

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// KMapKey ... compact representation of a K item to be used as key in a map
type KMapKey struct {
	TypeName string
	Value    string
}

// KUsableAsKey ... A K Item that can be used as key in a map
type KUsableAsKey interface {
	AsMapKey() KMapKey
}

// AsMapKey ... Convert KToken to map key
func (k *KToken) AsMapKey() KMapKey {
	typeName := fmt.Sprintf("KToken(%s)", k.Sort.Name())
	return KMapKey{TypeName: typeName, Value: k.Value}
}

// AsMapKey ... Convert Int to map key
func (k *Int) AsMapKey() KMapKey {
	return KMapKey{TypeName: "Int", Value: k.Value.String()}
}

// AsMapKey ... Convert Bool to map key
func (k *Bool) AsMapKey() KMapKey {
	return KMapKey{TypeName: "Bool", Value: fmt.Sprintf("%t", k.Value)}
}

// AsMapKey ... Convert String to map key
func (k *String) AsMapKey() KMapKey {
	return KMapKey{TypeName: "String", Value: k.Value}
}

// String ... string representation of the key
func (mapKey KMapKey) String() string {
	return fmt.Sprintf("%s_%s", mapKey.TypeName, mapKey.Value)
}

// ToKItem ... convert a map key back to a regular K item
func (mapKey KMapKey) ToKItem() (K, error) {
	if strings.HasPrefix(mapKey.TypeName, "KToken(") && strings.HasSuffix(mapKey.TypeName, ")") {
		sortName := strings.TrimPrefix(mapKey.TypeName, "KToken(")
		sortName = strings.TrimSuffix(sortName, ")")
		return &KToken{Sort: ParseSort(sortName), Value: mapKey.Value}, nil
	}
	switch mapKey.TypeName {
	case "Int":
		return ParseInt(mapKey.Value)
	case "Bool":
		b, err := strconv.ParseBool(mapKey.Value)
		if err != nil {
			return NoResult, err
		}
		return ToBool(b), nil
	case "String":
		return NewString(mapKey.Value), nil
	default:
		return NoResult, errors.New("unable to convert KMapKey to K. Unknown type")
	}

}
