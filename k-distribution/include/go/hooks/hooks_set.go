package main

type setHooksType int

const setHooks setHooksType = 0

func (setHooksType) in(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) unit(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) element(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) concat(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) difference(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) inclusion(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) intersection(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) choice(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) size(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) set2list(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (setHooksType) list2set(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

