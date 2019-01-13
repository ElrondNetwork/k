package main

type intHooksType int

const intHooks intHooksType = 0

func (intHooksType) tmod(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) emod(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) add(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 + i2)
	}
	panic("Not implemented")
}

func (intHooksType) le(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 + i2)
	}
	panic("Not implemented")
}

func (intHooksType) eq(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 == i2)
	}
	panic("Not implemented")
}

func (intHooksType) ne(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 != i2)
	}
	panic("Not implemented")
}

func (intHooksType) and(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 & i2)
	}
	panic("Not implemented")
}

func (intHooksType) mul(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 * i2)
	}
	panic("Not implemented")
}

func (intHooksType) sub(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 - i2)
	}
	panic("Not implemented")
}

func (intHooksType) tdiv(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) ediv(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) shl(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
    i2, ok2 := c2.(Int)
    if ok1 && ok2 {
        return Int(i1 << uint32(i2))
    }
    panic("Not implemented")
}

func (intHooksType) shr(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 >> uint32(i2))
	}
	panic("Not implemented")
}

func (intHooksType) lt(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 < i2)
	}
	panic("Not implemented")
}

func (intHooksType) ge(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 >= i2)
	}
	panic("Not implemented")
}

func (intHooksType) gt(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Bool(i1 > i2)
	}
	panic("Not implemented")
}

func (intHooksType) pow(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) powmod(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) xor(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) or(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	i1, ok1 := c1.(Int)
	i2, ok2 := c2.(Int)
	if ok1 && ok2 {
		return Int(i1 | i2)
	}
	panic("Not implemented")
}

func (intHooksType) not(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) abs(c K, lbl KLabel, sort Sort, config K) K {
	i, ok := c.(Int)
	if ok {
		if i < 0 {
			return Int(-i)
		}
		return Int(i)
	}
	panic("Not implemented")
}

func (intHooksType) max(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) min(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) log2(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) bitRange(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) signExtendBitRange(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) rand(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (intHooksType) srand(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}
