package %PACKAGE_INTERPRETER%

type kequalHooksType int

const kequalHooks kequalHooksType = 0

func (kequalHooksType) eq(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (kequalHooksType) ne(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (kequalHooksType) ite(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

