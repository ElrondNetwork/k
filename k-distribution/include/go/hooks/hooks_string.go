%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"math/big"
    "strconv"
    "strings"
)

type stringHooksType int

const stringHooks stringHooksType = 0

func (stringHooksType) concat(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return interpreter.Model.NewString(str1 + str2), nil
}

func (stringHooksType) lt(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 < str2), nil
}

func (stringHooksType) le(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 <= str2), nil
}

func (stringHooksType) gt(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 > str2), nil
}

func (stringHooksType) ge(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 >= str2), nil
}

func (stringHooksType) eq(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 == str2), nil
}

func (stringHooksType) ne(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str1, ok1 := interpreter.Model.GetString(c1)
	str2, ok2 := interpreter.Model.GetString(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(str1 != str2), nil
}

func (stringHooksType) chr(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i, ok := interpreter.Model.GetBigIntObject(c)
	if !ok {
		return invalidArgsResult()
	}

	b := byte(i.Value.Uint64())
	bytes := []byte{b}
	return interpreter.Model.NewString(string(bytes)), nil
}

func (stringHooksType) find(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok1 := interpreter.Model.GetString(c1)
	substr, ok2 := interpreter.Model.GetString(c2)
	firstIdx, ok3 := interpreter.Model.GetBigIntObject(c3)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !firstIdx.Value.IsUint64() {
		return invalidArgsResult()
	}
	firstIdxInt := firstIdx.Value.Uint64()
	if firstIdxInt > uint64(len(str)) {
		return invalidArgsResult()
	}

	result := strings.Index(str[firstIdxInt:], substr)
	if result == -1 {
		return m.IntMinusOne, nil
	}
	return interpreter.Model.FromUint64(firstIdxInt + uint64(result)), nil
}

func (stringHooksType) rfind(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok1 := interpreter.Model.GetString(c1)
	substr, ok2 := interpreter.Model.GetString(c2)
	lastIdx, ok3 := interpreter.Model.GetBigIntObject(c3)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !lastIdx.Value.IsUint64() {
		return invalidArgsResult()
	}
	lastIdxInt := lastIdx.Value.Uint64()
	if lastIdxInt > uint64(len(str)) {
		return invalidArgsResult()
	}
	result := strings.LastIndex(str[0:lastIdxInt], substr)
	if result == -1 {
		return m.IntMinusOne, nil
	}
	return interpreter.Model.FromInt(result), nil
}

func (stringHooksType) length(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	k, ok := interpreter.Model.GetString(c)
	if !ok {
		return invalidArgsResult()
	}
	return interpreter.Model.FromInt(len(k)), nil
}

func (stringHooksType) substr(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok1 := interpreter.Model.GetString(c1)
	from, ok2 := interpreter.Model.GetBigIntObject(c2) // from is inclusive
	to, ok3 := interpreter.Model.GetBigIntObject(c3)   // to is exclusive
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if !from.Value.IsUint64() || !to.Value.IsUint64() {
		return invalidArgsResult()
	}
	fromInt := from.Value.Uint64()
	toInt := to.Value.Uint64()
	length := uint64(len(str))
	if fromInt > toInt || fromInt > length {
		return invalidArgsResult()
	}
	if toInt > length {
		toInt = length
	}
	return interpreter.Model.NewString(str[fromInt:toInt]), nil
}

func (stringHooksType) ord(arg m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok := interpreter.Model.GetString(arg)
	if !ok {
		return invalidArgsResult()
	}
	asBytes := []byte(str)
	if len(asBytes) == 0 {
		return invalidArgsResult()
	}
	return interpreter.Model.IntFromByte(asBytes[0]), nil
}

func (stringHooksType) int2string(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i, ok := interpreter.Model.GetBigIntObject(c)
	if !ok {
		return invalidArgsResult()
	}
	return interpreter.Model.NewString(i.Value.String()), nil
}

func (stringHooksType) string2int(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2base(kstr m.KReference, kbase m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok1 := interpreter.Model.GetString(kstr)
	base, ok2 := interpreter.Model.GetBigIntObject(kbase)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !base.Value.IsUint64() {
		return invalidArgsResult()
	}
	baseVal := base.Value.Uint64()
	if baseVal < 2 || baseVal > 16 {
		return invalidArgsResult()
	}
	i := new(big.Int)
	var parseOk bool
	i, parseOk = i.SetString(str, int(baseVal))
	if !parseOk {
		return invalidArgsResult()
	}
	return interpreter.Model.FromBigInt(i), nil
}

func (stringHooksType) base2string(kint m.KReference, kbase m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i, ok1 := interpreter.Model.GetBigIntObject(kint)
	base, ok2 := interpreter.Model.GetBigIntObject(kbase)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if !base.Value.IsUint64() {
		return invalidArgsResult()
	}
	baseVal := base.Value.Uint64()
	if baseVal < 2 || baseVal > 16 {
		return invalidArgsResult()
	}
	str := i.Value.Text(int(baseVal))
	return interpreter.Model.NewString(str), nil
}

func (stringHooksType) string2token(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	str, ok := interpreter.Model.GetString(c)
	if !ok {
		return invalidArgsResult()
	}
	return interpreter.Model.NewKToken(sort, str), nil
}

func (stringHooksType) token2string(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	if k, typeOk := interpreter.Model.GetKTokenObject(c); typeOk {
		return interpreter.Model.NewString(k.Value), nil
	}
	if k, typeOk := m.CastToBool(c); typeOk {
		return interpreter.Model.NewString(strconv.FormatBool(k)), nil
	}
	if k, typeOk := interpreter.Model.GetString(c); typeOk {
		return interpreter.Model.NewString(k), nil // TODO: should do escaping
	}
	if k, typeOk := interpreter.Model.GetBigIntObject(c); typeOk {
		return interpreter.Model.NewString(k.Value.String()), nil
	}
	if _, typeOk := interpreter.Model.GetFloatObject(c); typeOk {
		return m.NoResult, &hookNotImplementedError{}
	}

	return invalidArgsResult()
}

func (stringHooksType) float2string(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) uuid(lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) floatFormat(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) string2float(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) replace(argS m.KReference, argToReplace m.KReference, argReplacement m.KReference, argCount m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	kS, ok1 := interpreter.Model.GetString(argS)
	kToReplace, ok2 := interpreter.Model.GetString(argToReplace)
	kReplacement, ok3 := interpreter.Model.GetString(argReplacement)
	kCount, ok4 := interpreter.Model.GetBigIntObject(argCount)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return invalidArgsResult()
	}
	count, countOk := kCount.ToInt32()
	if !countOk {
		return invalidArgsResult()
	}

	result := strings.Replace(kS, kToReplace, kReplacement, count)
	return interpreter.Model.NewString(result), nil
}

func (stringHooksType) replaceAll(argS m.KReference, argToReplace m.KReference, argReplacement m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	kS, ok1 := interpreter.Model.GetString(argS)
	kToReplace, ok2 := interpreter.Model.GetString(argToReplace)
	kReplacement, ok3 := interpreter.Model.GetString(argReplacement)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if len(kS) == 0 {
		return argS, nil // empty
	}
	count := strings.Count(kS, kToReplace)
	if count == 0 {
		return argS, nil
	}
	if count == 1 && strings.HasPrefix(kS, kToReplace) {
		if len(kReplacement) == 0 {
			// just cut off the prefix
			return interpreter.Model.NewString(kS[len(kToReplace):]), nil
		}
		return interpreter.Model.NewString(kReplacement + kS[len(kToReplace):]), nil
	}

	result := strings.ReplaceAll(kS, kToReplace, kReplacement)
	return interpreter.Model.NewString(result), nil
}

func (stringHooksType) replaceFirst(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) countAllOccurrences(argS m.KReference, argToCount m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	kS, ok1 := interpreter.Model.GetString(argS)
	kToCount, ok2 := interpreter.Model.GetString(argToCount)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}

	result := strings.Count(kS, kToCount)
	return interpreter.Model.FromInt(result), nil
}

func (stringHooksType) category(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) directionality(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) findChar(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (stringHooksType) rfindChar(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}
