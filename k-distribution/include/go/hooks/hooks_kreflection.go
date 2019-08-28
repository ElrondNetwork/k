%COMMENT%

package %PACKAGE%

import (
	m "%INCLUDE_MODEL%"
)

type kreflectionHooksType int

const kreflectionHooks kreflectionHooksType = 0

var constKReflectionSortInt = m.NewStringConstant("Int")
var constKReflectionSortString = m.NewStringConstant("String")
var constKReflectionSortBytes = m.NewStringConstant("Bytes")
var constKReflectionSortBool = m.NewStringConstant("Bool")

func (kreflectionHooksType) sort(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	if obj, t := interpreter.Model.GetKTokenObject(c); t {
		return interpreter.Model.NewString(obj.Sort.Name()), nil
	}
	if m.IsInt(c) {
		return constKReflectionSortInt, nil
	}
	if m.IsString(c) {
		return constKReflectionSortString, nil
	}
	if m.IsBytes(c) {
		return constKReflectionSortBytes, nil
	}
	if m.IsBool(c) {
		return constKReflectionSortBool, nil
	}
	if sortName, ok := interpreter.Model.CollectionSortName(c); ok {
		return interpreter.Model.NewString(sortName), nil
	}

	return m.NoResult, m.GetHookNotImplementedError()

}

func (kreflectionHooksType) getKLabel(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	if k, t := interpreter.Model.GetKApplyObject(c); t {
		return interpreter.Model.NewInjectedKLabel(k.Label), nil
	}
	return m.InternedBottom, nil
}

func (kreflectionHooksType) configuration(lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return config, nil
}

var freshCounter int

func (kreflectionHooksType) fresh(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	if str, t := interpreter.Model.GetString(c); t {
		sort := m.ParseSort(str)
		result, err := interpreter.freshFunction(sort, config, freshCounter)
		if err != nil {
			return m.NoResult, err
		}
		freshCounter++
		return result, nil
	}
	return m.NoResult, m.GetHookNotImplementedError()
}

func (kreflectionHooksType) isConcrete(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.BoolTrue, nil
}

func (kreflectionHooksType) getenv(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, m.GetHookNotImplementedError()
}

func (kreflectionHooksType) argv(lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, m.GetHookNotImplementedError()
}
