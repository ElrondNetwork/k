package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type listHooksType int

const listHooks listHooksType = 0

func (listHooksType) unit(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	data := make([]m.K, 0)
	return m.List{Sort: sort, Label: lbl.CollectionFor(), Data: data}, nil
}

func (listHooksType) element(e m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	data := make([]m.K, 1)
	data[0] = e
	return m.List{Sort: sort, Label: lbl.CollectionFor(), Data: data}, nil
}

func (listHooksType) concat(klist1 m.K, klist2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l1, isList1 := klist1.(m.List)
	l2, isList2 := klist2.(m.List)
	if !isList1 || !isList2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	data := make([]m.K, len(l1.Data)+len(l2.Data))
	for _, x := range l1.Data {
		data = append(data, x)
	}
	for _, x := range l2.Data {
		data = append(data, x)
	}
	return m.List{Sort: sort, Label: lbl.CollectionFor(), Data: data}, nil
}

func (listHooksType) in(e m.K, klist m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l, isList := klist.(m.List)
	if !isList {
		return m.NoResult, &hookInvalidArgsError{}
	}
	for _, x := range l.Data {
		if x == e {
			return m.Bool(true), nil
		}
	}
	return m.Bool(false), nil
}

func (listHooksType) get(klist m.K, index m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l, isList := klist.(m.List)
	i, isInt := index.(m.Int)
	if !isList || !isInt {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return l.Data[int(i)], nil
}

func (listHooksType) listRange(klist m.K, start m.K, end m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l, isList := klist.(m.List)
	si, isInt1 := start.(m.Int)
	ei, isInt2 := end.(m.Int)
	if !isList || !isInt1 || isInt2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.List{Sort: l.Sort, Label: l.Label, Data: l.Data[si:ei]}, nil
}

func (listHooksType) size(klist m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	l, isList := klist.(m.List)
	if !isList {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Int(len(l.Data)), nil
}

func (listHooksType) make(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) fill(c1 m.K, c2 m.K, c3 m.K, c4 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) update(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) updateAll(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}
