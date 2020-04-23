package rel

import (
	"fmt"
	"math"
	"reflect"
)

type Iota struct {
	start, end, step int
	inclusive        bool
	resultWrapper    func(...Value) Set
	strFormat        string
}

func NewIota(
	rangeValues *rangeDataValues,
	resultWrapper func(...Value) Set,
	stringFormat string,
) Expr {
	step := int(rangeValues.step.(Number))

	start, end := int(math.Inf(1)), int(math.Inf(-1))
	if step < 0 {
		start, end = end, start
	}

	if rangeValues.start != nil {
		start = int(rangeValues.start.(Number))
	}

	if rangeValues.end != nil {
		end = int(rangeValues.end.(Number))
	}

	if rangeValues.isInclusive() {
		if step > 0 && !math.IsInf(float64(end), 1) {
			end++
		} else if step < 0 && !math.IsInf(float64(end), -1) {
			end--
		}
	}

	return Iota{start, end, step, rangeValues.isInclusive(), resultWrapper, stringFormat}
}

func (it Iota) Count() int {
	return rangeLength(it.start, it.end, it.step)
}

func (it Iota) Has(v Value) bool {
	//TODO:
	return true
}

func (it Iota) Enumerator() ValueEnumerator {
	//TODO:
	return nil
}

func (it Iota) With(v Value) Set {
	// TODO:
	return it
}

func (it Iota) Without(v Value) Set {
	// TODO:
	return it
}

func (it Iota) Map(f func(Value) Value) Set {
	//TODO:
	return it
}

func (it Iota) Where(f func(Value) bool) Set {
	//TODO:
	return it
}

func (it Iota) Call(arg Value) Value {
	i := int(arg.(Number))
	if i >= it.Count() || i < 0 {
		// TODO:
		return nil
	}
	return Number(it.start + it.step*i)
}

func getIndex(start, step, maxLen int)

func (it Iota) CallSlice(start, end Value, step int, inclusive bool) Set {
	//TODO:
	return it
}

func (it Iota) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	//TODO:
	return nil, false
}

type IotaArrayEnumerator struct {
	index, length int
}

func (ia *IotaArrayEnumerator) MoveNext() bool {
	ia.index++
	return ia.index < ia.length
}

func (ia *IotaArrayEnumerator) Current() Value {
	return
}

func (it Iota) Eval(local Scope) (Value, error) {
	// data, err := it.dataRange.eval(local)
	// if err != nil {
	// 	return nil, err
	// }
	// values := generateSequence(data[0], data[1], int(data[2].(Number)), it.dataRange.isInclusive())
	// return it.resultWrapper(values...), nil
	return it, nil
}

var iotaKind = registerKind(210, reflect.TypeOf(Function{}))

// Kind returns a number that is unique for each major kind of Value.
func (it Iota) Kind() int {
	return iotaKind
}

func (it Iota) Less(v Value) bool {
	if it.Kind() != v.Kind() {
		return it.Kind() < v.Kind()
	}
	return it.String() < v.(Iota).String()
}

func (it Iota) Negate() Value {
	//TODO: negate Iota
	return nil
}

func (it Iota) Export() interface{} {
	//TODO: Export Iota
	return nil
}

func (it Iota) Equal(v interface{}) bool {
	//TODO: Equal Iota
	return false
}

func (it Iota) Hash(seed uintptr) uintptr {
	// TODO: get back on this
	return seed
}

func (it Iota) IsTrue() bool {
	if it.inclusive {
		return true
	}

	if it.start == it.end {
		return false
	}

	return isValidRange(it.start, it.end, it.step)
}

func generateSequence(start, end Value, step int, inclusive bool) []Value {
	startVal, endVal := math.Inf(1), math.Inf(-1)
	if step < 0 {
		startVal, endVal = endVal, startVal
	}

	if start != nil {
		startVal = float64(start.(Number))
	}

	if end != nil {
		endVal = float64(end.(Number))
	}

	vals := getIndexes(int(startVal), int(endVal), step, inclusive)
	wrappedVals := make([]Value, 0, len(vals))
	for _, v := range vals {
		wrappedVals = append(wrappedVals, Number(v))
	}
	return wrappedVals
}

func (it Iota) String() string {
	var start, end Expr
	step := Number(1)
	if !isInf(it.start) {
		start = NewNumber(float64(it.start))
	}
	if !isInf(it.end) {
		end = NewNumber(float64(it.end))
	}
	if it.step != 1 {
		step = NewNumber(float64(it.step))
	}
	return fmt.Sprintf(it.strFormat, NewRangeData(start, end, step, it.inclusive).string())
}

func isInf(x int) bool {
	return math.IsInf(float64(x), 1) || math.IsInf(float64(x), -1)
}
