package %PACKAGE_INTERPRETER%

type intHooksType int

const intHooks intHooksType = 0

func (intHooksType) tmod(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) emod(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) add(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 + i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) eq(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 == i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) ne(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 != i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) and(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 & i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) mul(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 * i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) sub(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 - i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) tdiv(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) ediv(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) shl(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
    i2, ok2 := c2.(Int)
    if ok1 && ok2 {
        return Int(i1 << uint32(i2)), nil
    }
    return noResult, &hookNotImplementedError{}
}

func (intHooksType) shr(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 >> uint32(i2)), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) le(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 <= i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) lt(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 < i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) ge(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 >= i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) gt(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 > i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) pow(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) powmod(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) xor(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) or(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 | i2), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) not(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) abs(c K, lbl KLabel, sort Sort, config K) (K, error) {
	i, ok := c.(Int)
	if ok {
		if i < 0 {
			return Int(-i), nil
		}
		return Int(i), nil
	}
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) max(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) min(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) log2(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) bitRange(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) signExtendBitRange(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) rand(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (intHooksType) srand(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}
