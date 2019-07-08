%COMMENT%

package %PACKAGE%

// KToken is a KObject representing a KToken item in K
type KToken struct {
	Value string
	Sort  Sort
}

// GetKTokenObject yields the cast object for a KApply reference, if possible.
func (ms *ModelState) GetKTokenObject(ref KReference) (*KToken, bool) {
	isKToken, constant, sort, length, index := parseKrefKToken(ref)
	if !isKToken {
		return nil, false
	}
	if constant {
		ref = unsetConstantFlag(ref)
		return constantsModel.GetKTokenObject(ref)
	}
	value := ""
	if length > 0 {
		value = string(ms.allBytes[index : index+length])
	}
	return &KToken{
		Sort:  Sort(sort),
		Value: value,
	}, true
}

// KTokenValue yields the value of a KToken object.
func (ms *ModelState) KTokenValue(ref KReference) string {
	isKToken, constant, _, length, index := parseKrefKToken(ref)
	if !isKToken {
		panic("KTokenValue called for reference to item other than KToken")
	}
	if constant {
		ref = unsetConstantFlag(ref)
		return constantsModel.KTokenValue(ref)
	}
	if length == 0 {
		return ""
	}
	return string(ms.allBytes[index : index+length])
}

// NewKToken creates a new object and returns the reference.
func (ms *ModelState) NewKToken(sort Sort, value string) KReference {
	length := uint64(len(value))
	if length == 0 {
		return createKrefKToken(false, uint64(sort), length, 0)
	}
	startIndex := uint64(len(ms.allBytes))
	ms.allBytes = append(ms.allBytes, []byte(value)...)
	return createKrefKToken(false, uint64(sort), length, startIndex)
}

// NewKTokenConstant creates a new KToken constant, which is saved statically.
// Do not use for anything other than constants, since these never get cleaned up.
func NewKTokenConstant(sort Sort, value string) KReference {
	ref := constantsModel.NewKToken(sort, value)
	ref = setConstantFlag(ref)
	return ref
}
