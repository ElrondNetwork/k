%COMMENT%

package %PACKAGE_MODEL%

func newBytesReference(refType kreferenceType, startIndex int, length int) KReference {
	return KReference{
		refType:        refType,
		constantObject: false,
		value1:         uint32(startIndex),
		value2:         uint32(length),
		value3:         0,
	}
}

func parseBytesReference(ref KReference) (startIndex int, length int) {
	startIndex = int(ref.value1)
	length = int(ref.value2)
	return
}

// StringEmpty is a reference to an empty string
var StringEmpty = newBytesReference(stringRef, 0, 0)

// BytesEmpty is a reference to a Bytes item with no bytes (length 0)
var BytesEmpty = newBytesReference(bytesRef, 0, 0)

// IsString returns true if reference points to a string
func IsString(ref KReference) bool {
	return ref.refType == stringRef
}

// IsBytes returns true if reference points to a byte array
func IsBytes(ref KReference) bool {
	return ref.refType == bytesRef
}

// GetString converts reference to a Go string, if possbile
func (ms *ModelState) GetString(ref KReference) (string, bool) {
	if ref.refType != stringRef {
		return "", false
	}
	if ref.constantObject {
		ref.constantObject = false
		return constantsModel.GetString(ref)
	}
	startIndex, length := parseBytesReference(ref)
	if length == 0 {
		return "", true
	}
	return string(ms.allBytes[startIndex : startIndex+length]), true
}

// GetBytes yields the cast object for a Bytes reference, if possible.
func (ms *ModelState) GetBytes(ref KReference) ([]byte, bool) {
	if ref.refType != bytesRef {
		return nil, false
	}
	if ref.constantObject {
		ref.constantObject = false
		return constantsModel.GetBytes(ref)
	}
	startIndex, length := parseBytesReference(ref)
	if length == 0 {
		return nil, true
	}
	return ms.allBytes[startIndex : startIndex+length], true
}

// NewString creates a new K string object from a Go string
func (ms *ModelState) NewString(str string) KReference {
	length := len(str)
	if length == 0 {
		return StringEmpty
	}
	startIndex := len(ms.allBytes)
	ms.allBytes = append(ms.allBytes, []byte(str)...)
	return newBytesReference(stringRef, startIndex, length)
}

// NewStringConstant creates a new string constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewStringConstant(s string) KReference {
	ref := constantsModel.NewString(s)
	ref.constantObject = true
	return ref
}

// NewBytes creates a new K string object from a Go string
func (ms *ModelState) NewBytes(value []byte) KReference {
	length := len(value)
	if length == 0 {
		return BytesEmpty
	}
	startIndex := len(ms.allBytes)
	ms.allBytes = append(ms.allBytes, value...)
	return newBytesReference(bytesRef, startIndex, length)
}

// Bytes2String converts a bytes reference to a string reference.
// The neat thing is, because we use the same underlying structure, no data needs to be copied.
func (ms *ModelState) Bytes2String(ref KReference) (KReference, bool) {
	if ref.refType != bytesRef {
		return NullReference, false
	}
	startIndex, length := parseBytesReference(ref)
	return newBytesReference(stringRef, startIndex, length), true
}

// String2Bytes converts a string reference to a bytes reference.
// The neat thing is, because we use the same underlying structure, no data needs to be copied.
func (ms *ModelState) String2Bytes(ref KReference) (KReference, bool) {
	if ref.refType != stringRef {
		return NullReference, false
	}
	startIndex, length := parseBytesReference(ref)
	return newBytesReference(bytesRef, startIndex, length), true
}

// StringSub yields a reference to a substring of a given string.
// Given the structure of our data, no data needs to be copied or moved in this operation.
func StringSub(ref KReference, fromIndex int, toIndex int) (KReference, bool) {
	return subString(stringRef, ref, fromIndex, toIndex)
}

// BytesSub yields a reference to a sub-slice of a given byte slice.
// Given the structure of our data, no data needs to be copied or moved in this operation.
func BytesSub(ref KReference, fromIndex int, toIndex int) (KReference, bool) {
	return subString(bytesRef, ref, fromIndex, toIndex)
}

func subString(expectedRefType kreferenceType, ref KReference, fromIndex int, toIndex int) (KReference, bool) {
	if ref.refType != expectedRefType {
		return NullReference, false
	}
	startIndex, length := parseBytesReference(ref)
	if fromIndex > toIndex || fromIndex < 0 || toIndex < 0 || fromIndex > length {
		return NullReference, false
	}
	if toIndex > length {
		toIndex = length
	}
	return newBytesReference(ref.refType, startIndex+fromIndex, toIndex-fromIndex), true
}

// StringLength yields the length of a string.
func StringLength(ref KReference) (int, bool) {
	if ref.refType != stringRef {
		return 0, false
	}
	_, length := parseBytesReference(ref)
	return length, true
}

// BytesLength yields the length of a byte array.
func BytesLength(ref KReference) (int, bool) {
	if ref.refType != bytesRef {
		return 0, false
	}
	_, length := parseBytesReference(ref)
	return length, true
}
