package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type mapHooksType int

const mapHooks mapHooksType = 0

func isValidKey(key m.K) bool {
	if _, t := key.(m.KApply); t {
		return false
	}
	if _, t := key.(m.KSequence); t {
		return false
	}
	if _, t := key.(m.Map); t {
		return false
	}
	if _, t := key.(m.Set); t {
		return false
	}
	if _, t := key.(m.List); t {
		return false
	}
	if _, t := key.(m.Array); t {
		return false
	}
	return true
}

// returns a map with 1 key to value mapping
func (mapHooksType) element(key m.K, val m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	if !isValidKey(key) {
		panic("Invaid key")
	}
	mp := make(map[m.K]m.K)
	mp[key] = val
	return m.Map{Sort: sort, Label: lbl.CollectionFor(), Data: mp}, nil
}

// returns an empty map
func (mapHooksType) unit(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	var mp map[m.K]m.K
	return m.Map{Sort: sort, Label: lbl.CollectionFor(), Data: mp}, nil
}

func (mh mapHooksType) lookup(kmap m.K, key m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return mh.lookupOrDefault(kmap, key, m.InternedBottom, lbl, sort, config)
}

func (mapHooksType) lookupOrDefault(kmap m.K, key m.K, defaultRes m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	if mp, isMap := kmap.(m.Map); isMap {
		if !isValidKey(key) {
			return defaultRes, nil
		}
		elem, found := mp.Data[key]
		if found {
			return elem, nil
		}
		return defaultRes, nil
	}

	if _, isBottom := kmap.(m.Map); isBottom {
		return m.InternedBottom, nil
	}

	return m.NoResult, &hookInvalidArgsError{}
}

func (mapHooksType) update(kmap m.K, newKey m.K, newValue m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !isValidKey(newKey) {
		panic("Invaid key")
	}
	// implementing it as an "immutable" map
	// that is, creating a copy for each update (for now)
	// not the most efficient, not sure if necessary, but it is the safest
	newData := make(map[m.K]m.K)
	for oldKey, oldValue := range mp.Data {
		newData[oldKey] = oldValue
	}
	newData[newKey] = newValue
	return m.Map{Sort: mp.Sort, Label: mp.Label, Data: newData}, nil
}

func (mapHooksType) remove(kmap m.K, key m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	// no updating of input map
	newData := make(map[m.K]m.K)
	for oldKey, oldValue := range mp.Data {
		if oldKey != key {
			newData[oldKey] = oldValue
		}
	}
	return m.Map{Sort: mp.Sort, Label: mp.Label, Data: newData}, nil
}

func (mapHooksType) concat(kmap1 m.K, kmap2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	m1, isMap1 := kmap1.(m.Map)
	m2, isMap2 := kmap2.(m.Map)
	if !isMap1 || !isMap2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return m.NoResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]m.K)
	for key, value := range m1.Data {
		data[key] = value
	}
	for key, value := range m2.Data {
		m1Val, exists := m1.Data[key]
		if exists {
			if m1Val != value {
				return m.NoResult, &hookInvalidArgsError{}
			}
		} else {
			data[key] = value
		}
	}
	return m.Map{Sort: m1.Sort, Label: m1.Label, Data: data}, nil
}

func (mapHooksType) difference(kmap1 m.K, kmap2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	m1, isMap1 := kmap1.(m.Map)
	m2, isMap2 := kmap2.(m.Map)
	if !isMap1 || !isMap2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return m.NoResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]m.K)
	for key, value := range m1.Data {
		_, exists := m2.Data[key]
		if !exists {
			data[key] = value
		}

	}
	return m.Map{Sort: m1.Sort, Label: m1.Label, Data: data}, nil
}

func (mapHooksType) keys(kmap m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	keySet := make(map[m.K]bool)
	for key := range mp.Data {
		keySet[key] = true
	}
	return m.Set{Sort: m.SortSet, Label: m.KLabelForSet, Data: keySet}, nil
}

func (mapHooksType) keysList(kmap m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var keyList []m.K
	for key := range mp.Data {
		keyList = append(keyList, key)
	}
	return m.List{Sort: m.SortList, Label: m.KLabelForList, Data: keyList}, nil
}

func (mapHooksType) inKeys(kmap m.K, key m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	_, keyExists := mp.Data[key]
	return m.Bool(keyExists), nil
}

func (mapHooksType) values(kmap m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var valueList []m.K
	for _, value := range mp.Data {
		valueList = append(valueList, value)
	}
	return m.List{Sort: m.SortList, Label: m.KLabelForList, Data: valueList}, nil
}

func (mapHooksType) choice(kmap m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	for key := range mp.Data {
		return key, nil
	}
	return m.NoResult, &hookInvalidArgsError{}
}

func (mapHooksType) size(kmap m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	if !isMap {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.NewIntFromInt(len(mp.Data)), nil
}

func (mapHooksType) inclusion(kmap1 m.K, kmap2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	m1, isMap1 := kmap1.(m.Map)
	m2, isMap2 := kmap2.(m.Map)
	if !isMap1 || !isMap2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return m.NoResult, &hookInvalidArgsError{}
	}
	for m2Key := range m2.Data {
		_, exists := m1.Data[m2Key]
		if !exists {
			return m.Bool(false), nil
		}
	}
	return m.Bool(true), nil
}

func (mapHooksType) updateAll(kmap1 m.K, kmap2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	m1, isMap1 := kmap1.(m.Map)
	m2, isMap2 := kmap2.(m.Map)
	if !isMap1 || !isMap2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if m1.Label != m2.Label {
		return m.NoResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]m.K)
	for key, value := range m1.Data {
		data[key] = value
	}
	for key, value := range m2.Data {
		data[key] = value
	}
	return m.Map{Sort: m1.Sort, Label: m1.Label, Data: data}, nil
}

func (mapHooksType) removeAll(kmap m.K, kset m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	mp, isMap := kmap.(m.Map)
	s, isSet := kset.(m.Set)
	if !isMap || !isSet {
		return m.NoResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]m.K)
	for key, value := range mp.Data {
		_, exists := s.Data[key]
		if !exists {
			data[key] = value
		}

	}
	return m.Map{Sort: mp.Sort, Label: mp.Label, Data: data}, nil
}
