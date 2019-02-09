package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type stringHooksType int

const stringHooks stringHooksType = 0

func (stringHooksType) concat(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) lt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) le(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) gt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) ge(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) eq(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) ne(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) chr(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) find(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) rfind(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) length(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) substr(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) ord(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) int2string(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2int(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2base(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) base2string(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2token(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) token2string(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) float2string(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) uuid(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) floatFormat(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2float(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) replace(c1 m.K, c2 m.K, c3 m.K, c4 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) replaceAll(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) replaceFirst(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) countAllOccurrences(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) category(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) directionality(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) findChar(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) rfindChar(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

