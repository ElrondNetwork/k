package main

type mapHooksType int

const mapHooks mapHooksType = 0

func (mapHooksType) element(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) unit(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) concat(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) lookup(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) lookupOrDefault(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) update(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) remove(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) difference(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) keys(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) keys_list(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) in_keys(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) values(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) choice(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) size(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) inclusion(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) updateAll(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (mapHooksType) removeAll(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

