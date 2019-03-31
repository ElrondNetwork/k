package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"testing"
)

func TestBytesEmpty(t *testing.T) {
	var bs m.K
	var err error
	bs, err = bytesHooks.empty(m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)
}

func TestBytes2Int(t *testing.T) {
	var res m.K
	var err error

	// unsigned
	res, err = bytesHooks.bytes2int(
		&m.Bytes{Value: []byte{1, 0}},
		&m.KApply{Label: m.LblBigEndianBytes},
		&m.KApply{Label: m.LblUnsignedBytes}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "256", res, err)

	res, err = bytesHooks.bytes2int(
		&m.Bytes{Value: []byte{1, 0}},
		&m.KApply{Label: m.LblLittleEndianBytes},
		&m.KApply{Label: m.LblUnsignedBytes}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "1", res, err)

	res, err = bytesHooks.bytes2int(
		&m.Bytes{Value: []byte{255}},
		&m.KApply{Label: m.LblBigEndianBytes},
		&m.KApply{Label: m.LblUnsignedBytes}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "255", res, err)

	// zero
	for _, b := range [][]byte{[]byte{}, []byte{0, 0}, []byte{0, 0, 0, 0, 0, 0}} {
		for _, c2 := range []m.K{&m.KApply{Label: m.LblBigEndianBytes}, &m.KApply{Label: m.LblLittleEndianBytes}} {
			for _, c3 := range []m.K{&m.KApply{Label: m.LblUnsignedBytes}, &m.KApply{Label: m.LblSignedBytes}} {
				res, err = bytesHooks.bytes2int(&m.Bytes{Value: b}, c2, c3, m.LblDummy, m.SortString, m.InternedBottom)
				assertIntOk(t, "0", res, err)
			}
		}
	}

	// -1
	for _, b := range [][]byte{[]byte{255}, []byte{255, 255, 255, 255, 255}} {
		for _, c2 := range []m.K{&m.KApply{Label: m.LblBigEndianBytes}, &m.KApply{Label: m.LblLittleEndianBytes}} {
			c3 := &m.KApply{Label: m.LblSignedBytes}
			res, err = bytesHooks.bytes2int(&m.Bytes{Value: b}, c2, c3, m.LblDummy, m.SortString, m.InternedBottom)
			assertIntOk(t, "-1", res, err)
		}
	}

	// other signed negative
	res, err = bytesHooks.bytes2int(
		&m.Bytes{Value: []byte{255, 254}},
		&m.KApply{Label: m.LblBigEndianBytes},
		&m.KApply{Label: m.LblSignedBytes}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "-2", res, err)

	res, err = bytesHooks.bytes2int(
		&m.Bytes{Value: []byte{255, 254}},
		&m.KApply{Label: m.LblLittleEndianBytes},
		&m.KApply{Label: m.LblSignedBytes}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "-257", res, err)
}

func TestInt2Bytes(t *testing.T) {
	var bs m.K
	var err error
	kappBigEndian := &m.KApply{Label: m.LblBigEndianBytes}

	// length 0, empty result
	bs, err = bytesHooks.int2bytes(m.IntZero, m.IntZero, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntZero, m.NewIntFromInt(12345), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntZero, m.NewIntFromInt(-12345), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	// one byte
	bs, err = bytesHooks.int2bytes(m.IntOne, m.IntZero, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.IntOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.IntMinusOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.NewIntFromInt(256), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.NewIntFromInt(257), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.NewIntFromInt(-256), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.NewIntFromInt(-257), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255}, bs, err)

	// 2 bytes
	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.IntZero, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0, 0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.IntOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0, 1}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.IntMinusOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255, 255}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(256), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(257), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 1}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(-256), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255, 0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(-257), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{254, 255}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(-255), kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255, 1}, bs, err)

	// more bytes
	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(5), m.IntZero, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0, 0, 0, 0, 0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(4), m.IntOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0, 0, 0, 1}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(6), m.IntMinusOne, kappBigEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255, 255, 255, 255, 255, 255}, bs, err)

	// little endian
	kappLittleEndian := &m.KApply{Label: m.LblLittleEndianBytes}
	bs, err = bytesHooks.int2bytes(m.IntZero, m.NewIntFromInt(-12345), kappLittleEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.int2bytes(m.IntOne, m.IntMinusOne, kappLittleEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{255}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.IntOne, kappLittleEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 0}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(2), m.NewIntFromInt(-256), kappLittleEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{0, 255}, bs, err)

	bs, err = bytesHooks.int2bytes(m.NewIntFromInt(4), m.IntOne, kappLittleEndian, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 0, 0, 0}, bs, err)
}

