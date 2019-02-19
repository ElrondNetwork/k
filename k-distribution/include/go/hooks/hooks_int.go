package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"math/big"
)

type intHooksType int

const intHooks intHooksType = 0

func (intHooksType) eq(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) == 0), nil
}

func (intHooksType) ne(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) != 0), nil
}

func (intHooksType) le(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) <= 0), nil
}

func (intHooksType) lt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) < 0), nil
}

func (intHooksType) ge(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) >= 0), nil
}

func (intHooksType) gt(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	return m.Bool(i1.Value.Cmp(i2.Value) > 0), nil
}

func (intHooksType) add(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Add(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) sub(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Sub(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) mul(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Mul(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (t intHooksType) tdiv(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	// TODO: investigate if we need another implementation here
	return t.ediv(c1, c2, lbl, sort, config)
}

func (t intHooksType) tmod(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	// TODO: investigate if we need another implementation here
	return t.emod(c1, c2, lbl, sort, config)
}

// euclidian division
func (intHooksType) ediv(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if i2.Value.Cmp(m.IntZero.Value) == 0 {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	var z big.Int
	z.Div(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) emod(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if i2.Value.Cmp(m.IntZero.Value) == 0 {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	var z big.Int
	z.Mod(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) pow(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) powmod(c1 m.K, c2 m.K, c3 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) shl(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !i2.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Lsh(i1.Value, uint(i2.Value.Uint64()))
	return m.NewInt(&z), nil
}

func (intHooksType) shr(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if !i2.Value.IsUint64() {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Rsh(i1.Value, uint(i2.Value.Uint64()))
	return m.NewInt(&z), nil
}

func (intHooksType) and(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.And(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) or(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Or(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) xor(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Xor(i1.Value, i2.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) not(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i, ok := c.(m.Int)
	if !ok {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Not(i.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) abs(c m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i, ok := c.(m.Int)
	if !ok {
		return m.NoResult, &hookInvalidArgsError{}
	}
	var z big.Int
	z.Abs(i.Value)
	return m.NewInt(&z), nil
}

func (intHooksType) max(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if i1.Value.Cmp(i2.Value) >= 0 {
		return c1, nil
	}
	return c2, nil
}

func (intHooksType) min(c1 m.K, c2 m.K, lbl m.KLabel, sort m.Sort, config m.K) (m.K, error) {
	i1, ok1 := c1.(m.Int)
	i2, ok2 := c2.(m.Int)
	if !ok1 || !ok2 {
		return m.NoResult, &hookInvalidArgsError{}
	}
	if i1.Value.Cmp(i2.Value) >= 0 {
		return c2, nil
	}
	return c1, nil
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
