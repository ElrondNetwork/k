%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type arrayHooksType int

const arrayHooks arrayHooksType = 0

func (arrayHooksType) make(maxSize m.KReference, defValue m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	maxSizeInt, ok := interpreter.Model.GetBigIntObject(maxSize)
	if !ok {
		return invalidArgsResult()
	}
	if !maxSizeInt.Value.IsUint64() {
		return invalidArgsResult()
	}
	maxSizeUint := maxSizeInt.Value.Uint64()
	return interpreter.Model.NewArray(sort, interpreter.Model.MakeDynamicArray(maxSizeUint, defValue)), nil
}

func (t arrayHooksType) makeEmpty(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return t.make(c, m.InternedBottom, lbl, sort, config, interpreter)
}

func (t arrayHooksType) ctor(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return t.makeEmpty(c2, lbl, sort, config, interpreter)
}

func (arrayHooksType) lookup(karr m.KReference, kidx m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	arr, ok1 := interpreter.Model.GetArrayObject(karr)
	idx, ok2 := interpreter.Model.GetBigIntObject(kidx)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !idx.Value.IsUint64() {
		return invalidArgsResult()
	}
	idxUint := idx.Value.Uint64()
	return arr.Data.Get(idxUint)
}

func (arrayHooksType) remove(karr m.KReference, kidx m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	arr, ok1 := interpreter.Model.GetArrayObject(karr)
	idx, ok2 := interpreter.Model.GetBigIntObject(kidx)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !idx.Value.IsUint64() {
		return invalidArgsResult()
	}
	idxUint := idx.Value.Uint64()
	err := arr.Data.Set(idxUint, arr.Data.Default)
	if err != nil {
		return m.NoResult, err
	}
	return karr, nil
}

func (arrayHooksType) update(karr m.KReference, kidx m.KReference, newVal m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	arr, ok1 := interpreter.Model.GetArrayObject(karr)
	idx, ok2 := interpreter.Model.GetBigIntObject(kidx)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !idx.Value.IsUint64() {
		return invalidArgsResult()
	}
	idxUint := idx.Value.Uint64()
	err := arr.Data.Set(idxUint, newVal)
	if err != nil {
		return m.NoResult, err
	}
	return karr, nil
}

func (arrayHooksType) updateAll(karr m.KReference, kidx m.KReference, klist m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	arr, ok1 := interpreter.Model.GetArrayObject(karr)
	idx, ok2 := interpreter.Model.GetBigIntObject(kidx)
	list, ok3 := interpreter.Model.GetListObject(klist)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !idx.Value.IsUint64() {
		return invalidArgsResult()
	}
	idxUint := idx.Value.Uint64()
	listLen := uint64(len(list.Data))
	arr.Data.UpgradeSize(idxUint + listLen - 1) // upgrade size all at once
	for i := uint64(0); i < listLen && idxUint+i < arr.Data.MaxSize; i++ {
		err := arr.Data.Set(idxUint+i, list.Data[i])
		if err != nil {
			return m.NoResult, err
		}
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) fill(karr m.KReference, kfrom m.KReference, kto m.KReference, elt m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	arr, ok1 := interpreter.Model.GetArrayObject(karr)
	from, ok2 := interpreter.Model.GetBigIntObject(kfrom)
	to, ok3 := interpreter.Model.GetBigIntObject(kto)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !from.Value.IsUint64() || !to.Value.IsUint64() {
		return invalidArgsResult()
	}
	fromInt := from.Value.Uint64()
	toInt := to.Value.Uint64()
	for i := fromInt; i < toInt && i < arr.Data.MaxSize; i++ {
		arr.Data.Set(i, elt)
	}
	return karr, nil
}

func (arrayHooksType) inKeys(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	idx, ok2 := interpreter.Model.GetBigIntObject(c1)
	arr, ok1 := interpreter.Model.GetArrayObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !idx.Value.IsUint64() {
		return invalidArgsResult()
	}
	idxUint := idx.Value.Uint64()
	val, err := arr.Data.Get(idxUint)
	if err != nil {
		return m.NoResult, err
	}
	hasValue := !interpreter.Model.Equals(val, arr.Data.Default)
	return m.ToKBool(hasValue), nil
}