func TestBytesSubstr(t *testing.T) {
	var bs m.K
	var err error

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(0), m.NewIntFromInt(5), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(3), m.NewIntFromInt(3), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(0), m.NewIntFromInt(0), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(5), m.NewIntFromInt(5), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(0), m.NewIntFromInt(2), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(4), m.NewIntFromInt(5), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{5}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(1), m.NewIntFromInt(5), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.substr(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(1), m.NewIntFromInt(4), m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{2, 3, 4}, bs, err)
}

func TestBytesReplaceAt(t *testing.T) {
	var bs m.K
	var err error

	bs, err = bytesHooks.replaceAt(
		&m.Bytes{Value: []byte{10, 20, 30, 40, 50}},
		m.NewIntFromInt(0),
		&m.Bytes{Value: []byte{11, 21, 31, 41, 51}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{11, 21, 31, 41, 51}, bs, err)

	bs, err = bytesHooks.replaceAt(
		&m.Bytes{Value: []byte{10, 20, 30, 40, 50}},
		m.NewIntFromInt(2),
		&m.Bytes{Value: []byte{33, 34, 35}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{10, 20, 33, 34, 35}, bs, err)

	bs, err = bytesHooks.replaceAt(
		m.BytesEmpty,
		m.NewIntFromInt(0),
		m.BytesEmpty, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.replaceAt(
		&m.Bytes{Value: []byte{10, 20, 30}},
		m.NewIntFromInt(1),
		&m.Bytes{Value: []byte{100}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{10, 100, 30}, bs, err)
}

func TestBytesLength(t *testing.T) {
	var res m.K
	var err error

	res, err = bytesHooks.length(m.BytesEmpty, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "0", res, err)

	res, err = bytesHooks.length(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "5", res, err)
}

func TestBytesPadRight(t *testing.T) {
	var bs m.K
	var err error
	padChar := m.NewIntFromInt(80)

	bs, err = bytesHooks.padRight(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(3), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.padRight(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(5), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.padRight(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(7), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5, 80, 80}, bs, err)

	bs, err = bytesHooks.padRight(m.BytesEmpty, m.NewIntFromInt(3), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{80, 80, 80}, bs, err)
}

func TestBytesPadLeft(t *testing.T) {
	var bs m.K
	var err error
	padChar := m.NewIntFromInt(80)

	bs, err = bytesHooks.padLeft(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(3), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.padLeft(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(5), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.padLeft(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.NewIntFromInt(7), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{80, 80, 1, 2, 3, 4, 5}, bs, err)

	bs, err = bytesHooks.padLeft(m.BytesEmpty, m.NewIntFromInt(3), padChar, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{80, 80, 80}, bs, err)
}

func TestBytesReverse(t *testing.T) {
	var bs m.K
	var err error

	bs, err = bytesHooks.reverse(m.BytesEmpty, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.reverse(&m.Bytes{Value: []byte{1, 2, 3, 4, 5}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{5, 4, 3, 2, 1}, bs, err)
}

func TestBytesConcat(t *testing.T) {
	var bs m.K
	var err error

	bs, err = bytesHooks.concat(m.BytesEmpty, m.BytesEmpty, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{}, bs, err)

	bs, err = bytesHooks.concat(&m.Bytes{Value: []byte{1, 2, 3}}, &m.Bytes{Value: []byte{4, 5}}, m.LblDummy, m.SortString, m.InternedBottom)
	assertBytesOk(t, []byte{1, 2, 3, 4, 5}, bs, err)
}

func assertBytesOk(t *testing.T, expectedBytes []byte, actual m.K, err error) {
	if err != nil {
		t.Error(err)
	}
	expected := &m.Bytes{Value: expectedBytes}
	if !expected.Equals(actual) {
		t.Errorf("Unexpected Bytes. Got: %s Want: %s.",
			m.PrettyPrint(actual),
			m.PrettyPrint(expected))
	}
}
