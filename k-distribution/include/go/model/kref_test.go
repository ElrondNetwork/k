%COMMENT%

package %PACKAGE%

import "testing"

func TestKrefBasic(t *testing.T) {
	testKrefBasic(t, bottomRef, true, 0)
	testKrefBasic(t, floatRef, false, 0)

	testKrefBasic(t, bottomRef, true, refBasicDataMask)
	testKrefBasic(t, boolRef, false, refBasicDataMask)
}

func testKrefBasic(t *testing.T, refType kreferenceType, constant bool, rest uint64) {
	ref := createKrefBasic(refType, constant, rest)
	decodedType, decodedConstant, decodedRest := parseKrefBasic(ref)
	if decodedType != refType {
		t.Error("testKrefBasic mismatch")
	}
	if decodedConstant != constant {
		t.Error("testKrefBasic mismatch")
	}
	if decodedRest != rest {
		t.Error("testKrefBasic mismatch")
	}
}

func TestKrefBigInt(t *testing.T) {
	testKrefBigInt(t, true, 0, 0)
	testKrefBigInt(t, false, 0, 0)

	testKrefBigInt(t, true, 100, 50)
	testKrefBigInt(t, false, 1000, 3)
}

func testKrefBigInt(t *testing.T, constant bool, recycleCount uint64, index uint64) {
	ref := createKrefBigInt(constant, recycleCount, index)
	isBigInt, constantOut, recycleCountOut, indexOut := parseKrefBigInt(ref)
	if !isBigInt {
		t.Error("testKrefBigInt bad refType")
	}
	if constantOut != constant {
		t.Error("testKrefBigInt mismatch")
	}
	if recycleCountOut != recycleCount {
		t.Error("testKrefBigInt mismatch")
	}
	if indexOut != index {
		t.Error("testKrefBigInt mismatch")
	}
}

func TestKrefCollection(t *testing.T) {
	testKrefCollection(t, listRef, Sort(5), KLabel(7), 123)
	testKrefList(t, Sort(2), KLabel(4))
}

func testKrefCollection(t *testing.T, refType kreferenceType, sort Sort, label KLabel, index uint64) {
	ms := NewModel()
	ms.NewList(sort, label, nil)
	ref := createKrefCollection(refType, sort, label, index)
	refTypeOut, sortOut, labelOut, indexOut := parseKrefCollection(ref)
	if refTypeOut != refType {
		t.Error("testKrefCollection bad refType")
	}
	if sortOut != sort {
		t.Error("testKrefCollection mismatch")
	}
	if labelOut != label {
		t.Error("testKrefCollection mismatch")
	}
	if indexOut != index {
		t.Error("testKrefCollection mismatch")
	}
}

func testKrefList(t *testing.T, sort Sort, label KLabel) {
	ms := NewModel()
	ref := ms.NewList(sort, label, nil)
	refTypeOut, sortOut, labelOut, _ := parseKrefCollection(ref)
	if refTypeOut != listRef {
		t.Error("testKrefList bad refType")
	}
	if sortOut != sort {
		t.Error("testKrefList mismatch")
	}
	if labelOut != label {
		t.Error("testKrefList mismatch")
	}
}
