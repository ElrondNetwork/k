package %PACKAGE_MODEL%

import (
	"math/big"
)

// IntZero ... K Int with value zero
var IntZero = &Int{Value: big.NewInt(0)}

// IntMinusOne ... K Int with value -1
var IntMinusOne = &Int{Value: big.NewInt(-1)}

// BoolTrue ... K boolean value with value true
var BoolTrue = &Bool{Value: true}

// BoolFalse ... K boolean value with value false
var BoolFalse = &Bool{Value: false}

// InternedBottom ... usually used as a dummy object
var InternedBottom = &Bottom{}

// NoResult ... what to return when a function returns an error
var NoResult = &Bottom{}

// NewInt ... provides new Int instance
func NewInt(bi *big.Int) *Int {
	return &Int{Value: bi}
}

// NewIntFromInt ... provides new Int instance
func NewIntFromInt(x int) *Int {
	return NewIntFromInt64(int64(x))
}

// NewIntFromInt64 ... provides new Int instance
func NewIntFromInt64(x int64) *Int {
	return &Int{Value: big.NewInt(x)}
}

// NewIntFromUint64 ... provides new Int instance
func NewIntFromUint64(x uint64) *Int {
	var z big.Int
	z.SetUint64(x)
	return &Int{Value: &z}
}

// ParseInt ... creates K int from string representation
func ParseInt(s string) (*Int, error) {
	b := big.NewInt(0)
	if s != "0" {
		b.UnmarshalText([]byte(s))
		if b.Cmp(IntZero.Value) == 0 {
			return IntZero, &parseIntError{parseVal: s}
		}
	}
	return NewInt(b), nil
}

// NewIntFromString ... same as ParseInt but panics instead of error
func NewIntFromString(s string) *Int {
	i, err := ParseInt(s)
	if err != nil {
		panic(err)
	}
	return i
}

// NewString ... Creates a new K string object from a Go string
func NewString(str string) *String {
	return &String{Value: str}
}

// String ... Yields a Go string representation of the K String
func (k *String) String() string {
	return k.Value
}

// ToBool ... Convert Go bool to K Bool
func ToBool(b bool) *Bool {
	if b {
		return BoolTrue
	}
	return BoolFalse
}

// IsEmpty ... returns true if KSequence has no elements
func (k *KSequence) IsEmpty() bool {
	return len(k.Ks) == 0
}
