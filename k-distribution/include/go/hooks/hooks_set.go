package main

type setHooksType int

const setHooks setHooksType = 0

func (setHooksType) in(e K, kset K, lbl KLabel, sort Sort, config K) (K, error) {
	s, isSet := kset.(Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	_, exists := s.data[e]
	return Bool(exists), nil
}

func (setHooksType) unit(lbl KLabel, sort Sort, config K) (K, error) {
	var data map[K]bool
	return Set{Sort: sort, Label: lbl.collectionFor(), data: data}, nil
}

// returns a set with 1 element
func (setHooksType) element(e K, lbl KLabel, sort Sort, config K) (K, error) {
	data := make(map[K]bool)
	data[e] = true
	return Set{Sort: sort, Label: lbl.collectionFor(), data: data}, nil
}

func (setHooksType) concat(kset1 K, kset2 K, lbl KLabel, sort Sort, config K) (K, error) {
	s1, isSet1 := kset1.(Set)
	s2, isSet2 := kset2.(Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]bool)
	for e1 := range s1.data {
		data[e1] = true
	}
	for e2 := range s2.data {
		data[e2] = true
	}
	return Set{Sort: sort, Label: lbl, data: data}, nil
}

func (setHooksType) difference(kset1 K, kset2 K, lbl KLabel, sort Sort, config K) (K, error) {
	s1, isSet1 := kset1.(Set)
	s2, isSet2 := kset2.(Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]bool)
	for e1 := range s1.data {
		_, existsInS2 := s2.data[e1]
		if !existsInS2 {
			data[e1] = true
		}
	}
	return Set{Sort: sort, Label: lbl, data: data}, nil
}

// tests if kset1 is a subset of kset2
func (setHooksType) inclusion(kset1 K, kset2 K, lbl KLabel, sort Sort, config K) (K, error) {
	s1, isSet1 := kset1.(Set)
	s2, isSet2 := kset2.(Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	for e1 := range s1.data {
		_, existsInS2 := s2.data[e1]
		if !existsInS2 {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func (setHooksType) intersection(kset1 K, kset2 K, lbl KLabel, sort Sort, config K) (K, error) {
	s1, isSet1 := kset1.(Set)
	s2, isSet2 := kset2.(Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]bool)
	for e1 := range s1.data {
		data[e1] = true
	}
	for e2 := range s2.data {
		data[e2] = true
	}
	return Set{Sort: sort, Label: lbl, data: data}, nil
}

func (setHooksType) choice(kset K, lbl KLabel, sort Sort, config K) (K, error) {
	s, isSet := kset.(Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	for e := range s.data {
		return e, nil
	}
	return noResult, &hookInvalidArgsError{}
}

func (setHooksType) size(kset K, lbl KLabel, sort Sort, config K) (K, error) {
	s, isSet := kset.(Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	return Int(len(s.data)), nil
}

func (setHooksType) set2list(kset K, lbl KLabel, sort Sort, config K) (K, error) {
	s, isSet := kset.(Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	var list []K
	for e := range s.data {
		list = append(list, e)
	}
	return List{Sort: sortList, Label: klabelForList, data: list}, nil
}

func (setHooksType) list2set(klist K, lbl KLabel, sort Sort, config K) (K, error) {
	l, isList := klist.(List)
	if !isList {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[K]bool)
	for _, e := range l.data {
		data[e] = true
	}
	return Set{Sort: sortSet, Label: klabelForSet, data: data}, nil
}
