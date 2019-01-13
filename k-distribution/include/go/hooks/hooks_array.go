package main

type arrayHooksType int

const arrayHooks arrayHooksType = 0

func (arrayHooksType) make(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) makeEmpty(len K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) lookup(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) remove(arr K, v K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) update(arr K, index K, newVal K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) updateAll(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) fill(c1 K, c2 K, c3 K, c4 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) in_keys(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (arrayHooksType) ctor(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}
