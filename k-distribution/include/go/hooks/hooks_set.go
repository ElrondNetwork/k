package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type setHooksType int

const setHooks setHooksType = 0

func (setHooksType) in(e m.K, kset m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s, isSet := kset.(m.Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	_, exists := s.Data[e]
	return m.Bool(exists), nil
}

func (setHooksType) unit(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	var data map[m.K]bool
	return m.Set{Sort: sort, Label: lbl.CollectionFor(), Data: data}, nil
}

// returns a set with 1 element
func (setHooksType) element(e m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	data := make(map[m.K]bool)
	data[e] = true
	return m.Set{Sort: sort, Label: lbl.CollectionFor(), Data: data}, nil
}

func (setHooksType) concat(kset1 m.K, kset2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s1, isSet1 := kset1.(m.Set)
	s2, isSet2 := kset2.(m.Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]bool)
	for e1 := range s1.Data {
		data[e1] = true
	}
	for e2 := range s2.Data {
		data[e2] = true
	}
	return m.Set{Sort: sort, Label: lbl, Data: data}, nil
}

func (setHooksType) difference(kset1 m.K, kset2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s1, isSet1 := kset1.(m.Set)
	s2, isSet2 := kset2.(m.Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]bool)
	for e1 := range s1.Data {
		_, existsInS2 := s2.Data[e1]
		if !existsInS2 {
			data[e1] = true
		}
	}
	return m.Set{Sort: sort, Label: lbl, Data: data}, nil
}

// tests if kset1 is a subset of kset2
func (setHooksType) inclusion(kset1 m.K, kset2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s1, isSet1 := kset1.(m.Set)
	s2, isSet2 := kset2.(m.Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	for e1 := range s1.Data {
		_, existsInS2 := s2.Data[e1]
		if !existsInS2 {
			return m.Bool(false), nil
		}
	}
	return m.Bool(true), nil
}

func (setHooksType) intersection(kset1 m.K, kset2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s1, isSet1 := kset1.(m.Set)
	s2, isSet2 := kset2.(m.Set)
	if !isSet1 || !isSet2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]bool)
	for e1 := range s1.Data {
		data[e1] = true
	}
	for e2 := range s2.Data {
		data[e2] = true
	}
	return m.Set{Sort: sort, Label: lbl, Data: data}, nil
}

func (setHooksType) choice(kset m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s, isSet := kset.(m.Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	for e := range s.Data {
		return e, nil
	}
	return noResult, &hookInvalidArgsError{}
}

func (setHooksType) size(kset m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s, isSet := kset.(m.Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	return m.Int(len(s.Data)), nil
}

func (setHooksType) set2list(kset m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	s, isSet := kset.(m.Set)
	if !isSet {
		return noResult, &hookInvalidArgsError{}
	}
	var list []m.K
	for e := range s.Data {
		list = append(list, e)
	}
	return m.List{Sort: m.SortList, Label: m.KLabelForList, Data: list}, nil
}

func (setHooksType) list2set(klist m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l, isList := klist.(m.List)
	if !isList {
		return noResult, &hookInvalidArgsError{}
	}
	data := make(map[m.K]bool)
	for _, e := range l.Data {
		data[e] = true
	}
	return m.Set{Sort: m.SortSet, Label: m.KLabelForSet, Data: data}, nil
}
