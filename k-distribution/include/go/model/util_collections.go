package %PACKAGE_MODEL%

import (
	"sort"
)

// MapKeyValuePair ... just a pair of key and value that was stored in a map
type MapKeyValuePair struct {
	KeyAsString string
	Key         K
	Value       K
}

// ToOrderedKeyValuePairs ... Yields a list of key-value pairs, ordered by the string representation of the keys
func (k *Map) ToOrderedKeyValuePairs() []MapKeyValuePair {
	result := make([]MapKeyValuePair, len(k.Data))

	var keysAsString []string
	stringKeysToPair := make(map[string]MapKeyValuePair)
	for key, val := range k.Data {
		keyAsString := key.String()
		keysAsString = append(keysAsString, keyAsString)
		keyAsK, err := key.ToKItem()
		if err != nil {
			panic(err)
		}
		pair := MapKeyValuePair{KeyAsString: keyAsString, Key: keyAsK, Value: val}
		stringKeysToPair[keyAsString] = pair
	}
	sort.Strings(keysAsString)
	for i, keyAsString := range keysAsString {
		pair, _ := stringKeysToPair[keyAsString]
		result[i] = pair
	}

	return result
}
