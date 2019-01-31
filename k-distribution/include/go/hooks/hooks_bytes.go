package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type bytesHooksType int

const bytesHooks bytesHooksType = 0

func (bytesHooksType) empty(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) bytes2int(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) int2bytes(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) bytes2string(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) string2bytes(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) substr(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) replaceAt(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) length(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) padRight(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) padLeft(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) reverse(c m.K,lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (bytesHooksType) concat(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

