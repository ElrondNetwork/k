package %PACKAGE_INTERPRETER%

type arrayHooksType int

const arrayHooks arrayHooksType = 0

func (arrayHooksType) make(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) makeEmpty(len K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) lookup(lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) remove(arr K, v K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) update(arr K, index K, newVal K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) updateAll(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) fill(c1 K, c2 K, c3 K, c4 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) inKeys(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (arrayHooksType) ctor(lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}
