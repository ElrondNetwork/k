package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"strings"
)

type bufferHooksType int

const bufferHooks bufferHooksType = 0

func (bufferHooksType) empty(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return &m.StringBuffer{Value: strings.Builder{}}, nil
}

func (bufferHooksType) concat(kbuf m.K, kstr m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	buf, ok1 := kbuf.(*m.StringBuffer)
	str, ok2 := kstr.(*m.String)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	buf.Value.WriteString(str.Value)
	return buf, nil
}

func (bufferHooksType) toString(kbuf m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	buf, ok := kbuf.(*m.StringBuffer)
	if !ok {
		return invalidArgsResult()
	}
	return m.NewString(buf.Value.String()), nil
}
