package main

var internedBottom K = Bottom{}

var noResult K = Bottom{}

var freshCounter int = 0

func doNothingWithVar(v K) {
}

type hookNotImplementedError struct {
}

func (e *hookNotImplementedError) Error() string {
	return "Hook not implemented."
}
