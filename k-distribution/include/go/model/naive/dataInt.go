%COMMENT%

package %PACKAGE_MODEL%

import (
	"math"
	"math/big"
)

// BigInt is a KObject representing a big int in K
type BigInt struct {
	Value *big.Int
}

// IsInt returns true if reference points to an integer
func IsInt(ref KReference) bool {
	_, typeOk := ref.(*BigInt)
	return typeOk
}

var bigOne = big.NewInt(1)

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

// NewIntConstant creates a new integer constant.
func NewIntConstant(stringRepresentation string) KReference {
	i, err := parseInt(stringRepresentation)
	if err != nil {
		panic(err)
	}
	return i
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

// IntFromByte provides a reference to a small integer
func (ms *ModelState) IntFromByte(x byte) KReference {
	var z big.Int
	z.SetUint64(uint64(x))
	return &BigInt{Value: &z}
}

// IntFromBytes provides a reference to an integer
func (ms *ModelState) IntFromBytes(bytes []byte) KReference {
	z := big.NewInt(0)
	z.SetBytes(bytes)
	return &BigInt{Value: z}
}

func parseInt(s string) (KReference, error) {
	if s == "0" {
		return IntZero, nil
	}
	b := big.NewInt(0)
	b.UnmarshalText([]byte(s))
	if b.Sign() == 0 {
		return IntZero, &parseIntError{parseVal: s}
	}
	return &BigInt{Value: b}, nil
}

// ParseInt creates K int from string representation
func (ms *ModelState) ParseInt(s string) (KReference, error) {
	return parseInt(s)
}

// IntFromString does the same as ParseInt but panics instead of returning an error
func (ms *ModelState) IntFromString(s string) KReference {
	i, err := ms.ParseInt(s)
	if err != nil {
		panic(err)
	}
	return i
}

// GetBigIntObject yields the cast object for a BigInt reference, if possible.
func (ms *ModelState) GetBigIntObject(ref KReference) (*BigInt, bool) {
	castObj, typeOk := ref.(*BigInt)
	if !typeOk {
		return nil, false
	}
	return castObj, true
}

// GetBigInt yields a big.Int cast from any K integer object, if possible.
func (ms *ModelState) GetBigInt(ref KReference) (*big.Int, bool) {
	bi, isBigInt := ms.GetBigIntObject(ref)
	if !isBigInt {
		return nil, false
	}
	return bi.Value, true
}

// IsZero returns true if an item represents number 0
func (k *BigInt) IsZero() bool {
	return k.Value.Sign() == 0
}

// IsPositive ... true if represented number is >= 0
func (k *BigInt) IsPositive() bool {
	return k.Value.Sign() >= 0
}

// IsNegative ... true if represented number is < 0
func (k *BigInt) IsNegative() bool {
	return k.Value.Sign() < 0
}

// ToUint32 ... converts to uint if possible, returns (0, false) if not
func (k *BigInt) ToUint32() (uint, bool) {
	if !k.Value.IsUint64() {
		return 0, false
	}

	u64 := k.Value.Uint64()
	if u64 > math.MaxUint32 {
		return 0, false
	}

	return uint(u64), true
}

// ToInt32 ... converts to int if possible, returns (0, false) if not
func (k *BigInt) ToInt32() (int, bool) {
	if !k.Value.IsInt64() {
		return 0, false
	}

	i64 := k.Value.Int64()
	if i64 < math.MinInt32 || i64 > math.MaxInt32 {
		return 0, false
	}

	return int(i64), true
}

// ToPositiveInt32 ... converts to int32 if possible, returns (0, false) if not
// also rejects negative numbers, so we don't have to test for that again
func (k *BigInt) ToPositiveInt32() (int, bool) {
	if !k.Value.IsInt64() {
		return 0, false
	}

	i64 := k.Value.Int64()
	if i64 < 0 || i64 > math.MaxInt32 {
		return 0, false
	}

	return int(i64), true
}

// ToByte converts to 1 byte if possible, returns (0, false) if not
func (k *BigInt) ToByte() (byte, bool) {
	if !k.Value.IsUint64() {
		return 0, false
	}

	u64 := k.Value.Uint64()
	if u64 > 255 {
		return 0, false
	}

	return byte(u64), true
}

// BigIntToTwosComplementBytes ... returns a byte array representation, 2's complement if number is negative
// big endian
func BigIntToTwosComplementBytes(i *big.Int, bytesLength int) []byte {
	var resultBytes []byte
	switch i.Sign() {
	case -1:
		// compute 2's complement
		plus1 := big.NewInt(0)
		plus1.Add(i, big.NewInt(1)) // add 1
		plus1Bytes := plus1.Bytes()
		offset := len(plus1Bytes) - bytesLength
		resultBytes = make([]byte, bytesLength)
		for i := 0; i < bytesLength; i++ {
			j := offset + i
			if j < 0 {
				resultBytes[i] = 255 // pad left with 11111111
			} else {
				resultBytes[i] = ^plus1Bytes[j] // also negate every bit
			}
		}
		break
	case 0:
		// just zeroes
		resultBytes = make([]byte, bytesLength)
		break
	case 1:
		originalBytes := i.Bytes()
		resultBytes = make([]byte, bytesLength)
		offset := len(originalBytes) - bytesLength
		for i := 0; i < bytesLength; i++ {
			j := offset + i
			if j < 0 {
				resultBytes[i] = 0 // pad left with 00000000
			} else {
				resultBytes[i] = originalBytes[j]
			}
		}
		break
	}

	return resultBytes
}

// TwosComplementBytesToBigInt ... convert a byte array to a number
// interprets input as a 2's complement representation if the first bit (most significant) is 1
// big endian
func TwosComplementBytesToBigInt(twosBytes []byte) *big.Int {
	testBit := twosBytes[0] >> 7
	result := new(big.Int)
	if testBit == 0 {
		// positive number, no further processing required
		result.SetBytes(twosBytes)
	} else {
		// convert to negative number
		notBytes := make([]byte, len(twosBytes))
		for i, b := range twosBytes {
			notBytes[i] = ^b // negate every bit
		}
		result.SetBytes(notBytes)
		result.Neg(result)
		result.Sub(result, bigOne) // -1
	}

	return result
}
