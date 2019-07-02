package %PACKAGE_MODEL%

import (
	"math/big"

	"github.com/JohnCGriffin/overflow"
)

var bigOne = big.NewInt(1)

// helper function for writing operations in fewer lines
func (*ModelState) bothSmall(ref1 KReference, ref2 KReference) (int32, int32, bool) {
	small1, isSmall1 := getSmallInt(ref1)
	if !isSmall1 {
		return 0, 0, false
	}
	small2, isSmall2 := getSmallInt(ref2)
	if !isSmall2 {
		return 0, 0, false
	}
	return small1, small2, true
}

// helper function for writing operations in fewer lines
func (ms *ModelState) bothBig(ref1 KReference, ref2 KReference) (*big.Int, *big.Int, bool) {
	big1, isInt1 := ms.GetBigInt(ref1)
	if !isInt1 {
		return nil, nil, false
	}
	big2, isInt2 := ms.GetBigInt(ref2)
	if !isInt2 {
		return nil, nil, false
	}
	return big1, big2, true
}

// IntEquals returns ref1 == ref2, if types ok
func (ms *ModelState) IntEquals(ref1 KReference, ref2 KReference) (bool, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		return small1 == small2, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		return big1.Cmp(big2) == 0, true
	}

	return false, false
}

// IntLt returns ref1 < ref2, if types ok
func (ms *ModelState) IntLt(ref1 KReference, ref2 KReference) (bool, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		return small1 < small2, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		return big1.Cmp(big2) < 0, true
	}

	return false, false
}

// IntLe returns ref1 <= ref2, if types ok
func (ms *ModelState) IntLe(ref1 KReference, ref2 KReference) (bool, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		return small1 <= small2, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		return big1.Cmp(big2) <= 0, true
	}

	return false, false
}

// IntGt returns ref1 > ref2, if types ok
func (ms *ModelState) IntGt(ref1 KReference, ref2 KReference) (bool, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		return small1 > small2, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		return big1.Cmp(big2) > 0, true
	}

	return false, false
}

// IntGe returns ref1 >= ref2, if types ok
func (ms *ModelState) IntGe(ref1 KReference, ref2 KReference) (bool, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		return small1 >= small2, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		return big1.Cmp(big2) >= 0, true
	}

	return false, false
}

