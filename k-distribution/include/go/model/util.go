package %PACKAGE_MODEL%

// BytesEmpty ... Bytes item with no bytes (length 0)
var BytesEmpty = &Bytes{Value: nil}

// InternedBottom ... usually used as a dummy object
var InternedBottom = &Bottom{}

// NoResult ... what to return when a function returns an error
var NoResult = &Bottom{}

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

// IsEmpty ... returns true if KSequence has no elements
func (k *KSequence) IsEmpty() bool {
	return len(k.Ks) == 0
}
