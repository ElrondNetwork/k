%COMMENT%

package %PACKAGE_MODEL%

import (
	"math/big"
)

// no small ints in this version
var maxSmallIntStringLength = -1

// BigInt is a KObject representing a big int in K
type BigInt struct {
	Value *big.Int
}

func fitsInSmallIntReference(i int32) bool {
	return false
}

func smallIntReference(i int32) KReference {
	panic("cannot call function smallIntReference in this implementation")
}

func getSmallInt(ref KReference) (int32, bool) {
	return 0, false
}

func (ms *ModelState) getBigIntObject(ref KReference) (*BigInt, bool) {
	castObj, typeOk := ref.(*BigInt)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// IsInt returns true if reference points to an integer
func IsInt(ref KReference) bool {
	_, typeOk := ref.(*BigInt)
	return typeOk
}

// IntZero is a reference to the constant integer 0
var IntZero = &BigInt{Value: big.NewInt(0)}

// IntOne is a reference to the constant integer 1
var IntOne = &BigInt{Value: big.NewInt(1)}

// IntMinusOne is a reference to the constant integer -1
var IntMinusOne = &BigInt{Value: big.NewInt(-1)}

// FromBigInt provides a reference to an integer (big or small)
func (ms *ModelState) FromBigInt(bi *big.Int) KReference {
	return &BigInt{Value: bi}
}

func parseBigInt(str string) (*big.Int, error) {
	if str == "0" {
		return big.NewInt(0), nil
	}
	b := big.NewInt(0)
	b.UnmarshalText([]byte(str))
	if b.Sign() == 0 {
		return nil, &parseIntError{parseVal: str}
	}
	return b, nil
}

// NewIntConstant creates a new integer constant.
func NewIntConstant(stringRepresentation string) KReference {
	b, err := parseBigInt(stringRepresentation)
	if err != nil {
		panic(err)
	}
	return &BigInt{Value: b}
}

// FromInt provides a reference to an integer (big or small)
func (ms *ModelState) FromInt(x int) KReference {
	return ms.FromInt64(int64(x))
}

// FromInt64 provides a reference to an integer (big or small)
func (ms *ModelState) FromInt64(x int64) KReference {
	return &BigInt{Value: big.NewInt(x)}
}

// FromUint64 provides a reference to an integer (big or small)
func (ms *ModelState) FromUint64(x uint64) KReference {
	var z big.Int
	z.SetUint64(x)
	return &BigInt{Value: &z}
}
