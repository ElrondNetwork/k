%COMMENT%

package %PACKAGE_MODEL%

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const maxSmallInt = math.MaxInt32
const minSmallInt = math.MinInt32

var maxSmallIntAsBigInt = big.NewInt(maxSmallInt)
var minSmallIntAsBigInt = big.NewInt(minSmallInt)

// only attempt to parse as small int strings shorter than this
var maxSmallIntStringLength = len(fmt.Sprintf("%d", maxSmallIntAsBigInt)) - 2

// BigInt is a KObject representing a big int in K
type BigInt struct {
	Value *big.Int
}

func (*BigInt) referenceType() kreferenceType {
	return bigIntRef
}

func fitsInSmallIntReference(i int32) bool {
	return i > minSmallInt && i < maxSmallInt
}

func smallIntReference(i int32) KReference {
	if i < 0 {
		return KReference{refType: smallNegativeIntRef, value1: uint32(-i)}
	}
	return KReference{refType: smallPositiveIntRef, value1: uint32(i)}
}

func getSmallInt(ref KReference) (int32, bool) {
	if ref.refType == smallPositiveIntRef {
		if ref.value1 > maxSmallInt {
			return 0, false
		}
		return int32(ref.value1), true
	}
	if ref.refType == smallNegativeIntRef {
		if ref.value1 > -minSmallInt {
			return 0, false
		}
		return -int32(ref.value1), true
	}
	return 0, false
}

func (ms *ModelState) getBigIntObject(ref KReference) (*BigInt, bool) {
	if ref.refType != bigIntRef {
		return nil, false
	}
	obj := ms.getReferencedObject(ref)
	castObj, typeOk := obj.(*BigInt)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// IsInt returns true if reference points to an integer
func IsInt(ref KReference) bool {
	return ref.refType == smallPositiveIntRef || ref.refType == smallNegativeIntRef || ref.refType == bigIntRef
}

// IntZero is a reference to the constant integer 0
var IntZero = smallIntReference(0)

// IntOne is a reference to the constant integer 1
var IntOne = smallIntReference(1)

// IntMinusOne is a reference to the constant integer -1
var IntMinusOne = smallIntReference(-1)

// FromBigInt provides a reference to an integer (big or small)
func (ms *ModelState) FromBigInt(bi *big.Int) KReference {
	// attempt to make it small
	if bi.IsInt64() {
		biInt64 := bi.Int64()
		if biInt64 >= minSmallInt && biInt64 <= maxSmallInt {
			return smallIntReference(int32(biInt64))
		}
	}
	// make it big
	return ms.addObject(&BigInt{Value: bi})
}

// NewIntConstant creates a new integer constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewIntConstant(stringRepresentation string) KReference {
	ref := constantsModel.IntFromString(stringRepresentation)
	ref.constantObject = true
	return ref
}

// FromInt converts a Go integer to an integer in the model
func (ms *ModelState) FromInt(x int) KReference {
	if x >= minSmallInt && x <= maxSmallInt {
		return smallIntReference(int32(x))
	}
	return ms.addObject(&BigInt{Value: big.NewInt(int64(x))})
}

// FromInt64 converts a int64 to an integer in the model
func (ms *ModelState) FromInt64(x int64) KReference {
	if x >= minSmallInt && x <= maxSmallInt {
		return smallIntReference(int32(x))
	}
	return ms.addObject(&BigInt{Value: big.NewInt(x)})
}

// FromUint64 converts a uint64 to an integer in the model
func (ms *ModelState) FromUint64(x uint64) KReference {
	if x <= maxSmallInt {
		return smallIntReference(int32(x))
	}
	var z big.Int
	z.SetUint64(x)
	return ms.addObject(&BigInt{Value: &z})
}

// IntFromByte converts a byte to an integer in the model
func (ms *ModelState) IntFromByte(x byte) KReference {
	return smallIntReference(int32(x))
}

// IntFromBytes converts a byte array to an integer in the model
func (ms *ModelState) IntFromBytes(bytes []byte) KReference {
	z := big.NewInt(0)
	z.SetBytes(bytes)
	return ms.FromBigInt(z)
}

