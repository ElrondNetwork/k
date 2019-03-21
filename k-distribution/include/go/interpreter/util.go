package %PACKAGE_INTERPRETER%

import (
    "fmt"
	m "%INCLUDE_MODEL%"
)

func trySplitToHeadTail(k m.K) (ok bool, head m.K, tail m.K) {
	if kseq, isKseq := k.(*m.KSequence); isKseq {
		switch len(kseq.Ks) {
		case 0:
			return false, m.NoResult, m.EmptyKSequence
		case 1:
			return true, kseq.Ks[0], m.EmptyKSequence
		case 2:
		    return true, kseq.Ks[0], kseq.Ks[1]
		default:
			return true, kseq.Ks[0], &m.KSequence{Ks: kseq.Ks[1:]}
		}
	}

	// treat non-KSequences as if they were KSequences with 1 element
	return true, k, m.EmptyKSequence
}

// appends all elements to a KSequence
// flattens if there are any KSequences among the elements (but only on 1 level, does not handle multiple nesting)
// never returns KSequence of 1 element, it returns the element directly instead
func assembleKSequence(elements ...m.K) m.K {
	var newKs []m.K
	for _, element := range elements {
		if kseqElem, isKseq := element.(*m.KSequence); isKseq {
			newKs = append(newKs, kseqElem.Ks...)
		} else {
			newKs = append(newKs, element)
		}
	}
	if len(newKs) == 0 {
		return m.EmptyKSequence
	}
	if len(newKs) == 1 {
		return newKs[0]
	}
	return &m.KSequence{Ks: newKs}
}

var freshCounter int

func isTrue(c m.K) bool {
	if b, typeOk := c.(*m.Bool); typeOk {
		return b.Value
	}
	return false
}

// helps us deal with unused variables in some situations
func doNothing(c m.K) {
}

// can be handy when debugging
func debugPrint(c m.K) {
	fmt.Println(m.PrettyPrint(c))
}
