package main

type mapHooksType int

const mapHooks mapHooksType = 0

// returns a map with 1 key to value mapping
func (mapHooksType) element(key K, val K, lbl KLabel, sort Sort, config K) (K, error) {
	m := make(map[K]K)
	m[key] = val
	return Map{Sort: sort, Label: lbl.collectionFor(), data: m}, nil
}

// returns an empty map
func (mapHooksType) unit(lbl KLabel, sort Sort, config K) (K, error) {
	var m map[K]K
	return Map{Sort: sort, Label: lbl.collectionFor(), data: m}, nil
}

func (mapHooksType) lookup(kmap K, key K, lbl KLabel, sort Sort, config K) (K, error) {
	if m, isMap := kmap.(Map); isMap {
		elem, found := m.data[key]
		if found {
			return elem, nil
		}
		return internedBottom, nil
	}

	if _, isBottom := kmap.(Map); isBottom {
		return internedBottom, nil
	}

	return noResult, &hookInvalidArgsError{}
}

func (mapHooksType) lookupOrDefault(kmap K, key K, defaultRes K, lbl KLabel, sort Sort, config K) (K, error) {
	if m, isMap := kmap.(Map); isMap {
		elem, found := m.data[key]
		if found {
			return elem, nil
		}
		return defaultRes, nil
	}

	if _, isBottom := kmap.(Map); isBottom {
		return internedBottom, nil
	}

	return noResult, &hookInvalidArgsError{}
}

func (mapHooksType) update(kmap K, newKey K, newValue K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	// implementing it as an "immutable" map
	// that is, creating a copy for each update (for now)
	// not the most efficient, not sure if necessary, but it is the safest
	newData := make(map[K]K)
	for oldKey, oldValue := range m.data {
		newData[oldKey] = oldValue
	}
	newData[newKey] = newValue
	return Map{Sort: m.Sort, Label: m.Label, data: newData}, nil
}

func (mapHooksType) remove(kmap K, key K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	// no updating of input map
	newData := make(map[K]K)
	for oldKey, oldValue := range m.data {
		if oldKey != key {
			newData[oldKey] = oldValue
		}
	}
	return Map{Sort: m.Sort, Label: m.Label, data: newData}, nil
}

func (mapHooksType) concat(kmap1 K, kmap2 K, lbl KLabel, sort Sort, config K) (K, error) {
	m1, isMap1 := kmap1.(Map)
	m2, isMap2 := kmap2.(Map)
	if !isMap1 || !isMap2 {
		return noResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]K)
	for key, value := range m1.data {
		data[key] = value
	}
	for key, value := range m2.data {
		m1Val, exists := m1.data[key]
		if exists {
			if m1Val != value {
				return noResult, &hookInvalidArgsError{}
			}
		} else {
			data[key] = value
		}
	}
	return Map{Sort: m1.Sort, Label: m1.Label, data: data}, nil
}

func (mapHooksType) difference(kmap1 K, kmap2 K, lbl KLabel, sort Sort, config K) (K, error) {
	m1, isMap1 := kmap1.(Map)
	m2, isMap2 := kmap2.(Map)
	if !isMap1 || !isMap2 {
		return noResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]K)
	for key, value := range m1.data {
		_, exists := m2.data[key]
		if !exists {
			data[key] = value
		}

	}
	return Map{Sort: m1.Sort, Label: m1.Label, data: data}, nil
}

func (mapHooksType) keys(kmap K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	keySet := make(map[K]bool)
	for key := range m.data {
		keySet[key] = true
	}
	return Set{Sort: sortSet, Label: lbl_Set_, data: keySet}, nil
}

func (mapHooksType) keysList(kmap K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	var keyList []K
	for key := range m.data {
		keyList = append(keyList, key)
	}
	return List{Sort: sortList, Label: lbl_List_, data: keyList}, nil
}

func (mapHooksType) inKeys(kmap K, key K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	_, keyExists := m.data[key]
	return Bool(keyExists), nil
}

func (mapHooksType) values(kmap K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	var valueList []K
	for _, value := range m.data {
		valueList = append(valueList, value)
	}
	return List{Sort: sortList, Label: lbl_List_, data: valueList}, nil
}

func (mapHooksType) choice(kmap K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	for key := range m.data {
		return key, nil
	}
	return noResult, &hookInvalidArgsError{}
}

func (mapHooksType) size(kmap K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	if !isMap {
		return noResult, &hookInvalidArgsError{}
	}
	return Int(len(m.data)), nil
}

func (mapHooksType) inclusion(kmap1 K, kmap2 K, lbl KLabel, sort Sort, config K) (K, error) {
	m1, isMap1 := kmap1.(Map)
	m2, isMap2 := kmap2.(Map)
	if !isMap1 || !isMap2 {
		return noResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return noResult, &hookInvalidArgsError{}
	}
	for m2Key := range m2.data {
		_, exists := m1.data[m2Key]
		if !exists {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func (mapHooksType) updateAll(kmap1 K, kmap2 K, lbl KLabel, sort Sort, config K) (K, error) {
	m1, isMap1 := kmap1.(Map)
	m2, isMap2 := kmap2.(Map)
	if !isMap1 || !isMap2 {
		return noResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]K)
	for key, value := range m1.data {
		data[key] = value
	}
	for key, value := range m2.data {
		data[key] = value
	}
	return Map{Sort: m1.Sort, Label: m1.Label, data: data}, nil
}

func (mapHooksType) removeAll(kmap K, kset K, lbl KLabel, sort Sort, config K) (K, error) {
	m, isMap := kmap.(Map)
	s, isSet := kset.(Set)
	if !isMap || !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]K)
	for key, value := range m.data {
		_, exists := s.data[key]
		if !exists {
			data[key] = value
		}

	}
	return Map{Sort: m.Sort, Label: m.Label, data: data}, nil
}
