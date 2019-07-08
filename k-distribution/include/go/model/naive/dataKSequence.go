%COMMENT%

package %PACKAGE%

// KSequence is a sequence of K items
type KSequence struct {
	Ks []K
}

// EmptyKSequence is the KSequence with no elements
var EmptyKSequence = &KSequence{Ks: nil}

// IsNonEmptyKSequenceMinimumLength returns true for any K sequence with length greater of equal than given argument.
// Returns false for EmptyKSequence.
// Especially used for pattern matching.
func (ms *ModelState) IsNonEmptyKSequenceMinimumLength(ref KReference, minimumLength int) bool {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		return false
	}
	if len(castObj.Ks) == 0 {
		return false
	}
	return len(castObj.Ks) >= minimumLength
}

// NewKSequence creates new KSequence instance with given references
func (ms *ModelState) NewKSequence(elements []KReference) KReference {
	return &KSequence{Ks: elements}
}

// KSequenceIsEmpty returns true if KSequence has no elements
func (ms *ModelState) KSequenceIsEmpty(ref KReference) bool {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		return false
	}
	return len(castObj.Ks) == 0
}

// KSequenceGet yields element at position
// Caution: no checks are performed that the position is valid
func (ms *ModelState) KSequenceGet(ref KReference, position int) KReference {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		panic("bad argument to KSequenceGet: ref is not a reference to a K sequence")
	}
	return castObj.Ks[position]
}

// KSequenceLength yields KSequence length
func (ms *ModelState) KSequenceLength(ref KReference) int {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		panic("bad argument to KSequenceLength: ref is not a reference to a K sequence")
	}
	return len(castObj.Ks)
}

// KSequenceToSlice converts KSequence to a slice of K items
func (ms *ModelState) KSequenceToSlice(ref KReference) []KReference {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		panic("bad argument to KSequenceToSlice: ref is not a reference to a K sequence")
	}
	return castObj.Ks
}

// KSequenceSub yields subsequence starting at position
func (ms *ModelState) KSequenceSub(ref KReference, startPosition int) KReference {
	castObj, typeOk := ref.(*KSequence)
	if !typeOk {
		panic("bad argument to KSequenceSub: ref is not a reference to a K sequence")
	}
	return &KSequence{Ks: castObj.Ks[startPosition:]}
}

// KSequenceSplitHeadTail  extracts first element of a KSequence, extracts the rest, if possible
// will treat non-KSequence as if they were KSequences of length 1
func (ms *ModelState) KSequenceSplitHeadTail(ref KReference) (ok bool, head KReference, tail KReference) {
	if kseq, isKseq := ref.(*KSequence); isKseq {
		switch len(kseq.Ks) {
		case 0:
			return false, NoResult, EmptyKSequence
		case 1:
			return true, kseq.Ks[0], EmptyKSequence
		case 2:
			return true, kseq.Ks[0], kseq.Ks[1]
		default:
			return true, kseq.Ks[0], &KSequence{Ks: kseq.Ks[1:]}
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, ref, EmptyKSequence
}

// AssembleKSequence appends all elements into a KSequence.
// It flattens any KSequences among the elements (but only on 1 level, does not handle multiple nesting).
// Never returns KSequence of 1 element, it returns the element directly instead
func (ms *ModelState) AssembleKSequence(elements ...KReference) KReference {
	var newKs []KReference
	for _, element := range elements {
		if kseqElem, isKseq := element.(*KSequence); isKseq {
			newKs = append(newKs, kseqElem.Ks...)
		} else {
			newKs = append(newKs, element)
		}
	}
	if len(newKs) == 0 {
		return EmptyKSequence
	}
	if len(newKs) == 1 {
		return newKs[0]
	}
	return &KSequence{Ks: newKs}
}
