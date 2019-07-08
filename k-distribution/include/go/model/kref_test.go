%COMMENT%

package %PACKAGE_MODEL%

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
		t.Error("testKrefBasic mismatch")
	}
	if recycleCountOut != recycleCount {
		t.Error("testKrefBasic mismatch")
	}
	if indexOut != index {
		t.Error("testKrefBasic mismatch")
	}
}