// ParseInt creates K int from string representation
func (ms *ModelState) ParseInt(str string) (KReference, error) {
	if str == "0" {
		return IntZero, nil
	}
	if len(str) < maxSmallIntStringLength {
		i, err := strconv.Atoi(str)
		if err != nil {
			return NullReference, &parseIntError{parseVal: str}
		}
		return smallIntReference(int32(i)), nil
	}

	b := big.NewInt(0)
	b.UnmarshalText([]byte(str))
	if b.Sign() == 0 {
		return IntZero, &parseIntError{parseVal: str}
	}
	return ms.FromBigInt(b), nil
}

// ParseIntFromBase creates K int from string representation in a given base
func (ms *ModelState) ParseIntFromBase(str string, base int) (KReference, error) {
	if base == 10 {
		return ms.ParseInt(str)
	}
	if str == "0" {
		return IntZero, nil
	}
	b := big.NewInt(0)
	_, ok := b.SetString(str, base)
	if !ok {
		return IntZero, &parseIntError{parseVal: str}
	}
	return ms.FromBigInt(b), nil
}

// IntFromString does the same as ParseInt but panics instead of returning an error
func (ms *ModelState) IntFromString(s string) KReference {
	i, err := ms.ParseInt(s)
	if err != nil {
		panic(err)
	}
	return i
}

// GetBigInt yields a big.Int cast from any K integer object, if possible.
func (ms *ModelState) GetBigInt(ref KReference) (*big.Int, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		return big.NewInt(int64(small)), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if !isBigInt {
		return nil, false
	}
	return bi.Value, true
}

// IsZero returns true if an item represents number 0
func (ms *ModelState) IsZero(ref KReference) bool {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		return small == 0
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if !isBigInt {
		return bi.Value.Sign() == 0
	}

	return false
}

// GetUint converts to uint if possible, returns (0, false) if not
func (ms *ModelState) GetUint(ref KReference) (uint, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small < 0 {
			return 0, false
		}
		return uint(small), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		if !bi.Value.IsUint64() {
			return 0, false
		}
		u64 := bi.Value.Uint64()
		if u64 > math.MaxUint32 {
			return 0, false
		}
		return uint(u64), true
	}

	return 0, false
}

// GetUint64 converts to uint64 if possible, returns (0, false) if not
func (ms *ModelState) GetUint64(ref KReference) (uint64, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small < 0 {
			return 0, false
		}
		return uint64(small), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		if !bi.Value.IsUint64() {
			return 0, false
		}
		return bi.Value.Uint64(), true
	}

	return 0, false
}

// GetInt converts to int if possible, returns (0, false) if not
func (ms *ModelState) GetInt(ref KReference) (int, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		return int(small), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		if !bi.Value.IsInt64() {
			return 0, false
		}
		i64 := bi.Value.Int64()
		if i64 >= math.MinInt32 && i64 <= math.MaxInt32 {
			return int(i64), true
		}
	}

	return 0, false
}

// GetPositiveInt converts to int32 if possible, returns (0, false) if not.
// Also rejects negative numbers, so we don't have to test for that again.
func (ms *ModelState) GetPositiveInt(ref KReference) (int, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small < 0 {
			return 0, false
		}
		return int(small), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		if !bi.Value.IsUint64() {
			return 0, false
		}
		u64 := bi.Value.Uint64()
		if u64 > math.MaxUint32 {
			return 0, false
		}
		return int(u64), true
	}

	return 0, false
}

// GetByte converts to 1 byte if possible, returns (0, false) if not
func (ms *ModelState) GetByte(ref KReference) (byte, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		if small < 0 || small > 255 {
			return 0, false
		}
		return byte(small), true
	}

	bi, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		if !bi.Value.IsUint64() {
			return 0, false
		}
		u64 := bi.Value.Uint64()
		if u64 > 255 {
			return 0, false
		}
		return byte(u64), true
	}

	return 0, false
}

// GetIntAsDecimalString converts a K integer to a decimal string representation, decimal, if possible.
func (ms *ModelState) GetIntAsDecimalString(ref KReference) (string, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		return fmt.Sprintf("%d", small), true
	}

	bigI, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		return bigI.Value.String(), true
	}

	return "", false
}

// GetIntToString converts a K integer to a string representation in given base, if possible.
func (ms *ModelState) GetIntToString(ref KReference, base int) (string, bool) {
	small, isSmall := getSmallInt(ref)
	if isSmall {
		return strconv.FormatInt(int64(small), base), true
	}

	bigI, isBigInt := ms.getBigIntObject(ref)
	if isBigInt {
		return bigI.Value.Text(base), true
	}

	return "", false
}