// IntAdd returns ref1 + ref2, if types ok
func (ms *ModelState) IntAdd(ref1 KReference, ref2 KReference) (KReference, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		result, of := overflow.Add32(small1, small2)
		if !of && fitsInSmallIntReference(result) {
			return smallIntReference(result), true
		}
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		if big1.Sign() == 0 {
			return ref2, true
		}
		if big2.Sign() == 0 {
			return ref1, true
		}

		var z big.Int
		z.Add(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntSub returns ref1 - ref2, if types ok
func (ms *ModelState) IntSub(ref1 KReference, ref2 KReference) (KReference, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		result, of := overflow.Sub32(small1, small2)
		if !of && fitsInSmallIntReference(result) {
			return smallIntReference(result), true
		}
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		if big2.Sign() == 0 {
			return ref1, true
		}

		var z big.Int
		z.Sub(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntMul returns ref1 x ref2, if types ok
func (ms *ModelState) IntMul(ref1 KReference, ref2 KReference) (KReference, bool) {
	small1, small2, smallOk := ms.bothSmall(ref1, ref2)
	if smallOk {
		result, of := overflow.Mul32(small1, small2)
		if !of && fitsInSmallIntReference(result) {
			return smallIntReference(result), true
		}
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		if big1.Sign() == 0 || big2.Sign() == 0 {
			return IntZero, true
		}
		if big1.Cmp(bigOne) == 0 {
			return ref2, true
		}
		if big2.Cmp(bigOne) == 0 {
			return ref1, true
		}

		var z big.Int
		z.Mul(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntDiv performs integer division.
// The result is truncated towards zero and obeys the rule of signs.
func (ms *ModelState) IntDiv(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		resultPositive := true
		if big1.Sign() < 0 {
			resultPositive = !resultPositive
			big1 = big.NewInt(0).Neg(big1)
		}
		if big2.Sign() < 0 {
			resultPositive = !resultPositive
			big2 = big.NewInt(0).Neg(big2)
		}

		var z big.Int
		z.Div(big1, big2)
		if !resultPositive {
			z.Neg(&z)
		}
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntMod performs integer remainder.
// The result of rem a b has the sign of a, and its absolute value is strictly smaller than the absolute value of b.
// The result satisfies the equality a = b * div a b + rem a b.
func (ms *ModelState) IntMod(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		arg1Negative := false
		if big1.Sign() < 0 {
			arg1Negative = true
			big1 = big.NewInt(0).Neg(big1)
		}
		if big2.Sign() < 0 {
			big2 = big.NewInt(0).Neg(big2)
		}

		var z big.Int
		z.Mod(big1, big2)
		if arg1Negative {
			z.Neg(&z)
		}
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntEuclidianDiv performs Euclidian division.
func (ms *ModelState) IntEuclidianDiv(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.Div(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntEuclidianMod performs Euclidian remainder.
func (ms *ModelState) IntEuclidianMod(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.Mod(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntPow returns ref1 ^ ref2, if types ok
func (ms *ModelState) IntPow(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.Exp(big1, big2, nil)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntPowMod returns (ref1 ^ ref2) mod ref3, if types ok
func (ms *ModelState) IntPowMod(ref1 KReference, ref2 KReference, ref3 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		big3, big3Ok := ms.GetBigInt(ref3)
		if big3Ok {
			var z big.Int
			z.Exp(big1, big2, big3)
			return ms.FromBigInt(&z), true
		}
	}

	return NullReference, false
}

// IntShl returns ref1 << ref2, if types ok
func (ms *ModelState) IntShl(ref1 KReference, ref2 KReference) (KReference, bool) {
	arg2, arg2Ok := ms.GetUint(ref2)
	if !arg2Ok {
		return NullReference, false
	}

	arg1, arg1Ok := ms.GetBigInt(ref1)
	if arg1Ok {
		var z big.Int
		z.Lsh(arg1, arg2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntShr returns ref1 >> ref2, if types ok
func (ms *ModelState) IntShr(ref1 KReference, ref2 KReference) (KReference, bool) {
	arg2, arg2Ok := ms.GetUint(ref2)
	if !arg2Ok {
		return NullReference, false
	}

	arg1, arg1Ok := ms.GetBigInt(ref1)
	if arg1Ok {
		var z big.Int
		z.Rsh(arg1, arg2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntAnd returns bitwise and, ref1 & ref2, if types ok
func (ms *ModelState) IntAnd(ref1 KReference, ref2 KReference) (KReference, bool) {
	if ms.IsZero(ref1) || ms.IsZero(ref2) {
		return IntZero, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.And(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntOr returns bitwise or, ref1 | ref2, if types ok
func (ms *ModelState) IntOr(ref1 KReference, ref2 KReference) (KReference, bool) {
	if ms.IsZero(ref1) {
		return ref2, true
	}
	if ms.IsZero(ref2) {
		return ref1, true
	}

	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.Or(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntXor returns bitwise xor, ref1 xor ref2, if types ok
func (ms *ModelState) IntXor(ref1 KReference, ref2 KReference) (KReference, bool) {
	big1, big2, bigOk := ms.bothBig(ref1, ref2)
	if bigOk {
		var z big.Int
		z.Xor(big1, big2)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntNot returns bitwise not, if type ok
func (ms *ModelState) IntNot(ref KReference) (KReference, bool) {
	arg, argOk := ms.GetBigInt(ref)
	if argOk {
		var z big.Int
		z.Not(arg)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntAbs returns the absoute value, if type ok
func (ms *ModelState) IntAbs(ref KReference) (KReference, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small >= 0 {
			return ref, true
		}
		return smallIntReference(-small), true
	}

	bigArg, bigOk := ms.GetBigInt(ref)
	if bigOk {
		if bigArg.Sign() >= 0 {
			return ref, true
		}
		var z big.Int
		z.Neg(bigArg)
		return ms.FromBigInt(&z), true
	}

	return NullReference, false
}

// IntLog2 basically counts the number of bits after the most significant bit.
// It is equal to a a truncated log2 of the number.
// Argument must be strictly positive.
func (ms *ModelState) IntLog2(ref KReference) (KReference, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small <= 0 {
			return NullReference, false
		}
		nrBits := 0
		for small > 0 {
			small = small >> 1
			nrBits++
		}
		return ms.FromInt(nrBits - 1), true
	}

	bigArg, bigOk := ms.GetBigInt(ref)
	if bigOk {
		if bigArg.Sign() <= 0 {
			return NullReference, false
		}
		bytes := bigArg.Bytes()
		leadingByte := bytes[0]
		nrLeadingBits := 0
		for leadingByte > 0 {
			leadingByte = leadingByte >> 1
			nrLeadingBits++
		}
		return ms.FromInt(nrLeadingBits + (len(bytes)-1)*8 - 1), true
	}

	return NullReference, false
}
