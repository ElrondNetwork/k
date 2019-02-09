package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type intHooksType int

const intHooks intHooksType = 0

func (intHooksType) tmod(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) emod(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) add(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 + i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) eq(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 == i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) ne(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 != i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) and(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 & i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) mul(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 * i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) sub(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 - i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) tdiv(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
		i1, ok1 := c1.(m.Int)
    	i2, ok2 := c2.(m.Int)
    	if ok1 && ok2 {
    		return m.Int(i1 / i2), nil
    	}
    	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) ediv(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
			i1, ok1 := c1.(m.Int)
        	i2, ok2 := c2.(m.Int)
        	if ok1 && ok2 {
        		return m.Int(i1 % i2), nil
        	}
        	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) shl(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
    i2, ok2 := c2.(m.Int)
    if ok1 && ok2 {
        return m.Int(i1 << uint32(i2)), nil
    }
    return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) shr(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 >> uint32(i2)), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) le(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 <= i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) lt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 < i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) ge(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 >= i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) gt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Bool(i1 > i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) pow(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) powmod(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) xor(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) or(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if ok1 && ok2 {
		return m.Int(i1 | i2), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) not(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) abs(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i, ok := c.(m.Int)
	if ok {
		if i < 0 {
			return m.Int(-i), nil
		}
		return m.Int(i), nil
	}
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) max(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) min(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) log2(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) bitRange(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) signExtendBitRange(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) rand(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) srand(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}
