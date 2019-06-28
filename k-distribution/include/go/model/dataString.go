%COMMENT%

package %PACKAGE_MODEL%

import (
    "strings"
)

// StringEmpty is a reference to an empty string
var StringEmpty = addConstantObject(&String{Value: ""})

// BytesEmpty is a reference to a Bytes item with no bytes (length 0)
var BytesEmpty = addConstantObject(&Bytes{Value: nil})

// IsString returns true if reference points to a string
func IsString(ref KReference) bool {
	return ref.refType == stringRef
}

// IsBytes returns true if reference points to a byte array
func IsBytes(ref KReference) bool {
	return ref.refType == bytesRef
}

// IsStringBuffer returns true if reference points to a string buffer
func IsStringBuffer(ref KReference) bool {
	return ref.refType == stringBufferRef
}

// GetBigIntObject yields the cast object for a String reference, if possible.
func (ms *ModelState) GetStringObject(ref KReference) (*String, bool) {
	if ref.refType == stringRef {
		obj := ms.getObject(ref)
        castObj, typeOk := obj.(*String)
        if !typeOk {
            panic("wrong object type for reference")
        }
        return castObj, true
	}

	return nil, false
}

// GetString converts reference to a Go string, if possbile
func (ms *ModelState) GetString(ref KReference) (string, bool) {
	castObj, typeOk := ms.GetStringObject(ref)
	if !typeOk {
		return "", false
	}
	return castObj.Value, true
}

// GetBytesObject yields the cast object for a Bytes reference, if possible.
func (ms *ModelState) GetBytesObject(ref KReference) (*Bytes, bool) {
	if ref.refType != bytesRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*Bytes)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// GetStringBufferObject yields the cast object for a StringBuffer reference, if possible.
func (ms *ModelState) GetStringBufferObject(ref KReference) (*StringBuffer, bool) {
	if ref.refType != stringBufferRef {
		return nil, false
	}
	obj := ms.getObject(ref)
	castObj, typeOk := obj.(*StringBuffer)
	if !typeOk {
		panic("wrong object type for reference")
	}
	return castObj, true
}

// NewString creates a new K string object from a Go string
func (ms *ModelState) NewString(str string) KReference {
	return ms.addObject(&String{Value: str})
}

// NewStringConstant creates a new string constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewStringConstant(s string) KReference {
	ref := constantsModel.NewString(s)
	ref.constantObject = true
	return ref
}

// NewStringBuffer creates a new object and returns the reference.
func (ms *ModelState) NewStringBuffer() KReference {
	return ms.addObject(&StringBuffer{Value: strings.Builder{}})
}