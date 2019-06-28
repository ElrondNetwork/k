%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"math/big"
)

type intHooksType int

const intHooks intHooksType = 0

var bigIntZero = big.NewInt(0)
var bigIntOne = big.NewInt(1)

func (intHooksType) eq(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) == 0), nil
}

func (intHooksType) ne(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) != 0), nil
}

func (intHooksType) le(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) <= 0), nil
}

func (intHooksType) lt(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) < 0), nil
}

func (intHooksType) ge(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) >= 0), nil
}

func (intHooksType) gt(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	return m.ToKBool(i1.Value.Cmp(i2.Value) > 0), nil
}

func (intHooksType) add(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Add(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) sub(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Sub(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) mul(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Mul(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

// Integer division. The result is truncated towards zero and obeys the rule of signs.
func (t intHooksType) tdiv(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i2.IsZero() {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	resultPositive := true
	if i1.IsNegative() {
		resultPositive = !resultPositive
	}
	if i2.IsNegative() {
		resultPositive = !resultPositive
	}
	var i1Abs, i2Abs, z big.Int
	i1Abs.Abs(i1.Value)
	i2Abs.Abs(i2.Value)

	z.Div(&i1Abs, &i2Abs)
	if !resultPositive {
		z.Neg(&z)
	}
	return interpreter.Model.FromBigInt(&z), nil
}

// Integer remainder. The result of rem a b has the sign of a, and its absolute value is strictly smaller than the absolute value of b.
// The result satisfies the equality a = b * div a b + rem a b.
func (t intHooksType) tmod(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i2.IsZero() {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	var i1Abs, i2Abs, z big.Int
	i1Abs.Abs(i1.Value)
	i2Abs.Abs(i2.Value)

	z.Mod(&i1Abs, &i2Abs)
	if i1.IsNegative() {
		z.Neg(&z)
	}
	return interpreter.Model.FromBigInt(&z), nil
}

// Euclidian division
func (intHooksType) ediv(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i2.IsZero() {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	var z big.Int
	z.Div(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

// Euclidian remainder
func (intHooksType) emod(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i2.Value.Sign() == 0 {
		return m.NoResult, &hookDivisionByZeroError{}
	}
	var z big.Int
	z.Mod(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) pow(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Exp(i1.Value, i2.Value, nil)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) powmod(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	i3, ok3 := interpreter.Model.GetBigIntObject(c3)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Exp(i1.Value, i2.Value, i3.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) shl(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	arg2, arg2Ok := i2.ToUint32()
	if !arg2Ok {
		return invalidArgsResult()
	}
	var z big.Int
	z.Lsh(i1.Value, arg2)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) shr(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	arg2, arg2Ok := i2.ToUint32()
	if !arg2Ok {
		return invalidArgsResult()
	}
	var z big.Int
	z.Rsh(i1.Value, arg2)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) and(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.And(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) or(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Or(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) xor(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	var z big.Int
	z.Xor(i1.Value, i2.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) not(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i, ok := interpreter.Model.GetBigIntObject(c)
	if !ok {
		return invalidArgsResult()
	}
	var z big.Int
	z.Not(i.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) abs(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i, ok := interpreter.Model.GetBigIntObject(c)
	if !ok {
		return invalidArgsResult()
	}
	var z big.Int
	z.Abs(i.Value)
	return interpreter.Model.FromBigInt(&z), nil
}

func (intHooksType) max(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i1.Value.Cmp(i2.Value) >= 0 {
		return c1, nil
	}
	return c2, nil
}

func (intHooksType) min(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	i1, ok1 := interpreter.Model.GetBigIntObject(c1)
	i2, ok2 := interpreter.Model.GetBigIntObject(c2)
	if !ok1 || !ok2 {
		return invalidArgsResult()
	}
	if i1.Value.Cmp(i2.Value) >= 0 {
		return c2, nil
	}
	return c1, nil
}

func (intHooksType) log2(karg m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	bigi, ok := interpreter.Model.GetBigIntObject(karg)
	if !ok {
		return invalidArgsResult()
	}
	if bigi.Value.Sign() <= 0 {
		return invalidArgsResult()
	}
	bytes := bigi.Value.Bytes()
	leadingByte := bytes[0]
	nrBytes := 0
	for leadingByte > 0 {
		leadingByte = leadingByte >> 1
		nrBytes++
	}
	return interpreter.Model.FromInt(nrBytes + (len(bytes)-1)*8 - 1), nil
}

func (intHooksType) bitRange(argI m.KReference, argOffset m.KReference, argLen m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	// rule bitRangeInt(I::Int, IDX::Int, LEN::Int) => (I >>Int IDX) modInt (1 <<Int LEN)
	ki, ok1 := interpreter.Model.GetBigIntObject(argI)
	koff, ok2 := interpreter.Model.GetBigIntObject(argOffset)
	klen, ok3 := interpreter.Model.GetBigIntObject(argLen)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if ki.IsZero() {
		return m.IntZero, nil // any operation on zero will result in zero
	}

	if koff.IsNegative() {
		return invalidArgsResult()
	}
	offset, offsetOk := koff.ToInt32()
	if !offsetOk {
		if ki.IsPositive() {
			// means it doesn't fit in an int32, so a huge number
			// huge offset means that certainly no 1 bits will be caught
			// scenario occurs in tests/VMTests/vmIOandFlowOperations/byte1/byte1.iele.json
			// but only if the number is positive, otherwise the result would be a ridiculously large number of 1's
			return m.IntZero, nil
		}
		return invalidArgsResult()
	}

	length, lengthOk := klen.ToPositiveInt32()
	if !lengthOk {
		return invalidArgsResult()
	}
	if length == 0 {
		return m.IntZero, nil
	}
	if offset&7 != 0 || length&7 != 0 {
		// this is a quick check that they are both divisible by 8
		// as long as they are divisible by 8, we can operate on whole bytes
		// if they are not, things get more complicated, will only implement when necessary
		return m.NoResult, &hookNotImplementedError{}
	}
	offsetBytes := offset >> 3 // divide by 8 to get number of bytes
	lengthBytes := length >> 3 // divide by 8 to get number of bytes

	resultBytes := m.BigIntToTwosComplementBytes(ki.Value, lengthBytes+offsetBytes)
	if offsetBytes != 0 {
		resultBytes = resultBytes[0:lengthBytes]
	}

	result := new(big.Int)
	result.SetBytes(resultBytes)
	return interpreter.Model.FromBigInt(result), nil
}

func (intHooksType) signExtendBitRange(argI m.KReference, argOffset m.KReference, argLen m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	// rule signExtendBitRangeInt(I::Int, IDX::Int, LEN::Int) => (bitRangeInt(I, IDX, LEN) +Int (1 <<Int (LEN -Int 1))) modInt (1 <<Int LEN) -Int (1 <<Int (LEN -Int 1))
	ki, ok1 := interpreter.Model.GetBigIntObject(argI)
	koff, ok2 := interpreter.Model.GetBigIntObject(argOffset)
	klen, ok3 := interpreter.Model.GetBigIntObject(argLen)
	if !ok1 || !ok2 || !ok3 {
		return invalidArgsResult()
	}
	if ki.IsZero() {
		return m.IntZero, nil // any operation on zero will result in zero
	}

	if koff.IsNegative() {
		return invalidArgsResult()
	}
	offset, offsetOk := koff.ToInt32()
	if !offsetOk {
		if ki.IsPositive() {
			// means it doesn't fit in an int32, so a huge number
			// huge offset means that certainly no 1 bits will be caught
			// scenario occurs in tests/VMTests/vmIOandFlowOperations/byte1/byte1.iele.json
			// but only if the number is positive, otherwise the result would be a ridiculously large number of 1's
			return m.IntZero, nil
		}
		return invalidArgsResult()
	}

	length, lengthOk := klen.ToPositiveInt32()
	if !lengthOk {
		return invalidArgsResult()
	}
	if length == 0 {
		return m.IntZero, nil
	}
	if offset&7 != 0 || length&7 != 0 {
		// this is a quick check that they are both divisible by 8
		// as long as they are divisible by 8, we can operate on whole bytes
		// if they are not, things get more complicated, will only implement when necessary
		return m.NoResult, &hookNotImplementedError{}
	}
	offsetBytes := offset >> 3 // divide by 8 to get number of bytes
	lengthBytes := length >> 3 // divide by 8 to get number of bytes

	resultBytes := m.BigIntToTwosComplementBytes(ki.Value, lengthBytes+offsetBytes)
	if offsetBytes != 0 {
		resultBytes = resultBytes[0:lengthBytes]
	}

	result := m.TwosComplementBytesToBigInt(resultBytes)
	return interpreter.Model.FromBigInt(result), nil
}

func (intHooksType) rand(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (intHooksType) srand(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}
