package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
    "strconv"
    "strings"
)

type stringHooksType int

const stringHooks stringHooksType = 0

func (stringHooksType) concat(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	k1, ok1 := c1.(*m.String)
	k2, ok2 := c2.(*m.String)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.NewString(k1.String() + k2.String()), nil
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
	k1, ok1 := c1.(*m.String)
	k2, ok2 := c2.(*m.String)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToBool(k1.String() == k2.String()), nil
}

func (stringHooksType) ne(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	k1, ok1 := c1.(*m.String)
	k2, ok2 := c2.(*m.String)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToBool(k1.String() != k2.String()), nil
}

func (stringHooksType) chr(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i, ok := c.(*m.Int)
	if !ok {
		return invalidArgsResult()
	}
	r := rune(i.Value.Uint64())
	return m.NewString(string(r)), nil
}

func (stringHooksType) find(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	str, ok1 := c1.(*m.String)
	substr, ok2 := c2.(*m.String)
	firstIdx, ok3 := c3.(*m.Int)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !firstIdx.Value.IsUint64() {
		return invalidArgsResult()
	}
	firstIdxInt := firstIdx.Value.Uint64()
	if firstIdxInt > uint64(len(str.String())) {
		return invalidArgsResult()
	}

	result := strings.Index(str.String()[firstIdxInt:], substr.String())
	if result == -1 {
		return m.IntMinusOne, nil
	}
	return m.NewIntFromUint64(firstIdxInt + uint64(result)), nil
}

func (stringHooksType) rfind(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	str, ok1 := c1.(*m.String)
	substr, ok2 := c2.(*m.String)
	lastIdx, ok3 := c3.(*m.Int)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !lastIdx.Value.IsUint64() {
		return invalidArgsResult()
	}
	lastIdxInt := lastIdx.Value.Uint64()
	if lastIdxInt > uint64(len(str.String())) {
		return invalidArgsResult()
	}
	result := strings.LastIndex(str.String()[0:lastIdxInt], substr.String())
	if result == -1 {
		return m.IntMinusOne, nil
	}
	return m.NewIntFromInt(result), nil
}

func (stringHooksType) length(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	k, ok := c.(*m.String)
	if !ok {
		return invalidArgsResult()
	}
	return m.NewIntFromInt(len(k.String())), nil
}

func (stringHooksType) substr(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	str, ok1 := c1.(*m.String)
	from, ok2 := c2.(*m.Int) // from is inclusive
	to, ok3 := c3.(*m.Int)   // to is exclusive
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !from.Value.IsUint64() || !to.Value.IsUint64() {
		return invalidArgsResult()
	}
	fromInt := from.Value.Uint64()
	toInt := to.Value.Uint64()
	length := uint64(len(str.String()))
	if fromInt > toInt || fromInt > length {
		return invalidArgsResult()
	}
	if toInt > length {
		toInt = length
	}
	return m.NewString(str.String()[fromInt:toInt]), nil
}

func (stringHooksType) ord(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) int2string(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2int(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2base(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) base2string(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2token(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	k, ok := c.(*m.String)
	if !ok {
		return invalidArgsResult()
	}
	return &m.KToken{Sort: sort, Value: k.String()}, nil
}

func (stringHooksType) token2string(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	if k, typeOk := c.(*m.KToken); typeOk {
		return m.NewString(k.Value), nil
	}
	if k, typeOk := c.(*m.Bool); typeOk {
		return m.NewString(strconv.FormatBool(k.Value)), nil
	}
	if k, typeOk := c.(*m.String); typeOk {
		return k, nil // TODO: should do escaping
	}
	if k, typeOk := c.(*m.Int); typeOk {
		return m.NewString(k.Value.String()), nil
	}
	if _, typeOk := c.(*m.Float); typeOk {
		return m.NoResult, &hookNotImplementedError{}
	}

	return invalidArgsResult()
}

func (stringHooksType) float2string(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) uuid(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) floatFormat(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2float(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
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

func (stringHooksType) category(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) directionality(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) findChar(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) rfindChar(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}
