%COMMENT%

package %PACKAGE%

// KVariable is a KObject representing a KVariable item in K
type KVariable struct {
	Name string
}

// NewKVariable creates a new object and returns the reference.
func (ms *ModelState) NewKVariable(name string) KReference {
	return &KVariable{Name: name}
}