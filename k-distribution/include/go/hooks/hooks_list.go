package main

type listHooksType int

const listHooks listHooksType = 0

func (listHooksType) unit(lbl KLabel, sort Sort, config K) (K, error) {
	data := make([]K, 0)
	return List{Sort: sort, Label: lbl.collectionFor(), data: data}, nil
}

func (listHooksType) element(e K, lbl KLabel, sort Sort, config K) (K, error) {
	data := make([]K, 1)
	data[0] = e
	return List{Sort: sort, Label: lbl.collectionFor(), data: data}, nil
}

func (listHooksType) concat(klist1 K, klist2 K, lbl KLabel, sort Sort, config K) (K, error) {
	l1, isList1 := klist1.(List)
	l2, isList2 := klist2.(List)
	if !isList1 || !isList2 {
		return noResult, &hookInvalidArgsError{}
	}
	data := make([]K, len(l1.data)+len(l2.data))
	for _, x := range l1.data {
		data = append(data, x)
	}
	for _, x := range l2.data {
		data = append(data, x)
	}
	return List{Sort: sort, Label: lbl.collectionFor(), data: data}, nil
}

func (listHooksType) in(e K, klist K, lbl KLabel, sort Sort, config K) (K, error) {
	l, isList := klist.(List)
	if !isList {
		return noResult, &hookInvalidArgsError{}
	}
	for _, x := range l.data {
		if x == e {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}

func (listHooksType) get(klist K, index K, lbl KLabel, sort Sort, config K) (K, error) {
	l, isList := klist.(List)
	i, isInt := index.(Int)
	if !isList || !isInt {
		return noResult, &hookInvalidArgsError{}
	}
	return l.data[int(i)], nil
}

func (listHooksType) listRange(klist K, start K, end K, lbl KLabel, sort Sort, config K) (K, error) {
	l, isList := klist.(List)
	si, isInt1 := start.(Int)
	ei, isInt2 := end.(Int)
	if !isList || !isInt1 || isInt2 {
		return noResult, &hookInvalidArgsError{}
	}
	return List{Sort: l.Sort, Label: l.Label, data: l.data[si:ei]}, nil
}

func (listHooksType) size(klist K, lbl KLabel, sort Sort, config K) (K, error) {
	l, isList := klist.(List)
	if !isList {
		return noResult, &hookInvalidArgsError{}
	}
	return Int(len(l.data)), nil
}

func (listHooksType) make(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (listHooksType) fill(c1 K, c2 K, c3 K, c4 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (listHooksType) update(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (listHooksType) updateAll(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}
