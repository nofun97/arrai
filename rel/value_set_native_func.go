package rel

import (
	"encoding/binary"
	"fmt"

	"github.com/OneOfOne/xxhash"
	"github.com/go-errors/errors"
)

// NativeFunction represents a binary relation uniquely mapping inputs to outputs.
type NativeFunction struct {
	name string
	fn   func(Value) Value
}

// NewNativeFunction returns a new function.
func NewNativeFunction(name string, fn func(Value) Value) Expr {
	return &NativeFunction{name, fn}
}

// Name returns a native function's name.
func (f *NativeFunction) Name() string {
	return f.name
}

// Fn returns a native function's implementation.
func (f *NativeFunction) Fn() func(Value) Value {
	return f.fn
}

// Hash computes a hash for a NativeFunction.
func (f *NativeFunction) Hash(seed uint32) uint32 {
	xx := xxhash.NewS32(seed + 0x48acc265)
	binary.Write(xx, binary.LittleEndian, f.String())
	return xx.Sum32()
}

// Equal tests two Values for equality. Any other type returns false.
func (f *NativeFunction) Equal(i interface{}) bool {
	if g, ok := i.(*NativeFunction); ok {
		return f == g
	}
	return false
}

// String returns a string representation of the expression.
func (f *NativeFunction) String() string {
	return fmt.Sprintf("%s", f.name)
}

// Eval returns the Value
func (f *NativeFunction) Eval(local, global *Scope) (Value, error) {
	return f, nil
}

// Kind returns a number that is unique for each major kind of Value.
func (f *NativeFunction) Kind() int {
	return 203
}

// Bool always returns true.
func (f *NativeFunction) Bool() bool {
	return true
}

// Less returns true iff g is not a number or f.number < g.number.
func (f *NativeFunction) Less(g Value) bool {
	if f.Kind() != g.Kind() {
		return f.Kind() < g.Kind()
	}
	return f.String() < g.String()
}

// Negate returns {(negateTag): f}.
func (f *NativeFunction) Negate() Value {
	return NewTuple(NewAttr(negateTag, f))
}

// Export exports a NativeFunction.
func (f *NativeFunction) Export() interface{} {
	return f.fn
}

// Call calls the NativeFunction with the given parameter.
func (f *NativeFunction) Call(expr Expr, local, global *Scope) (Value, error) {
	if expr == nil {
		return nil, errors.Errorf("missing function arg")
	}
	value, err := expr.Eval(local, global)
	if err != nil {
		return nil, err
	}
	return f.fn(value), nil
}