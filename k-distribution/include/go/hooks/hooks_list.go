%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type listHooksType int

const listHooks listHooksType = 0

func (listHooksType) unit(lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	data := make([]m.KReference, 0)
	return interpreter.Model.NewList(sort, m.CollectionFor(lbl), data), nil
}

func (listHooksType) element(e m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	data := make([]m.KReference, 1)
	data[0] = e
	return interpreter.Model.NewList(sort, m.CollectionFor(lbl), data), nil
}

func (listHooksType) concat(klist1 m.KReference, klist2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	l1, isList1 := interpreter.Model.GetListObject(klist1)
	l2, isList2 := interpreter.Model.GetListObject(klist2)
	if !isList1 || !isList2 {
		return invalidArgsResult()
	}
	if len(l1.Data) == 0 {
		return klist2, nil
	}
	if len(l2.Data) == 0 {
		return klist1, nil
	}
	data := append(l1.Data, l2.Data...)
	return interpreter.Model.NewList(sort, m.CollectionFor(lbl), data), nil
}

func (listHooksType) in(e m.KReference, klist m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	l, isList := interpreter.Model.GetListObject(klist)
	if !isList {
		return invalidArgsResult()
	}
	for _, x := range l.Data {
		if x == e {
			return m.BoolTrue, nil
		}
	}
	return m.BoolFalse, nil
}

func (listHooksType) get(klist m.KReference, kindex m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	l, isList := interpreter.Model.GetListObject(klist)
	i, isInt := interpreter.Model.GetBigIntObject(kindex)
	if !isList || !isInt {
		return invalidArgsResult()
	}
	if !i.Value.IsInt64() {
		return invalidArgsResult()
	}
	index := int(i.Value.Int64())
	if index < 0 {
		index = len(l.Data) + index // count from the end, e.g. -1 is the last element
	}
	if index < 0 || index >= len(l.Data) {
		return m.InternedBottom, nil
	}
	return l.Data[index], nil
}

func (listHooksType) listRange(klist m.KReference, start m.KReference, end m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	l, isList := interpreter.Model.GetListObject(klist)
	si, isInt1 := interpreter.Model.GetBigIntObject(start)
	ei, isInt2 := interpreter.Model.GetBigIntObject(end)
	if !isList || !isInt1 || !isInt2 || !si.Value.IsUint64() || !ei.Value.IsUint64() {
		return invalidArgsResult()
	}
	siUint := si.Value.Uint64()
	eiUint := ei.Value.Uint64()
	return interpreter.Model.NewList(l.Sort, l.Label, l.Data[siUint:eiUint]), nil
}

func (listHooksType) size(klist m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	l, isList := interpreter.Model.GetListObject(klist)
	if !isList {
		return invalidArgsResult()
	}
	return interpreter.Model.FromInt(len(l.Data)), nil
}

func (listHooksType) make(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) fill(c1 m.KReference, c2 m.KReference, c3 m.KReference, c4 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) update(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (listHooksType) updateAll(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}
