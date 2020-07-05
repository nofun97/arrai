package rel

import (
	"math"
	"strings"

	"github.com/go-errors/errors"
)

type SliceData struct {
	start, end, step Expr
	original         string
}

func NewSliceData(start, end, step Expr) *SliceData {
	// string created at the start to retain original input data
	var originalSlice strings.Builder
	if start != nil {
		originalSlice.WriteString(start.String())
	}
	originalSlice.WriteByte(':')
	if end != nil {
		originalSlice.WriteString(end.String())
	}
	if step != nil {
		originalSlice.WriteByte(':')
		originalSlice.WriteString(step.String())
	}

	return &SliceData{start, end, step, originalSlice.String()}
}

func (s *SliceData) SetStart(start Expr) {
	s.start = start
}

func (s *SliceData) SetEnd(end Expr) {
	s.end = end
}

func (s *SliceData) Eval(local Scope) error {
	eval := func(e *Expr, t string) error {
		if (*e) == nil {
			return nil
		}
		val, err := (*e).Eval(local)
		if err != nil {
			return err
		}
		if _, err = assertNumber(val, t); err != nil {
			return err
		}
		*e = val
		return nil
	}
	if err := eval(&s.start, "start"); err != nil {
		return err
	}
	if err := eval(&s.end, "end"); err != nil {
		return err
	}
	if err := eval(&s.step, "step"); err != nil {
		return err
	}
	return nil
}

func (s *SliceData) String() string {
	return s.original
}

func assertNumber(expr Expr, t string) (Number, error) {
	num, isNumber := expr.(Number)
	if !isNumber {
		return 0, errors.Errorf("slice %s does not evaluate to a number: %s", t, num)
	}
	return num, nil
}

func (s *SliceData) CreateSliceIterator() (*SliceIterator, error) {
	if s.start == nil || s.end == nil {
		panic("CreateSliceIterator: start and end of slice cannot be nil")
	}

	var start, step, end Number
	var err error
	// there's eval, is this check redundant?
	start, err = assertNumber(s.start, "start")
	if err != nil {
		return nil, err
	}

	end, err = assertNumber(s.end, "end")
	if err != nil {
		return nil, err
	}

	if s.step == nil {
		if isValidRange(start, end, 1) {
			s.step = NewNumber(1)
		} else {
			s.step = NewNumber(-1)
		}
	}

	step, err = assertNumber(s.step, "step")
	if err != nil {
		return nil, err
	}

	if !isValidRange(start, end, step) {
		return nil, errors.Errorf("slice does not evaluate to a valid range: %s:%s:%s", start, end, step)
	}
	return &SliceIterator{start, end, step, 0}, nil
}

func isValidRange(start, end, step Number) bool {
	if step < 0 {
		return start >= end
	} else if step > 0 {
		return start <= end
	}
	return false
}

type SliceIterator struct {
	start, limit, increment, index Number
}

func (s *SliceIterator) Count() int {
	return int(
		math.Abs(
			math.Abs(float64(s.start)-float64(s.limit))/float64(s.increment) + 1,
		),
	)
}

func (s *SliceIterator) Current() Value {
	return s.CallIndex(s.index)
}

func (s *SliceIterator) MoveNext() bool {
	if int(s.index)+1 > s.Count() {
		return false
	}
	s.index++
	return true
}

func (s *SliceIterator) CallIndex(i Number) Number {
	n := s.start + s.increment*i
	if (s.increment < 0 && n < s.limit) || (s.increment > 0 && n > s.limit) {
		//TODO: should this return error?
		return s.limit
	} else if s.increment == 0 {
		panic("SliceIterator.CallIndex: increment cannot be 0")
	}
	return n
}

func initArraySlice(s *SliceData, offset, max int) {
	if s.start == nil {
		s.start = NewNumber(float64(offset))
	}
	if s.end == nil {
		s.end = NewNumber(float64(offset + max))
	}
}


