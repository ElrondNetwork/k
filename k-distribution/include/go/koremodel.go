package main

// K ... Defines a K entity, this is either a KItem, or a KSequence of KItems
type K interface {
}

type KItem interface {
}

type KApply struct {
	label string
	list  []K
}

type InjectedKLabel struct {
	label string
}

type KToken struct {
	value string
	sort  string
}

type KVariable struct {
	name string
}

// KSequence ... a sequence of K items
type KSequence struct {
	ks []K
}
