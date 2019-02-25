package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type arrayHooksType int

const arrayHooks arrayHooksType = 0

func (arrayHooksType) make(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	len, ok := c1.(m.Int)
	if !ok {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !len.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	lenUint := len.Value.Uint64()
	data := make([]*m.K, lenUint)
	for i := uint64(0); i < lenUint; i++ {
		data[i] = &c2
	}
	return m.Array{Sort: sort, Data: data, Default: &c2}, nil
}

func (t arrayHooksType) makeEmpty(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return t.make(c, m.InternedBottom, lbl, sort, config)
}

func (t arrayHooksType) ctor(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return t.makeEmpty(c2, lbl, sort, config)
}

func (arrayHooksType) lookup(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	arr, ok1 := c1.(m.Array)
	idx, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !idx.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	idxInt := idx.Value.Uint64()
	if idxInt >= uint64(len(arr.Data)) {
		return *arr.Default, nil
	}
	return *arr.Data[idxInt], nil
}

func (arrayHooksType) remove(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	arr, ok1 := c1.(m.Array)
	idx, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !idx.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	idxInt := idx.Value.Uint64()
	if idxInt < uint64(len(arr.Data)) {
		arr.Data[idxInt] = arr.Default
	}
	return arr, nil
}

func (arrayHooksType) update(c1 m.K, c2 m.K, newVal m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	arr, ok1 := c1.(m.Array)
	idx, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !idx.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	idxInt := idx.Value.Uint64()
	if idxInt >= 0 || idxInt < uint64(len(arr.Data)) {
		arr.Data[idxInt] = &newVal
	}
	return arr, nil
}

func (arrayHooksType) updateAll(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	arr, ok1 := c1.(m.Array)
	idx, ok2 := c2.(m.Int)
	list, ok3 := c3.(m.List)
	if !ok1 || !ok2 || !ok3 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !idx.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	idxInt := idx.Value.Uint64()
	for i := uint64(0); i < uint64(len(list.Data)) && idxInt+i < uint64(len(arr.Data)); i++ {
		arr.Data[idxInt+i] = &list.Data[i]
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) fill(c1 m.K, c2 m.K, c3 m.K, elt m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	arr, ok1 := c1.(m.Array)
	from, ok2 := c2.(m.Int)
	to, ok3 := c3.(m.Int)
	if !ok1 || !ok2 || !ok3 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !from.Value.IsUint64() || !to.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	fromInt := from.Value.Uint64()
	toInt := to.Value.Uint64()
	for i := fromInt; i < toInt && i < uint64(len(arr.Data)); i++ {
		arr.Data[i] = &elt
	}
	return arr, nil
}

func (arrayHooksType) inKeys(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	idx, ok2 := c1.(m.Int)
	arr, ok1 := c2.(m.Array)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !idx.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	idxInt := idx.Value.Uint64()
	if idxInt >= uint64(len(arr.Data)) {
		return m.Bool(false), nil
	}
	if *arr.Data[idxInt] == *arr.Default {
		return m.Bool(true), nil
	}
	return m.Bool(false), nil
}
