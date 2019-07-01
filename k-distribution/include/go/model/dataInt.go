%COMMENT%

package %PACKAGE_MODEL%

import (
	"fmt"
	"math"
	"math/big"
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
