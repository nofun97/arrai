package rel

import (
	"fmt"

	"github.com/arr-ai/wbnf/parser"
	"github.com/go-errors/errors"
)

type SliceExpr struct {
	ExprScanner
	lhs   Expr
	slice *SliceData
}

func NewSliceExpr(scanner parser.Scanner, lhs Expr, slice *SliceData) Expr {
	return SliceExpr{ExprScanner{scanner}, lhs, slice}
}

func (s SliceExpr) Eval(local Scope) (Value, error) {
	set, err := s.lhs.Eval(local)
	if err != nil {
		return nil, WrapContext(err, s, local)
	}
	if _, isSet := set.(Set); !isSet {
		return nil, WrapContext(errors.Errorf("only set can be sliced: %s", set), s, local)
	}

	if err := s.slice.Eval(local); err != nil {
		return nil, WrapContext(err, s, local)
	}

	return set.(Set).CallSlice(s.slice)
}

func (s SliceExpr) String() string {
	return fmt.Sprintf("%s(%s)", s.lhs, s.slice)
}
