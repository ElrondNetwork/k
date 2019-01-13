package main

type kreflectionHooksType int

const kreflectionHooks kreflectionHooksType = 0

func (kreflectionHooksType) sort(c K, lbl KLabel, sort Sort, config K) K {
	switch k := c.(type) {
	case KToken:
		return String(k.Sort.name())
	case Int:
		return String("Int")
	case String:
		return String("String")
	case Bytes:
		return String("Bytes")
	case Bool:
		return String("Bool")
	case Map:
		return String(k.Sort.name())
	case List:
		return String(k.Sort.name())
	case Set:
		return String(k.Sort.name())
	default:
		panic("Not implemented")
	}
}

func (kreflectionHooksType) getKLabel(c K, lbl KLabel, sort Sort, config K) K {
	if k, t := c.(KApply); t {
		return InjectedKLabel{Label: k.Label}
	}
	return internedBottom
}

func (kreflectionHooksType) configuration(lbl KLabel, sort Sort, config K) K {
	return config
}

func (kreflectionHooksType) fresh(c K, lbl KLabel, sort Sort, config K) K {
	if k, t := c.(String); t {
		sort := parseSort(string(k))
		result := freshFunction(sort, config, freshCounter)
		freshCounter++
		return result
	}
	panic("Not implemented")
}

func (kreflectionHooksType) isConcrete(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	return Bool(true)
}

func (kreflectionHooksType) getenv(c K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (kreflectionHooksType) argv(lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}
