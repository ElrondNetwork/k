%COMMENT%

package %PACKAGE_MODEL%

// we keep all KSequences concatenated into 1 really big slice
// each KSequence object is actually a pointer to the first element, its index in this slice
// each KSequence is terminated by a nil element
// the first element of allKs is the terminator for the empty sequence
var allKs = []K{nil}

// ClearModel ... clean up any data left from previous executions, to save memory
func ClearModel() {
	allKs = []K{nil}
}

// EmptyKSequence ... the KSequence with no elements
var EmptyKSequence = KSequence(0)

// NewKSequence ... creates new KSequence instance with elements
func NewKSequence(elements []K) KSequence {
	newKsPointer := len(allKs)
	allKs = append(allKs, elements...)
	allKs = append(allKs, nil)
	return KSequence(newKsPointer)
}

// IsEmpty ... returns true if KSequence has no elements
func (k KSequence) IsEmpty() bool {
	return allKs[int(k)] == nil
}

// Get ... element at position
// Caution: no checks are performed that the position is valid
func (k KSequence) Get(position int) K {
	return allKs[int(k)+position]
}

// Length ... KSequence length
func (k KSequence) Length() int {
	length := 0
	i := int(k)
	for allKs[i] != nil {
		i++
		length++
	}
	return length
}

// ToSlice ... converts KSequence to a slice of K items
func (k KSequence) ToSlice() []K {
	ptr := int(k)
	length := k.Length()
	return allKs[ptr : ptr+length]
}

// SubSequence ... subsequence starting at position
func (k KSequence) SubSequence(startPosition int) KSequence {
	return KSequence(int(k) + startPosition)
}

// TrySplitToHeadTail ... extracts first element of a KSequence, extracts the rest, if possible
// will treat non-KSequence as if they were KSequences of length 1
func TrySplitToHeadTail(k K) (ok bool, head K, tail K) {
	if kseq, isKseq := k.(KSequence); isKseq {
		ptr := int(kseq)
		head := allKs[ptr]
		if head == nil {
			// empty KSequence, no result
			return false, NoResult, EmptyKSequence
		}

		second := allKs[ptr+1]
		if second != nil && allKs[ptr+2] == nil {
			// the KSequence has length 2
			// this case is special because here the tail is not a KSequence
			return true, head, second
		}

		return true, head, KSequence(ptr + 1)
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
			ksLength := kseqElem.Length()
			if ksLength > 0 {
				elemPointer := int(kseqElem)
				elemAsSlice := allKs[elemPointer : elemPointer+ksLength]
				newKs = append(newKs, elemAsSlice...)
			}
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
