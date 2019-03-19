package %PACKAGE_INTERPRETER%

import (
	"fmt"
    m "%INCLUDE_MODEL%"
	"testing"
)

func TestParseIntOk(t *testing.T) {
	strs := []string{
		"0", "123",
		"-123",
		"57896044618658097711785492504343953926634992332820282019728792003956564819968",
		"-57896044618658097711785492504343953926634992332820282019728792003956564819968"}
	for _, s := range strs {
		i, err := m.ParseInt(s)
		assertIntOk(t, s, i, err)
	}
}

func TestParseIntError(t *testing.T) {
	strs := []string{"abc", "-0", ""}
	for _, s := range strs {
		_, err := m.ParseInt(s)
		if err == nil {
			t.Errorf("Error expected when parsing %s", s)
		}
	}
}

func TestIntHooks1(t *testing.T) {
	a := m.NewIntFromInt(1)
	b := m.NewIntFromInt(2)
	var z m.K
	var err error

	z, err = intHooks.eq(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.ne(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.le(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.lt(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.ge(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.gt(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.add(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "3", z, err)

	z, err = intHooks.sub(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "-1", z, err)

	z, err = intHooks.mul(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", z, err)

	z, err = intHooks.tdiv(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.tmod(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.ediv(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.emod(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.shl(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "4", z, err)

	z, err = intHooks.shr(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.and(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.or(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "3", z, err)

	z, err = intHooks.xor(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "3", z, err)

	z, err = intHooks.not(b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "-3", z, err)

	z, err = intHooks.abs(b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", z, err)

	z, err = intHooks.max(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", z, err)

	z, err = intHooks.min(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

}

func TestIntHooks2(t *testing.T) {
	a := m.NewIntFromInt(1)
	b := m.NewIntFromInt(1)

	var z m.K
	var err error

	z, err = intHooks.eq(m.NewIntFromInt(1), m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.ne(m.NewIntFromInt(1), m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.le(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.lt(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.ge(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, z, err)

	z, err = intHooks.gt(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, z, err)

	z, err = intHooks.add(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", z, err)

	z, err = intHooks.sub(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.mul(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.tdiv(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.tmod(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.ediv(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.emod(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.shl(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", z, err)

	z, err = intHooks.shr(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.and(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.or(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.xor(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", z, err)

	z, err = intHooks.not(a, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "-2", z, err)

	z, err = intHooks.abs(a, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.max(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)

	z, err = intHooks.min(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", z, err)
}

func TestIntHooksPow(t *testing.T) {
	a := m.NewIntFromInt(2)
	b := m.NewIntFromInt(10)
	c := m.NewIntFromInt(1000)
	var z m.K
	var err error

	z, err = intHooks.pow(a, b, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1024", z, err)

	z, err = intHooks.pow(b, a, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "100", z, err)

	z, err = intHooks.powmod(a, b, c, m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "24", z, err)
}

func TestIntLog2(t *testing.T) {
	var log m.K
	var err error

	log, err = intHooks.log2(m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "0", log, err)

	log, err = intHooks.log2(m.NewIntFromInt(2), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", log, err)

	log, err = intHooks.log2(m.NewIntFromInt(3), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "1", log, err)

	log, err = intHooks.log2(m.NewIntFromInt(4), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "2", log, err)

	log, err = intHooks.log2(m.NewIntFromInt(255), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "7", log, err)

	log, err = intHooks.log2(m.NewIntFromInt(256), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "8", log, err)

	for i := 1000; i < 1009; i++ {
		big, _ := intHooks.shl(m.NewIntFromInt(1), m.NewIntFromInt(i), m.LblDummy, m.SortInt, m.InternedBottom)
		log, err = intHooks.log2(big, m.LblDummy, m.SortInt, m.InternedBottom)
		assertIntOk(t, fmt.Sprintf("%d", i), log, err)

		big, _ = intHooks.sub(big, m.IntOne, m.LblDummy, m.SortInt, m.InternedBottom)
		log, err = intHooks.log2(big, m.LblDummy, m.SortInt, m.InternedBottom)
		assertIntOk(t, fmt.Sprintf("%d", i-1), log, err)
	}
}

func assertIntOk(t *testing.T, expectedAsStr string, actual m.K, err error) {
	if err != nil {
		t.Error(err)
	}
	expectedK := m.NewIntFromString(expectedAsStr)
	if !actual.Equals(expectedK) {
		t.Errorf("Unexpected result. Got:%s Want:%s", m.PrettyPrint(actual), m.PrettyPrint(expectedK))
	}
}
