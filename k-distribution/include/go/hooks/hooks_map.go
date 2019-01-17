package main

type mapHooksType int

const mapHooks mapHooksType = 0

func (mapHooksType) element(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) unit(lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) concat(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) lookup(kmap K, key K, lbl KLabel, sort Sort, config K) (K, error) {
	if m, isMap := kmap.(Map); isMap {
		elem, found := m.data[key]
		if found {
			return elem, nil
		}
		return internedBottom, nil
	}

	if _, isBottom := kmap.(Map); isBottom {
		return internedBottom, nil
	}

	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) lookupOrDefault(kmap K, key K, defaultRes K, lbl KLabel, sort Sort, config K) (K, error) {
	if m, isMap := kmap.(Map); isMap {
		elem, found := m.data[key]
		if found {
			return elem, nil
		}
		return defaultRes, nil
	}

	if _, isBottom := kmap.(Map); isBottom {
		return internedBottom, nil
	}

	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) update(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) remove(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) difference(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) keys(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) keys_list(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) in_keys(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) values(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) choice(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) size(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) inclusion(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) updateAll(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (mapHooksType) removeAll(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}
