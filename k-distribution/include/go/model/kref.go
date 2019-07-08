%COMMENT%

package %PACKAGE_MODEL%

// kreferenceType identifies the type of K item referenced by a KReference
type kreferenceType uint64

const (
	nullRef kreferenceType = iota
	emptyKseqRef
	nonEmptyKseqRef
	kapplyRef
	injectedKLabelRef
	ktokenRef
	kvariableRef
	mapRef
	setRef
	listRef
	arrayRef
	smallPositiveIntRef
	smallNegativeIntRef
	bigIntRef
	mintRef
	floatRef
	stringRef
	stringBufferRef
	bytesRef
	boolRef
	bottomRef
)

func isCollectionType(refType kreferenceType) bool {
	return refType == mapRef ||
		refType == listRef ||
		refType == setRef ||
		refType == arrayRef
}

// KReference is a reference to a K item.
// For some types, like bool and small int, the entire state can be kept in the reference object.
// For the others, the reference contains enough data to find the object in the model state.
type KReference uint64

// type KReference struct {
// 	refType        kreferenceType
// 	constantObject bool
// 	value1         uint32
// 	value2         uint32
// 	value3         uint32
// }

// NullReference is the zero-value of KReference. It doesn't point to anything.
// It has type nullRef.
var NullReference = KReference(0)

// The basic encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: is constant = 1, not constant = 0
// - the remaining 58 LS bits: type-specific data

const refTypeShift = 59 // shift right this many bits to get the refType
const refBasicDataShift = 58
const refBasicDataMask = (1 << refBasicDataShift) - 1

func getRefType(ref KReference) kreferenceType {
	return kreferenceType(uint64(ref) >> refTypeShift)
}

func parseKrefBasic(ref KReference) (refType kreferenceType, constant bool, rest uint64) {
	refRaw := uint64(ref)
	rest = refRaw & refBasicDataMask
	refRaw >>= refBasicDataShift
	constant = refRaw&1 == 1
	refRaw >>= 1
	refType = kreferenceType(refRaw)
	return
}

func createKrefBasic(refType kreferenceType, constant bool, rest uint64) KReference {
	refRaw := uint64(refType)
	refRaw <<= 1
	if constant {
		refRaw |= 1
	}
	refRaw <<= refBasicDataShift
	refRaw |= rest
	return KReference(refRaw)
}

func setConstantFlag(ref KReference) KReference {
	refRaw := uint64(ref)
	refRaw |= (1 << refBasicDataShift)
	return KReference(refRaw)
}

func unsetConstantFlag(ref KReference) KReference {
	refRaw := uint64(ref)
	refRaw &^= (1 << refBasicDataShift) // bit clear operator
	return KReference(refRaw)
}

// big int reference structure (from MSB to LSB):
// - 5 bits: reference type
// - 1 bit: is constant = 1, not constant = 0
// - 26 bits: recycle count
// - 32 bits: index

const refBigIntRecycleCountBits = 26
const refBigIntRecycleCountMask = (1 << refBigIntRecycleCountBits) - 1
const refBigIntIndexBits = 32
const refBigIntIndexMask = (1 << refBigIntIndexBits) - 1

func createKrefBigInt(constant bool, recycleCount uint64, index uint64) KReference {
	ref := uint64(bigIntRef) << 1
	if constant {
		ref |= 1
	}
	ref <<= refBigIntRecycleCountBits
	ref |= recycleCount
	ref <<= refBigIntIndexBits
	ref |= index
	return KReference(ref)
}

func parseKrefBigInt(ref KReference) (isBigInt bool, constant bool, recycleCount uint64, index uint64) {
	refRaw := uint64(ref)
	index = refRaw & refBigIntIndexMask
	refRaw >>= refBigIntIndexBits
	recycleCount = refRaw & refBigIntRecycleCountMask
	refRaw >>= refBigIntRecycleCountBits
	constant = refRaw&1 == 1
	refRaw >>= 1
	isBigInt = refRaw == uint64(bigIntRef)
	return
}

// The collection encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: ignored
// - 13 bits: Label
// - 13 bits: Sort
// - 32 bits: object index

const refCollectionSortShift = 13
const refCollectionSortMask = (1 << refCollectionSortShift) - 1
const refCollectionLabelShift = 13
const refCollectionLabelMask = (1 << refCollectionLabelShift) - 1
const refCollectionIndexShift = 32
const refCollectionIndexMask = (1 << refCollectionIndexShift) - 1

func parseKrefCollection(ref KReference) (refType kreferenceType, sort Sort, label KLabel, index uint64) {
	refRaw := uint64(ref)
	index = refRaw & refCollectionIndexMask
	refRaw >>= refCollectionIndexShift
	label = KLabel(refRaw & refCollectionLabelMask)
	refRaw >>= refCollectionLabelShift
	sort = Sort(refRaw & refCollectionSortMask)
	refRaw >>= refCollectionSortShift
	refRaw >>= 1 // ignore constant flag
	refType = kreferenceType(refRaw)
	return
}

func createKrefCollection(refType kreferenceType, sort Sort, label KLabel, index uint64) KReference {
	refRaw := uint64(refType)
	refRaw <<= 1
	refRaw <<= refCollectionSortShift
	refRaw |= uint64(sort)
	refRaw <<= refCollectionLabelShift
	refRaw |= uint64(label)
	refRaw <<= refCollectionIndexShift
	refRaw |= index
	return KReference(refRaw)
}

// The KApply encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: ignored
// - 13 bits: Sort
// - 13 bits: arity
// - 32 bits: arguments index

