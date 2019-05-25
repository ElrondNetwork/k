%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"testing"
)

func TestAssembleKSequenceEmpty1(t *testing.T) {
	kseq := m.AssembleKSequence()
	assertKSequenceEmpty(t, kseq)
}

func TestAssembleKSequenceEmpty2(t *testing.T) {
	kseq := m.AssembleKSequence(m.EmptyKSequence)
	assertKSequenceEmpty(t, kseq)
}

func TestAssembleKSequenceEmpty3(t *testing.T) {
	kseq := m.AssembleKSequence(m.EmptyKSequence, m.EmptyKSequence)
	assertKSequenceEmpty(t, kseq)
}

func TestAssembleKSequenceEmpty4(t *testing.T) {
	kseq := m.AssembleKSequence(m.EmptyKSequence, m.EmptyKSequence, m.EmptyKSequence, m.EmptyKSequence, m.EmptyKSequence)
	assertKSequenceEmpty(t, kseq)
}

func TestAssembleKSequence1(t *testing.T) {
	kseq := m.AssembleKSequence(m.NewIntFromInt(1), m.NewIntFromInt(2))
	assertKSequenceOfInts(t, kseq, 1, 2)
}

func TestAssembleKSequence2(t *testing.T) {
	kseq := m.AssembleKSequence(m.EmptyKSequence, m.NewIntFromInt(1), m.NewIntFromInt(2), m.EmptyKSequence)
	assertKSequenceOfInts(t, kseq, 1, 2)
}

func TestAssembleKSequence3(t *testing.T) {
	kseq := m.AssembleKSequence(m.NewIntFromInt(1))
	assertIntOk(t, "1", kseq, nil)
}

func TestAssembleKSequence4(t *testing.T) {
	kseq := m.AssembleKSequence(m.NewIntFromInt(1), m.EmptyKSequence)
	assertIntOk(t, "1", kseq, nil)
}

func TestAssembleKSequence5(t *testing.T) {
	kseq := m.AssembleKSequence(m.EmptyKSequence, m.NewIntFromInt(1))
	assertIntOk(t, "1", kseq, nil)
}

func TestAssembleKSequenceNest1(t *testing.T) {
	kseq1 := m.AssembleKSequence(m.NewIntFromInt(1), m.NewIntFromInt(2))
	kseq2 := m.AssembleKSequence(m.NewIntFromInt(3), m.NewIntFromInt(4))
	kseq3 := m.AssembleKSequence(kseq1, kseq2)
	assertKSequenceOfInts(t, kseq3, 1, 2, 3, 4)
}

func TestAssembleKSequenceNest2(t *testing.T) {
	kseq1 := m.AssembleKSequence(m.NewIntFromInt(1), m.NewIntFromInt(2), m.EmptyKSequence)
	kseq2 := m.AssembleKSequence(m.NewIntFromInt(3), m.EmptyKSequence, m.NewIntFromInt(4))
	kseq3 := m.AssembleKSequence(kseq1, m.EmptyKSequence, kseq2)
	assertKSequenceOfInts(t, kseq3, 1, 2, 3, 4)
}

func TestAssembleKSequenceNest3(t *testing.T) {
	kseq1 := m.AssembleKSequence(m.EmptyKSequence)
	kseq2 := m.AssembleKSequence(m.NewIntFromInt(3), m.EmptyKSequence, m.NewIntFromInt(4))
	kseq3 := m.AssembleKSequence(kseq1, m.EmptyKSequence, kseq2)
	assertKSequenceOfInts(t, kseq3, 3, 4)
}

func assertKSequenceEmpty(t *testing.T, actual m.K) {
	expected := m.EmptyKSequence
	if !actual.Equals(expected) {
		t.Errorf("Unexpected result. Got:%s Want:%s", m.PrettyPrint(actual), m.PrettyPrint(expected))
	}
}

func TestSubSequence1(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := kseq.SubSequence(0)
	assertKSequenceOfInts(t, sub, 1, 2)
}

func TestSubSequence2(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := kseq.SubSequence(1)
	assertKSequenceOfInts(t, sub, 2)
}

func TestSubSequence3(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := kseq.SubSequence(2)
	assertKSequenceEmpty(t, sub)
}

func TestSubSequenceAssemble1(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := m.AssembleKSequence(kseq.SubSequence(0))
	assertKSequenceOfInts(t, sub, 1, 2)
}

func TestSubSequenceAssemble2(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := m.AssembleKSequence(kseq.Get(0), kseq.SubSequence(1))
	assertKSequenceOfInts(t, sub, 1, 2)
}

func TestSubSequenceAssemble3(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})
	sub := m.AssembleKSequence(kseq.Get(0), kseq.Get(1), kseq.SubSequence(2))
	assertKSequenceOfInts(t, sub, 1, 2)
}

func TestSubSequenceEmpty(t *testing.T) {
	sub := m.EmptyKSequence.SubSequence(0)
	assertKSequenceEmpty(t, sub)

	// sub = m.EmptyKSequence.SubSequence(1)
	// assertKSequenceEmpty(t, sub)
}

func TestKSequenceLength(t *testing.T) {
	kseq := m.NewKSequence([]m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)})

	if kseq.Length() != 2 {
		t.Errorf("Unexpected result length. Got:%d Want:%d", kseq.Length(), 2)
	}
}

func assertKSequenceOfInts(t *testing.T, actual m.K, ints ...int) {
	var ks []m.K
	for _, i := range ints {
		ks = append(ks, m.NewIntFromInt(i))
	}
	expected := m.NewKSequence(ks)
	if !actual.Equals(expected) {
		t.Errorf("Unexpected result. Got:%s Want:%s", m.PrettyPrint(actual), m.PrettyPrint(expected))
	}
}
