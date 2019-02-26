package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"testing"
)

func TestStringConcat(t *testing.T) {
	result, err := stringHooks.concat(m.NewString("abc"), m.NewString("def"), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "abcdef", result, err)
}

func TestStringEq(t *testing.T) {
	var result m.K
	var err error

	result, err = stringHooks.eq(m.NewString("abc"), m.NewString("abc"), m.LblDummy, m.SortString, m.InternedBottom)
	assertBoolOk(t, true, result, err)

	result, err = stringHooks.ne(m.NewString("abc"), m.NewString("abc"), m.LblDummy, m.SortString, m.InternedBottom)
	assertBoolOk(t, false, result, err)

	result, err = stringHooks.eq(m.NewString("yes"), m.NewString("no"), m.LblDummy, m.SortString, m.InternedBottom)
	assertBoolOk(t, false, result, err)

	result, err = stringHooks.ne(m.NewString(""), m.NewString("s"), m.LblDummy, m.SortString, m.InternedBottom)
	assertBoolOk(t, true, result, err)
}

func TestStringChr(t *testing.T) {
	var result m.K
	var err error

	result, err = stringHooks.chr(m.NewIntFromInt(97), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "a", result, err)

	result, err = stringHooks.chr(m.NewIntFromInt(32), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, " ", result, err)
}

func TestStringFind(t *testing.T) {
	var result m.K
	var err error

	str := m.NewString("abcabcabcd")
	substr := m.NewString("abc")

	result, err = stringHooks.find(str, substr, m.NewIntFromInt(0), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "0", result, err)

	result, err = stringHooks.find(str, substr, m.NewIntFromInt(1), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "3", result, err)

	result, err = stringHooks.find(str, substr, m.NewIntFromInt(3), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "3", result, err)

	result, err = stringHooks.find(str, substr, m.NewIntFromInt(7), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "-1", result, err)

	result, err = stringHooks.rfind(str, substr, m.NewIntFromInt(10), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "6", result, err)

	result, err = stringHooks.rfind(str, substr, m.NewIntFromInt(6), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "3", result, err)

	result, err = stringHooks.rfind(str, substr, m.NewIntFromInt(2), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "-1", result, err)

	result, err = stringHooks.rfind(str, substr, m.NewIntFromInt(0), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "-1", result, err)
}

func TestStringLength(t *testing.T) {
	len, err := stringHooks.length(m.NewString("abc"), m.LblDummy, m.SortString, m.InternedBottom)
	assertIntOk(t, "3", len, err)
}

func TestStringSubstr(t *testing.T) {
	var result m.K
	var err error

	str := m.NewString("abcdef")

	result, err = stringHooks.substr(str, m.NewIntFromInt(0), m.NewIntFromInt(2), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "ab", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(0), m.NewIntFromInt(6), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "abcdef", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(0), m.NewIntFromInt(1000), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "abcdef", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(2), m.NewIntFromInt(3), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "c", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(2), m.NewIntFromInt(2), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(0), m.NewIntFromInt(0), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "", result, err)

	result, err = stringHooks.substr(str, m.NewIntFromInt(6), m.NewIntFromInt(6), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "", result, err)
}

func TestString2Token(t *testing.T) {
	result, err := stringHooks.string2token(m.NewString("abc"), m.LblDummy, m.SortString, m.InternedBottom)
	if err != nil {
		t.Error(err)
	}
	expected := m.KToken{Sort: m.SortString, Value: "abc"}
	if result != expected {
		t.Errorf("Wrong KToken. Got: %s Want: %s.",
			result.PrettyTreePrint(0),
			expected.PrettyTreePrint(0))
	}
}

func TestToken2String(t *testing.T) {
	var result m.K
	var err error

	ktoken := m.KToken{Sort: m.SortKResult, Value: "token!"}
	result, err = stringHooks.token2string(ktoken, m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "token!", result, err)

	result, err = stringHooks.token2string(m.NewIntFromInt(56), m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "56", result, err)

	result, err = stringHooks.token2string(m.BoolTrue, m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "true", result, err)

	result, err = stringHooks.token2string(m.BoolFalse, m.LblDummy, m.SortString, m.InternedBottom)
	assertStringOk(t, "false", result, err)

}

func assertStringOk(t *testing.T, expectedStr string, actual m.K, err error) {
	if err != nil {
		t.Error(err)
	}
	k, isString := actual.(m.String)
	if !isString {
		t.Error("Result is not a String.")
		return
	}
	if expectedStr != k.String() {
		t.Errorf("Unexpected String. Got: %s Want: %s.",
			k.String(),
			expectedStr)
	}
}
