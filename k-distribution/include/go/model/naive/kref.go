%COMMENT%

package %PACKAGE%

// KReference is a reference to a K item.
// In this implementation, it is just an alias to the regular pointer to the item.
type KReference = K

// NullReference is the zero-value of KReference. It doesn't point to anything.
var NullReference KReference = nil

// NoResult is the result when a function returns an error
var NoResult = InternedBottom
