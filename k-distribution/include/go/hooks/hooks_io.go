%COMMENT%

package %PACKAGE_INTERPRETER%

import (
	m "%INCLUDE_MODEL%"
)

type ioHooksType int

const ioHooks ioHooksType = 0

func (ioHooksType) close(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) getc(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) open(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) putc(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) read(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) seek(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) seekEnd(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) tell(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) write(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) lock(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) unlock(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) log(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) stat(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) lstat(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) opendir(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) parse(c1 m.KReference, c2 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) parseInModule(c1 m.KReference, c2 m.KReference, c3 m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

func (ioHooksType) system(c m.KReference,lbl m.KLabel, sort m.Sort, config m.KReference, interpreter *Interpreter) (m.KReference, error) {
	return m.NoResult, &hookNotImplementedError{}
}

