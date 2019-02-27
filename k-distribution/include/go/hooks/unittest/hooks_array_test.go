package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
	"testing"
)

func TestArrayMake(t *testing.T) {
	var arr m.K
	var err error
	var int3 = m.NewIntFromInt(3)
	var int5 = m.NewIntFromInt(5)
	var bottom = m.InternedBottom

	arr, err = arrayHooks.makeEmpty(int3, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, bottom, bottom}, arr, err)

	arr, err = arrayHooks.make(int3, int5, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, int5, []m.K{int5, int5, int5}, arr, err)

	arr, err = arrayHooks.ctor(bottom, int5, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, bottom, bottom, bottom, bottom}, arr, err)

}

// Without default (default = bottom)
func TestArrayMakeUpdateRemoveLookup1(t *testing.T) {
	var arr m.K
	var err error
	var elem m.K
	var int3 = m.NewIntFromInt(3)
	var int5 = m.NewIntFromInt(5)
	var int7 = m.NewIntFromInt(7)
	var bottom = m.InternedBottom

	// create
	arr, err = arrayHooks.makeEmpty(int3, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, bottom, bottom}, arr, err)

	// updates
	arrayHooks.update(arr, m.NewIntFromInt(1), int5, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, int5, bottom}, arr, err)

	arrayHooks.update(arr, m.NewIntFromInt(2), int7, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, int5, int7}, arr, err)

	// test some lookups
	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(0), m.LblDummy, m.SortInt, m.InternedBottom)
	assertBottomOk(t, elem, err)

	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "5", elem, err)

	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(2), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "7", elem, err)

	// remove
	arrayHooks.remove(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, bottom, int7}, arr, err)

	// lookup again
	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertBottomOk(t, elem, err)
}

// With default
func TestArrayMakeUpdateRemoveLookup2(t *testing.T) {
	var arr m.K
	var err error
	var elem m.K
	var int3 = m.NewIntFromInt(3)
	var int5 = m.NewIntFromInt(5)
	var int7 = m.NewIntFromInt(7)
	var defElem = m.NewIntFromInt(20)

	// create
	arr, err = arrayHooks.make(int3, defElem, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, defElem, []m.K{defElem, defElem, defElem}, arr, err)

	// updates
	arrayHooks.update(arr, m.NewIntFromInt(1), int5, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, defElem, []m.K{defElem, int5, defElem}, arr, err)

	arrayHooks.update(arr, m.NewIntFromInt(2), int7, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, defElem, []m.K{defElem, int5, int7}, arr, err)

	// test some lookups
	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(0), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "20", elem, err)

	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "5", elem, err)

	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(2), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "7", elem, err)

	// remove
	arrayHooks.remove(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, defElem, []m.K{defElem, defElem, int7}, arr, err)

	// lookup again
	elem, err = arrayHooks.lookup(arr, m.NewIntFromInt(1), m.LblDummy, m.SortInt, m.InternedBottom)
	assertIntOk(t, "20", elem, err)
}

func TestArrayUpdateAll1(t *testing.T) {
	var arr m.K
	var err error
	var bottom = m.InternedBottom

	arr, _ = arrayHooks.makeEmpty(m.NewIntFromInt(4), m.LblDummy, m.SortInt, m.InternedBottom)
	list1 := &m.List{Sort: m.SortInt, Data: []m.K{m.NewIntFromInt(1), m.NewIntFromInt(2)}}
	arrayHooks.updateAll(arr, m.NewIntFromInt(1), list1, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, m.NewIntFromInt(1), m.NewIntFromInt(2), bottom}, arr, err)
}

func TestArrayUpdateAll2(t *testing.T) {
	var arr m.K
	var err error
	var bottom = m.InternedBottom

	arr, _ = arrayHooks.makeEmpty(m.NewIntFromInt(4), m.LblDummy, m.SortInt, m.InternedBottom)
	list2 := &m.List{Sort: m.SortInt, Data: []m.K{m.NewIntFromInt(1), m.NewIntFromInt(2), m.NewIntFromInt(3), m.NewIntFromInt(4)}}
	arrayHooks.updateAll(arr, m.NewIntFromInt(1), list2, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, m.NewIntFromInt(1), m.NewIntFromInt(2), m.NewIntFromInt(3)}, arr, err)
}

func TestArrayFill1(t *testing.T) {
	var arr m.K
	var err error
	var bottom = m.InternedBottom
	var fill = m.NewIntFromInt(123)

	arr, _ = arrayHooks.makeEmpty(m.NewIntFromInt(4), m.LblDummy, m.SortInt, m.InternedBottom)
	arrayHooks.fill(arr, m.NewIntFromInt(1), m.NewIntFromInt(3), fill, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, fill, fill, bottom}, arr, err)
}

func TestArrayFill2(t *testing.T) {
	var arr m.K
	var err error
	var bottom = m.InternedBottom
	var fill = m.NewIntFromInt(123)

	arr, _ = arrayHooks.makeEmpty(m.NewIntFromInt(4), m.LblDummy, m.SortInt, m.InternedBottom)
	arrayHooks.fill(arr, m.NewIntFromInt(1), m.NewIntFromInt(10), fill, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, bottom, []m.K{bottom, fill, fill, fill}, arr, err)
}

func TestArrayInKeys(t *testing.T) {
	var arr m.K
	var err error
	var result m.K
	var defElem = m.NewIntFromInt(20)

	// create
	arr, err = arrayHooks.make(m.NewIntFromInt(2), defElem, m.LblDummy, m.SortInt, m.InternedBottom)
	assertArrayOk(t, defElem, []m.K{defElem, defElem}, arr, err)

	// updates
	arrayHooks.update(arr, m.IntZero, m.NewIntFromInt(5), m.LblDummy, m.SortInt, m.InternedBottom)

	// test
	result, err = arrayHooks.inKeys(m.IntZero, arr, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, false, result, err)

	result, err = arrayHooks.inKeys(m.NewIntFromInt(1), arr, m.LblDummy, m.SortInt, m.InternedBottom)
	assertBoolOk(t, true, result, err)
}

func assertArrayOk(t *testing.T, expectedDefault m.K, expectedElems []m.K, a m.K, err error) {
	if err != nil {
		t.Error(err)
	}
	arr, isArray := a.(*m.Array)
	if !isArray {
		t.Error("Result is not an Array.")
		return
	}
	if !expectedDefault.Equals(arr.Default) {
		t.Errorf("Unexpected Array default. Got: %s Want: %s.",
			arr.Default.PrettyTreePrint(0),
			expectedDefault.PrettyTreePrint(0))
	}
	if len(expectedElems) != len(arr.Data) {
		t.Errorf("Unexpected Array length. Got: %d Want: %d.",
			len(arr.Data),
			len(expectedElems))
		return
	}
	for i := 0; i < len(expectedElems); i++ {
		if !arr.Data[i].Equals(expectedElems[i]) {
			t.Errorf("Unexpected element at position %d. Got: %s Want: %s.",
				i,
				arr.Data[i].PrettyTreePrint(0),
				expectedElems[i].PrettyTreePrint(0))
		}
	}
}

func assertBottomOk(t *testing.T, actual m.K, err error) {
	if err != nil {
		t.Error(err)
	}
	_, isBottom := actual.(*m.Bottom)
	if !isBottom {
		t.Errorf("Bottom expected. Got: %s", actual.PrettyTreePrint(0))
		return
	}
}
