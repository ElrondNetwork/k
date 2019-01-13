package main

type ioHooksType int

const ioHooks ioHooksType = 0

func (ioHooksType) close(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) getc(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) open(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) putc(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) read(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) seek(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) seekEnd(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) tell(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) write(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) lock(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) unlock(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) log(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) stat(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) lstat(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) opendir(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) parse(c1 K, c2 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) parseInModule(c1 K, c2 K, c3 K, lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

func (ioHooksType) system(c K,lbl KLabel, sort Sort, config K) K {
	panic("Not implemented")
}

