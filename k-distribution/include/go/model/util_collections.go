%COMMENT%

package %PACKAGE_MODEL%

import (
	"sort"
)

// ExtractKApplyArgs ... checks that a K item is a KApply and returns its arguments if so
func ExtractKApplyArgs(subject K, expectedLabel KLabel, expectedArity int) ([]K, bool) {
	kapp, isKapp := subject.(*KApply)
	if !isKapp {
		return nil, false
	}
	if kapp.Label != expectedLabel || len(kapp.List) != expectedArity {
		return nil, false
	}
	return kapp.List, true
}

// ExtractListData ... checks that a K item is a list and returns its contents if so
func ExtractListData(subject K, expectedSort Sort, expectedLabel KLabel) ([]K, bool) {
	klist, isList := subject.(*List)
	if !isList {
		return nil, false
	}
	if klist.Sort != expectedSort || klist.Label != expectedLabel {
		return nil, false
	}
	return klist.Data, true
}

// ExtractMapData ... checks that a K item is a map and returns its contents if so
func ExtractMapData(subject K, expectedSort Sort, expectedLabel KLabel) (map[KMapKey]K, bool) {
	kmap, isMap := subject.(*Map)
	if !isMap {
		return nil, false
	}
	if kmap.Sort != expectedSort || kmap.Label != expectedLabel {
		return nil, false
	}
	return kmap.Data, true
}

// MapKeyValuePair ... just a pair of key and value that was stored in a map
type MapKeyValuePair struct {
	KeyAsString string
	Key         K
	Value       K
}

// ToOrderedKeyValuePairs ... Yields a list of key-value pairs
// ordered by the string representation of the keys
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

// ToOrderedElements ... Yields a list of the K items in the map (converted from KMapKey)
// ordered by the KMapKey string representation of the elements
func (k *Set) ToOrderedElements() []K {
	result := make([]K, len(k.Data))

	var keysAsString []string
	stringKeysToElem := make(map[string]K)
	for key := range k.Data {
		keyAsString := key.String()
		keysAsString = append(keysAsString, keyAsString)
		keyAsK, err := key.ToKItem()
		if err != nil {
			panic(err)
		}
		stringKeysToElem[keyAsString] = keyAsK
	}
	sort.Strings(keysAsString)
	for i, keyAsString := range keysAsString {
		elem, _ := stringKeysToElem[keyAsString]
		result[i] = elem
	}

	return result
}
