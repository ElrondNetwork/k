%COMMENT%

package %PACKAGE_MODEL%

// NewKSequence ... creates new KSequence instance with elements
func NewKSequence(elements []K) KSequence {
	return KSequence(elements)
}

// IsEmpty ... returns true if KSequence has no elements
func (k KSequence) IsEmpty() bool {
	ks := []K(k)
	return len(ks) == 0
}

// IsEmpty ... KSequence length
func (k KSequence) Length() int {
	return len([]K(k))
}

// IsEmpty ... element at position
func (k KSequence) Get(position int) K {
    ks := []K(k)
	return ks[position]
}

// SubSequence ... subsequence starting at position
func (k KSequence) SubSequence(startPosition int) KSequence {
    ks := []K(k)
	return KSequence(ks[startPosition:])
}

// TrySplitToHeadTail ... extracts first element of a KSequence, extracts the rest, if possible
// will treat non-KSequence as if they were KSequences of length 1
func TrySplitToHeadTail(k K) (ok bool, head K, tail K) {
	if kseq, isKseq := k.(KSequence); isKseq {
		switch kseq.Length() {
		case 0:
			return false, NoResult, EmptyKSequence
		case 1:
			return true, kseq.Get(0), EmptyKSequence
		case 2:
			return true, kseq.Get(0), kseq.Get(1)
		default:
		    ks := []K(kseq)
			return true, kseq.Get(0), KSequence(ks[1:])
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, k, EmptyKSequence
}

// AssembleKSequence ... appends all elements into a KSequence
// flattens if there are any KSequences among the elements (but only on 1 level, does not handle multiple nesting)
// never returns KSequence of 1 element, it returns the element directly instead
func AssembleKSequence(elements ...K) K {
	var newKs []K
	for _, element := range elements {
		if kseqElem, isKseq := element.(KSequence); isKseq {
			newKs = append(newKs, []K(kseqElem)...)
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
	return NewKSequence(newKs)
}