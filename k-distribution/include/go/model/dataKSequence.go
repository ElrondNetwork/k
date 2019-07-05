%COMMENT%

package %PACKAGE_MODEL%

// EmptyKSequence is the KSequence with no elements.
// To simplify things, it is a separate reference type.
var EmptyKSequence = KReference{refType: emptyKseqRef}

type ksequenceElem struct {
	head KReference
	tail KReference
}

func createNonEmptyKseqRef(elemIndex int, length int) KReference {
	return KReference{refType: nonEmptyKseqRef, value1: uint32(elemIndex), value2: uint32(length)}
}

func nonEmptyKseqRefParse(ref KReference) (ok bool, elemIndex int, length int) {
	if ref.refType != nonEmptyKseqRef {
		return false, 0, 0
	}
	return true, int(ref.value1), int(ref.value2)
}

// IsNonEmptyKSequenceMinimumLength returns true for any K sequence with length greater of equal than given argument.
// Returns false for EmptyKSequence.
// Especially used for pattern matching.
func (ms *ModelState) IsNonEmptyKSequenceMinimumLength(ref KReference, minimumLength int) bool {
	if ref.refType == emptyKseqRef {
		return false
	}
	isNonEmptyKseq, _, length := nonEmptyKseqRefParse(ref)
	if !isNonEmptyKseq {
		return minimumLength == 1
	}
	return length >= minimumLength
}

// NewKSequence creates new KSequence instance with given references
func (ms *ModelState) NewKSequence(elements []KReference) KReference {
	return ms.AssembleKSequence(elements...)
}

// KSequenceIsEmpty returns true if KSequence has no elements
func (ms *ModelState) KSequenceIsEmpty(ref KReference) bool {
	return ref.refType == emptyKseqRef
}

// KSequenceGet yields element at position.
func (ms *ModelState) KSequenceGet(ref KReference, position int) KReference {
	for i := 0; i < position; i++ {
		isNonEmptyKseq, elemIndex, _ := nonEmptyKseqRefParse(ref)
		if !isNonEmptyKseq {
			panic("bad argument to KSequenceGet: position exceeds K sequence length")
		}
		elem := ms.allKsElements[elemIndex]
		ref = elem.tail
	}

	isNonEmptyKseq, elemIndex, _ := nonEmptyKseqRefParse(ref)
	if !isNonEmptyKseq {
		return ref
	}
	return ms.allKsElements[elemIndex].head
}

// KSequenceLength yields KSequence length
func (ms *ModelState) KSequenceLength(ref KReference) int {
	isNonEmptyKseq, _, length := nonEmptyKseqRefParse(ref)
	if !isNonEmptyKseq {
		panic("bad argument to KSequenceLength: ref is not a reference to a K sequence")
	}
	return length
}

// KSequenceToSlice converts KSequence to a slice of K items
func (ms *ModelState) KSequenceToSlice(ref KReference) []KReference {
	if ref.refType == emptyKseqRef {
		return nil
	}
	isNonEmptyKseq, elemIndex, length := nonEmptyKseqRefParse(ref)
	if !isNonEmptyKseq {
		panic("bad argument to KSequenceToSlice: ref is not a reference to a K sequence")
	}

	var result []KReference
	for isNonEmptyKseq {
		elem := ms.allKsElements[elemIndex]
		result = append(result, elem.head)
		ref = elem.tail
		isNonEmptyKseq, elemIndex, _ = nonEmptyKseqRefParse(ref)
	}

	// last element is not a K sequence
	result = append(result, ref)

	if len(result) != length {
		panic("K sequence reference length does not match actual length of K sequence")
	}

	return result
}

// KSequenceSub yields subsequence starting at position
func (ms *ModelState) KSequenceSub(ref KReference, startPosition int) KReference {
	for i := 0; i < startPosition; i++ {
		isNonEmptyKseq, elemIndex, _ := nonEmptyKseqRefParse(ref)
		if !isNonEmptyKseq {
			if i == startPosition-1 {
				return EmptyKSequence
			}
			panic("bad argument to KSequenceSub: startPosition exceeds original K sequence")
		} else {
			elem := ms.allKsElements[elemIndex]
			ref = elem.tail
		}
	}

	return ref
}

// KSequenceSplitHeadTail  extracts first element of a KSequence, extracts the rest, if possible
// will treat non-KSequence as if they were KSequences of length 1
func (ms *ModelState) KSequenceSplitHeadTail(ref KReference) (ok bool, head KReference, tail KReference) {
	if ref.refType == emptyKseqRef {
		return false, NoResult, EmptyKSequence
	}

	if ref.refType == nonEmptyKseqRef {
		elem := ms.allKsElements[ref.value1]
		return true, elem.head, elem.tail
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, ref, EmptyKSequence
}

// AssembleKSequence appends all given arguments into a KSequence.
// It flattens any KSequences among the arguments.
// Never returns KSequence of 1 element, it returns the element directly instead
func (ms *ModelState) AssembleKSequence(refs ...KReference) KReference {
	head := EmptyKSequence
	resultLength := 0

	for i := len(refs) - 1; i >= 0; i-- {
		ref := refs[i]
		if ref.refType == emptyKseqRef {
			// nothing, ignore
		} else {
			if head.refType == emptyKseqRef {
				// first to be added
				head = ref
				if ref.refType == nonEmptyKseqRef {
					// result extends another K sequence
					resultLength = int(head.value2)
				} else {
					resultLength = 1
				}
			} else {
				// append to the simple linked list
				// like the cons in cons lists
				if ref.refType == nonEmptyKseqRef {
					// flatten K sequence given as argument
					// concatenate entire sub-sequence to beginning of result sequence
					slice := ms.KSequenceToSlice(ref)
					slice = append(slice, head)
					head = ms.AssembleKSequence(slice...)
				} else {
					// add 1 element to beginning of list
					newHead := ksequenceElem{
						head: ref,
						tail: head,
					}
					resultLength++
					newIndex := len(ms.allKsElements)
					ms.allKsElements = append(ms.allKsElements, newHead)
					head = createNonEmptyKseqRef(newIndex, resultLength)
				}
			}
		}
	}

	return head
}
