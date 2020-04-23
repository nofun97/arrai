package rel

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-errors/errors"
)

type RangeData struct {
	start, end, step Expr
	inclusive        bool
}

type rangeDataValues struct {
	start, end, step Value
	inclusive        bool
}

func NewRangeData(start, end, step Expr, inclusive bool) *RangeData {
	return &RangeData{start, end, step, inclusive}
}

func (r *RangeData) eval(local Scope) (*rangeDataValues, error) {
	var start, end, step Value
	var err error

	if r.start != nil {
		start, err = r.start.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := start.(Number); !isNumber {
			return nil, errors.Errorf("lower bound does not evaluate to a Number: %s", start)
		}
	}

	if r.end != nil {
		end, err = r.end.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := end.(Number); !isNumber {
			return nil, errors.Errorf("upper bound does not evaluate to a Number: %s", end)
		}
	}

	if r.step != nil {
		step, err = r.step.Eval(local)
		if err != nil {
			return nil, err
		}
		if _, isNumber := step.(Number); !isNumber {
			return nil, errors.Errorf("step does not evaluate to a Number: %s", step)
		}
	} else {
		step = Number(1)
	}

	return &rangeDataValues{start, end, step, r.inclusive}, nil
}

func (rd *rangeDataValues) isInclusive() bool {
	return rd.inclusive
}

func (r *RangeData) string() string {
	str := strings.Builder{}
	switch {
	case r.start == nil && r.end == nil:
		str.WriteString(";")
	case r.start != nil && r.end == nil:
		str.WriteString(fmt.Sprintf("%s;", r.start))
	case r.start == nil && r.end != nil:
		str.WriteString(fmt.Sprintf(";%s", r.end))
	default:
		str.WriteString(fmt.Sprintf("%s;%s", r.start, r.end))
	}
	if r.step != nil {
		str.WriteString(fmt.Sprintf(";%s", r.step))
	}
	return str.String()
}

// resolveArrayIndexes returns an array of indexes to get for array.
func resolveArrayIndexes(start, end Value, step, offset, maxLen int, inclusive bool) []int {
	if maxLen == 0 {
		return []int{}
	}
	startIndex, endIndex := initDefaultArrayIndex(start, end, offset, maxLen+offset, step)

	if startIndex == endIndex {
		if inclusive {
			return []int{startIndex}
		}
		return []int{}
	}

	return getIndexes(startIndex, endIndex, step, inclusive)
}

// initDefaultArrayIndex initialize the start and end values for arrays.
func initDefaultArrayIndex(start, end Value, minLen, maxLen, step int) (startIndex int, endIndex int) {
	if start != nil {
		startIndex = resolveIndex(int(start.(Number)), minLen, maxLen)
		if startIndex == maxLen {
			startIndex--
		}
	} else {
		if step > 0 {
			startIndex = minLen
		} else {
			startIndex = maxLen - 1
		}
	}

	if end != nil {
		endIndex = resolveIndex(int(end.(Number)), minLen, maxLen)
	} else {
		// TODO: apply inclusivity to the undefined end index
		if step > 0 {
			endIndex = maxLen
		} else {
			endIndex = minLen - 1
		}
	}
	return
}

// resolveIndex solves the edge cases of index values.
func resolveIndex(i, minVal, maxVal int) int {
	if i > maxVal {
		return maxVal
	} else if i < 0 {
		if -i > maxVal {
			return minVal
		}
		return maxVal + i
	}
	return i
}

// getIndexes returns a range of numbers between start and end with the provided step.
// inclusive decides whether end can be included or not.
func getIndexes(start, end, step int, inclusive bool) []int {
	if !isValidRange(start, end, step) {
		return []int{}
	}
	if inclusive {
		if step > 0 {
			end++
		} else {
			end--
		}
	}

	length := int(math.Abs(float64(start - end)))
	if step != 1 && step != -1 {
		length = int(math.Ceil(float64(length) / math.Abs(float64(step))))
	}
	indexes := make([]int, 0, length)
	for i := 0; i < length; i++ {
		indexes = append(indexes, start+step*i)
	}

	return indexes
}

// isValidRange checks whether start, end, and step are valid values.
func isValidRange(start, end, step int) bool {
	return step != 0 && ((start > end && step < 0) || (start < end && step > 0))
}
