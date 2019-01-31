package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type bufferHooksType int

const bufferHooks bufferHooksType = 0

func (bufferHooksType) empty(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bufferHooksType) concat(buf m.K, elem m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bufferHooksType) toString(buf m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

