package main

var internedBottom K = Bottom{}

var noResult K = Bottom{}

var freshCounter int = 0

func isTrue(c K) bool {
    if b, typeOk := c.(Bool); typeOk {
        return bool(b)
    }
    return false
}

type hookNotImplementedError struct {
}

func (e *hookNotImplementedError) Error() string {
	return "Hook not implemented."
}
