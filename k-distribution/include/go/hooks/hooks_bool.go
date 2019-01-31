package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type boolHooksType int

const boolHooks boolHooksType = 0

func (boolHooksType) and(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool(b1 && b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (h boolHooksType) andThen(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return h.and(c1, c2, lbl, sort, config)
}

func (boolHooksType) or(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool(b1 || b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (h boolHooksType) orElse(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return h.or(c1, c2, lbl, sort, config)
}

func (boolHooksType) not(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b, ok := c.(m.Bool)
	if ok {
		return m.Bool(!b), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) implies(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool((!b1) || b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) ne(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool(b1 != b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) eq(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool(b1 == b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) xor(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	b1, ok1 := c1.(m.Bool)
	b2, ok2 := c2.(m.Bool)
	if ok1 && ok2 {
		return m.Bool(b1 != b2), nil
	}
	return noResult, &hookNotImplementedError{}
}
