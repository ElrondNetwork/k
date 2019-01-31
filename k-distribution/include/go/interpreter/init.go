package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

var internedBottom m.K = m.Bottom{}

var noResult m.K = m.Bottom{}

var emptyKSequence = m.KSequence{Ks: nil}

func trySplitToHeadTail(k m.K) (ok bool, head m.K, tail m.KSequence) {
	if kseq, isKseq := k.(m.KSequence); isKseq {
		switch len(kseq.Ks) {
		case 0:
			return false, noResult, emptyKSequence
		case 1:
			return true, kseq.Ks[0], emptyKSequence
		default:
			return true, kseq.Ks[0], m.KSequence{Ks: kseq.Ks[1:]}
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, k, emptyKSequence
}

func assembleFromHeadTail(head m.K, tail m.KSequence) m.K {
	if tail.IsEmpty() {
		// output the element intself instead of a m.KSequence with 1 element
		return head
	}
	return m.KSequence{Ks: append([]m.K{head}, tail.Ks...)}
}

var freshCounter int

func isTrue(c m.K) bool {
	if b, typeOk := c.(m.Bool); typeOk {
		return bool(b)
	}
	return false
}