const refKApplyLabelShift = 13
const refKApplyLabelMask = (1 << refKApplyLabelShift) - 1
const refKApplyArityShift = 13
const refKApplyArityMask = (1 << refKApplyArityShift) - 1
const refKApplyIndexShift = 32
const refKApplyIndexMask = (1 << refKApplyIndexShift) - 1
const refKApplyTypeAsUint = uint64(kapplyRef)

func parseKrefKApply(ref KReference) (isKApply bool, label KLabel, arity uint64, index uint64) {
	refRaw := uint64(ref)
	index = refRaw & refKApplyIndexMask
	refRaw >>= refKApplyIndexShift
	arity = refRaw & refKApplyArityMask
	refRaw >>= refKApplyArityShift
	label = KLabel(refRaw & refKApplyLabelMask)
	refRaw >>= refKApplyLabelShift
	refRaw >>= 1 // ignore constant flag
	isKApply = refRaw == refKApplyTypeAsUint
	return
}

func createKrefKApply(label KLabel, arity uint64, index uint64) KReference {
	refRaw := refKApplyTypeAsUint
	refRaw <<= 1
	refRaw <<= refKApplyLabelShift
	refRaw |= uint64(label)
	refRaw <<= refKApplyArityShift
	refRaw |= arity
	refRaw <<= refKApplyIndexShift
	refRaw |= index
	return KReference(refRaw)
}

// The K sequence encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: ignored
// - 26 bits: length
// - 32 bits: head (element) index

const refNonEmptyKseqLengthShift = 26
const refNonEmptyKseqLengthMask = (1 << refNonEmptyKseqLengthShift) - 1
const refNonEmptyKseqIndexShift = 32
const refNonEmptyKseqIndexMask = (1 << refNonEmptyKseqIndexShift) - 1
const refNonEmptyKseqTypeAsUint = uint64(nonEmptyKseqRef)

func parseKrefKseq(ref KReference) (refType kreferenceType, elemIndex uint64, length uint64) {
	refRaw := uint64(ref)
	elemIndex = refRaw & refNonEmptyKseqIndexMask
	refRaw >>= refNonEmptyKseqIndexShift
	length = refRaw & refNonEmptyKseqLengthMask
	refRaw >>= refNonEmptyKseqLengthShift
	refRaw >>= 1 // ignore constant flag
	refType = kreferenceType(refRaw)
	return
}

func createKrefNonEmptyKseq(elemIndex uint64, length uint64) KReference {
	refRaw := refNonEmptyKseqTypeAsUint
	refRaw <<= 1
	refRaw <<= refNonEmptyKseqLengthShift
	refRaw |= length
	refRaw <<= refNonEmptyKseqIndexShift
	refRaw |= elemIndex
	return KReference(refRaw)
}

// The K token encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: is constant = 1, not constant = 0
// - 13 bits: sort
// - 13 bits: length
// - 32 bits: value string start index in allBytes

const refKTokenSortShift = 13
const refKTokenSortMask = (1 << refKTokenSortShift) - 1
const refKTokenLengthShift = 13
const refKTokenLengthMask = (1 << refKTokenLengthShift) - 1
const refKTokenIndexShift = 32
const refKTokenIndexMask = (1 << refKTokenIndexShift) - 1
const refKTokenTypeAsUint = uint64(ktokenRef)

func parseKrefKToken(ref KReference) (isKToken bool, constant bool, sort Sort, length uint64, index uint64) {
	refRaw := uint64(ref)
	index = refRaw & refKTokenIndexMask
	refRaw >>= refKTokenIndexShift
	length = refRaw & refKTokenLengthMask
	refRaw >>= refKTokenLengthShift
	sort = Sort(refRaw & refKTokenSortMask)
	refRaw >>= refKTokenSortShift
	constant = refRaw&1 == 1
	refRaw >>= 1
	isKToken = refRaw == refKTokenTypeAsUint
	return
}

func createKrefKToken(constant bool, sort Sort, length uint64, index uint64) KReference {
	refRaw := refKTokenTypeAsUint
	refRaw <<= 1
	if constant {
		refRaw |= 1
	}
	refRaw <<= refKTokenSortShift
	refRaw |= uint64(sort)
	refRaw <<= refKTokenLengthShift
	refRaw |= length
	refRaw <<= refKTokenIndexShift
	refRaw |= index
	return KReference(refRaw)
}

// The byte array and string encoding is as follows (from MSB to LSB):
// - first 5 bits: reference type
// - next 1 bit: is constant = 1, not constant = 0
// - 26 bits: length
// - 32 bits: head (element) index

const refBytesLengthShift = 26
const refBytesLengthMask = (1 << refBytesLengthShift) - 1
const refBytesIndexShift = 32
const refBytesIndexMask = (1 << refBytesIndexShift) - 1

func parseKrefBytes(ref KReference) (refType kreferenceType, constant bool, startIndex uint64, length uint64) {
	refRaw := uint64(ref)
	startIndex = refRaw & refBytesIndexMask
	refRaw >>= refBytesIndexShift
	length = refRaw & refBytesLengthMask
	refRaw >>= refBytesLengthShift
	constant = refRaw&1 == 1
	refRaw >>= 1
	refType = kreferenceType(refRaw)
	return
}

func createKrefBytes(refType kreferenceType, constant bool, startIndex uint64, length uint64) KReference {
	refRaw := uint64(refType)
	refRaw <<= 1
	if constant {
		refRaw |= 1
	}
	refRaw <<= refBytesLengthShift
	refRaw |= length
	refRaw <<= refBytesIndexShift
	refRaw |= startIndex
	return KReference(refRaw)
}
