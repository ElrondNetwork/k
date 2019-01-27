package %PACKAGE_INTERPRETER%

type kreflectionHooksType int

const kreflectionHooks kreflectionHooksType = 0

func (kreflectionHooksType) sort(c K, lbl KLabel, sort Sort, config K) (K, error) {
	switch k := c.(type) {
	case KToken:
		return String(k.Sort.name()), nil
	case Int:
		return String("Int"), nil
	case String:
		return String("String"), nil
	case Bytes:
		return String("Bytes"), nil
	case Bool:
		return String("Bool"), nil
	case Map:
		return String(k.Sort.name()), nil
	case List:
		return String(k.Sort.name()), nil
	case Set:
		return String(k.Sort.name()), nil
	default:
		return noResult, &hookNotImplementedError{}
	}
}

func (kreflectionHooksType) getKLabel(c K, lbl KLabel, sort Sort, config K) (K, error) {
	if k, t := c.(KApply); t {
		return InjectedKLabel{Label: k.Label}, nil
	}
	return internedBottom, nil
}

func (kreflectionHooksType) configuration(lbl KLabel, sort Sort, config K) (K, error) {
	return config, nil
}

func (kreflectionHooksType) fresh(c K, lbl KLabel, sort Sort, config K) (K, error) {
	if k, t := c.(String); t {
		sort := parseSort(string(k))
		result, err := freshFunction(sort, config, freshCounter)
		if err != nil {
		    return noResult, err
		}
		freshCounter++
		return result, nil
	}
	return noResult, &hookNotImplementedError{}
}

func (kreflectionHooksType) isConcrete(c1 K, c2 K, lbl KLabel, sort Sort, config K) (K, error) {
	return Bool(true), nil
}

func (kreflectionHooksType) getenv(c K, lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}

func (kreflectionHooksType) argv(lbl KLabel, sort Sort, config K) (K, error) {
	return noResult, &hookNotImplementedError{}
}
