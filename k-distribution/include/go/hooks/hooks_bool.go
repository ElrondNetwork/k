package main

type boolHooksType int

const boolHooks boolHooksType = 0

func (boolHooksType) and(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool(b1 && b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (h boolHooksType) andThen(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return h.and(c1, c2, lbl, sort, config)
}

func (boolHooksType) or(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool(b1 || b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (h boolHooksType) orElse(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return h.or(c1, c2, lbl, sort, config)
}

func (boolHooksType) not(c K, lbl KLabel, sort Sort, config K) (K, error) {
	b, ok := c.(Bool)
	if ok {
		return Bool(!b), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) implies(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool((!b1) || b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) ne(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool(b1 != b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) eq(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool(b1 == b2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (boolHooksType) xor(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	b1, ok1 := c1.(Bool)
	b2, ok2 := c2.(Bool)
	if ok1 && ok2 {
		return Bool(b1 != b2), nil
	}
	return noResult, &hookNotImplementedError{}
}
