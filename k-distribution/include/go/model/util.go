package %PACKAGE_MODEL%

import (
    "math"
	"math/big"
	"sort"
)

// IntZero ... K Int with value zero
var IntZero = &Int{Value: big.NewInt(0)}

// IntOne ... K Int with value 1
var IntOne = &Int{Value: big.NewInt(1)}

// IntMinusOne ... K Int with value -1
var IntMinusOne = &Int{Value: big.NewInt(-1)}

// BytesEmpty ... Bytes item with no bytes (length 0)
var BytesEmpty = &Bytes{Value: nil}

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

// IsZero ... true if item represents number 0
func (k *Int) IsZero() bool {
	return k.Value.Sign() == 0
}

// IsNegative ... true if represented number is < 0
func (k *Int) IsNegative() bool {
	return k.Value.Sign() < 0
}

// ToUint32 ... converts to uint if possible, returns (0, false) if not
func (k *Int) ToUint32() (uint, bool) {
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
func (k *Int) ToInt32() (int, bool) {
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
func (k *Int) ToPositiveInt32() (int, bool) {
	if !k.Value.IsInt64() {
		return 0, false
	}

	i64 := k.Value.Int64()
	if i64 < 0 || i64 > math.MaxInt32 {
		return 0, false
	}

	return int(i64), true
}

// ToByte ... converts to 1 byte if possible, returns (0, false) if not
func (k *Int) ToByte() (byte, bool) {
	if !k.Value.IsUint64() {
		return 0, false
	}

	u64 := k.Value.Uint64()
	if u64 > 255 {
		return 0, false
	}

	return byte(u64), true
}

// IsEmpty ... returns true if Bytes is the empty byte slice
func (k *Bytes) IsEmpty() bool {
	return len(k.Value) == 0
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

// IsTrue ... Checks if argument is identical to the K Bool with the value true
func IsTrue(c K) bool {
	if b, typeOk := c.(*Bool); typeOk {
		return b.Value
	}
	return false
}

// IsEmpty ... returns true if KSequence has no elements
func (k *KSequence) IsEmpty() bool {
	return len(k.Ks) == 0
}


// MapKeyValuePair ... just a pair of key and value that was stored in a map
type MapKeyValuePair struct {
	KeyAsString string
	Key         K
	Value       K
}

// ToOrderedKeyValuePairs ... Yields a list of key-value pairs, ordered by the string representation of the keys
func (k *Map) ToOrderedKeyValuePairs() []MapKeyValuePair {
	result := make([]MapKeyValuePair, len(k.Data))

	var keysAsString []string
	stringKeysToPair := make(map[string]MapKeyValuePair)
	for key, val := range k.Data {
		keyAsString := key.String()
		keysAsString = append(keysAsString, keyAsString)
		keyAsK, err := key.ToKItem()
		if err != nil {
			panic(err)
		}
		pair := MapKeyValuePair{KeyAsString: keyAsString, Key: keyAsK, Value: val}
		stringKeysToPair[keyAsString] = pair
	}
	sort.Strings(keysAsString)
	for i, keyAsString := range keysAsString {
		pair, _ := stringKeysToPair[keyAsString]
		result[i] = pair
	}

	return result
}
