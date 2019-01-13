package main

type bufferHooksType int

const bufferHooks bufferHooksType = 0

func (bufferHooksType) empty(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (bufferHooksType) concat(buf K, elem K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (bufferHooksType) toString(buf K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

