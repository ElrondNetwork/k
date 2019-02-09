package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type arrayHooksType int

const arrayHooks arrayHooksType = 0

func (arrayHooksType) make(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) makeEmpty(len m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) lookup(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) remove(arr m.K, v m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) update(arr m.K, index m.K, newVal m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) updateAll(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) fill(c1 m.K, c2 m.K, c3 m.K, c4 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) inKeys(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (arrayHooksType) ctor(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}
