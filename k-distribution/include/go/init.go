package main

var internedBottom K = Bottom{}

var noResult K = Bottom{}

var emptyKSequence = KSequence{ks: nil}

func trySplitToHeadTail(k K) (ok bool, head K, tail KSequence) {
	if kseq, isKseq := k.(KSequence); isKseq {
		switch len(kseq.ks) {
		case 0:
			return false, noResult, emptyKSequence
		case 1:
			return true, kseq.ks[0], emptyKSequence
		default:
			return true, kseq.ks[0], KSequence{ks: kseq.ks[1:]}
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, k, emptyKSequence
}

func assembleFromHeadTail(head K, tail KSequence) K {
	if tail.isEmpty() {
		// output the element intself instead of a KSequence with 1 element
		return head
	}
	return KSequence{ks: append([]K{head}, tail.ks...)}
}

var freshCounter int

func isTrue(c K) bool {
	if b, typeOk := c.(Bool); typeOk {
		return bool(b)
	}
	return false
}
