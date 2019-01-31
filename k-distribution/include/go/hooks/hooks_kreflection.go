package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type kreflectionHooksType int

const kreflectionHooks kreflectionHooksType = 0

func (kreflectionHooksType) sort(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	switch k := c.(type) {
	case m.KToken:
		return m.String(k.Sort.Name()), nil
	case m.Int:
		return m.String("Int"), nil
	case m.String:
		return m.String("String"), nil
	case m.Bytes:
		return m.String("Bytes"), nil
	case m.Bool:
		return m.String("Bool"), nil
	case m.Map:
		return m.String(k.Sort.Name()), nil
	case m.List:
		return m.String(k.Sort.Name()), nil
	case m.Set:
		return m.String(k.Sort.Name()), nil
	default:
		return noResult, &hookNotImplementedError{}
	}
}

func (kreflectionHooksType) getKLabel(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	if k, t := c.(m.KApply); t {
		return m.InjectedKLabel{Label: k.Label}, nil
	}
	return internedBottom, nil
}

func (kreflectionHooksType) configuration(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return config, nil
}

func (kreflectionHooksType) fresh(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	if k, t := c.(m.String); t {
		sort := m.ParseSort(string(k))
		result, err := freshFunction(sort, config, freshCounter)
		if err != nil {
		    return noResult, err
		}
		freshCounter++
		return result, nil
	}
	return noResult, &hookNotImplementedError{}
}

func (kreflectionHooksType) isConcrete(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.Bool(true), nil
}

func (kreflectionHooksType) getenv(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}

func (kreflectionHooksType) argv(lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return noResult, &hookNotImplementedError{}
}
