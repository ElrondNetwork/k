%COMMENT%

package %PACKAGE_MODEL%

// KToken is a KObject representing a KToken item in K
type KToken struct {
	Value string
	Sort  Sort
}

func newKTokenReference(sortInt int, startIndex int, length int) KReference {
	return KReference{
		refType:        ktokenRef,
		constantObject: false,
		value1:         uint32(startIndex),
		value2:         uint32(length),
		value3:         uint32(sortInt),
	}
}

func parseKTokenReference(ref KReference) (sortInt int, startIndex int, length int) {
	startIndex = int(ref.value1)
	length = int(ref.value2)
	sortInt = int(ref.value3)
	return
}

// GetKTokenObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetKTokenObject(ref KReference) (*KToken, bool) {
	if ref.refType != ktokenRef {
		return nil, false
	}
	if ref.constantObject {
		ref.constantObject = false
		return constantsModel.GetKTokenObject(ref)
	}
	sortInt, startIndex, length := parseKTokenReference(ref)
	value := ""
	if length > 0 {
		value = string(ms.allBytes[startIndex : startIndex+length])
	}
	return &KToken{
		Sort:  Sort(sortInt),
		Value: value,
	}, true
}

// NewKToken creates a new object and returns the reference.
func (ms *ModelState) NewKToken(sort Sort, value string) KReference {
	length := len(value)
	if length == 0 {
		return newKTokenReference(int(sort), 0, 0)
	}
	startIndex := len(ms.allBytes)
	ms.allBytes = append(ms.allBytes, []byte(value)...)
	return newKTokenReference(int(sort), startIndex, length)
}

// NewKTokenConstant creates a new KToken constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewKTokenConstant(sort Sort, value string) KReference {
	ref := constantsModel.NewKToken(sort, value)
	ref.constantObject = true
	return ref
}
