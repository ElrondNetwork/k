%COMMENT%

package %PACKAGE_MODEL%

// we keep all KSequences concatenated into one 2d structure
// each KSequence object is actually a pointer to the first element, its index in this slice
// the first element of allKs is the empty sequence
var allKs = [][]K{[]K{}}

// ClearModel ... clean up any data left from previous executions, to save memory
func ClearModel() {
	allKs = [][]K{[]K{}}
}

// EmptyKSequence ... the KSequence with no elements
var EmptyKSequence = KSequence{sequenceIndex: 0, headIndex: 0}

// NewKSequence ... creates new KSequence instance with elements
func NewKSequence(elements []K) KSequence {
	newSequenceIndex := len(allKs)
	allKs = append(allKs, elements)
	return KSequence{sequenceIndex: newSequenceIndex, headIndex: 0}
}

// IsEmpty ... returns true if KSequence has no elements
func (k KSequence) IsEmpty() bool {
	return k.Length() == 0
}

// Get ... element at position
// Caution: no checks are performed that the position is valid
func (k KSequence) Get(position int) K {
	seq := allKs[k.sequenceIndex]
	return seq[k.headIndex+position]
}

// Length ... KSequence length
func (k KSequence) Length() int {
	return len(allKs[k.sequenceIndex]) - k.headIndex
}

// ToSlice ... converts KSequence to a slice of K items
func (k KSequence) ToSlice() []K {
	return allKs[k.sequenceIndex][k.headIndex:]
}

// SubSequence ... subsequence starting at position
func (k KSequence) SubSequence(startPosition int) KSequence {
	return KSequence{sequenceIndex: k.sequenceIndex, headIndex: k.headIndex + startPosition}
}

// TrySplitToHeadTail ... extracts first element of a KSequence, extracts the rest, if possible
// will treat non-KSequence as if they were KSequences of length 1
func TrySplitToHeadTail(k K) (ok bool, head K, tail K) {
	if kseq, isKseq := k.(KSequence); isKseq {
		seq := allKs[kseq.sequenceIndex]
		length := len(seq) - kseq.headIndex
		switch length {
		case 0:
			// empty KSequence, no result
			return false, NoResult, EmptyKSequence
		case 1:
			return true, seq[kseq.headIndex], EmptyKSequence
		case 2:
			// the KSequence has length 2
			// this case is special because here the tail is not a KSequence
			return true, seq[kseq.headIndex], seq[kseq.headIndex+1]
		default:
			// advance head
			return true, seq[kseq.headIndex], KSequence{kseq.sequenceIndex, kseq.headIndex + 1}
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
			newKs = append(newKs, kseqElem.ToSlice()...)
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
