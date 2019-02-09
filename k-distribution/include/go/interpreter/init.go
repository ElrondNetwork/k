package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

func trySplitToHeadTail(k m.K) (ok bool, head m.K, tail m.KSequence) {
	if kseq, isKseq := k.(m.KSequence); isKseq {
		switch len(kseq.Ks) {
		case 0:
			return false, m.NoResult, m.EmptyKSequence
		case 1:
			return true, kseq.Ks[0], m.EmptyKSequence
		default:
			return true, kseq.Ks[0], m.KSequence{Ks: kseq.Ks[1:]}
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, k, m.EmptyKSequence
}

func assembleFromHeadAndTail(head m.K, tail m.K) m.K {
	if kseqTail, isKseq := tail.(m.KSequence); isKseq {
		if kseqTail.IsEmpty() {
			// output the head itself instead of a KSequence with 1 element
			return head
		}
		return m.KSequence{Ks: append([]m.K{head}, kseqTail.Ks...)}
	}

	// tail is not KSequence, so we end up with a KSequence of 2 elements: head and tail
	return m.KSequence{Ks: []m.K{head, tail}}
}

func assembleFromHeadSliceAndTail(headSlice []m.K, tail m.K) m.K {
	if kseqTail, isKseq := tail.(m.KSequence); isKseq {
		if kseqTail.IsEmpty() {
			// output the head itself instead of a KSequence with 1 element
			return m.KSequence{Ks: headSlice}
		}
		return m.KSequence{Ks: append(headSlice, kseqTail.Ks...)}
	}

	// tail is not KSequence, so we end up with a KSequence of 2 elements: head and tail
	return m.KSequence{Ks: append(headSlice, tail)}
}

var freshCounter int

func isTrue(c m.K) bool {
	if b, typeOk := c.(m.Bool); typeOk {
		return bool(b)
	}
	return false
}

// helps us deal with unused variables in some situations
func doNothing(c m.K) {
}
