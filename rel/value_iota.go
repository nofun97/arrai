package rel

import (
	"fmt"
	"math"
	"reflect"
)

type Iota struct {
	start, end, step, length int
	inclusive                bool
	resultWrapper            func(...Value) Set
	strFormat                string
}

func NewIota(
	rangeValues *rangeDataValues,
	resultWrapper func(...Value) Set,
	stringFormat string,
) Set {
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

	if !isValidRange(start, end, step) {
		return None
	}

	length := rangeLength(start, end, step)
	if isInf(end) && rangeValues.isInclusive() {
		length++
	}
	return Iota{start, end, step, length, rangeValues.isInclusive(), resultWrapper, stringFormat}
}

func (it Iota) Count() int {
	return it.length
}

func (it Iota) Has(v Value) bool {
	if _, isNumber := v.(Number); !isNumber {
		return false
	}
	index := int(math.Abs(float64((int(v.(Number)) - it.start) / it.step)))
	return index < it.length
}

func (it Iota) Enumerator() ValueEnumerator {
	if e, enumeratable := it.ArrayEnumerator(); enumeratable {
		return e
	}
	return nil
}

func (it Iota) With(v Value) Set {
	//TODO: how to handle lazy unbounded set
	if it.Has(v) {
		return it
	}
	if _, isNumber := v.(Number); !isNumber || it.inRange(v) {
		return newSetFromSet(it).With(v)
	}
	num := int(v.(Number))
	if it.step > 0 {
		var rv *rangeDataValues
		if it.getMaxNum() + it.step == num {
			rv = &rangeDataValues{
				Number(it.start), Number(it.)
			}
		}
	}
	return newSetFromSet(it).With(v)
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
	return getValueAtIndex(int(arg.(Number)), it.start, it.step, it.Count())
}

func getValueAtIndex(index, start, step, maxLen int) Value {
	if index >= maxLen || index < 0 {
		// TODO:
		return nil
	}
	return Number(start + step*index)
}

func (it Iota) CallSlice(start, end Value, step int, inclusive bool) Set {
	//TODO:
	return it
}

func (it Iota) ArrayEnumerator() (OffsetValueEnumerator, bool) {
	if it.length == 0 {
		return nil, false
	}
	return &IotaArrayEnumerator{0, it.start, it.step, it.length}, true
}

type IotaArrayEnumerator struct {
	index, minVal, increment, length int
}

func (ia *IotaArrayEnumerator) MoveNext() bool {
	ia.index++
	return ia.index < ia.length
}

func (ia *IotaArrayEnumerator) Current() Value {
	return getValueAtIndex(ia.index, ia.minVal, ia.increment, ia.length)
}

func (ia *IotaArrayEnumerator) Offset() int {
	//TODO:
	return 0
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

// func generateSequence(start, end Value, step int, inclusive bool) []Value {
// 	startVal, endVal := math.Inf(1), math.Inf(-1)
// 	if step < 0 {
// 		startVal, endVal = endVal, startVal
// 	}

// 	if start != nil {
// 		startVal = float64(start.(Number))
// 	}

// 	if end != nil {
// 		endVal = float64(end.(Number))
// 	}

// 	vals := getIndexes(int(startVal), int(endVal), step, inclusive)
// 	wrappedVals := make([]Value, 0, len(vals))
// 	for _, v := range vals {
// 		wrappedVals = append(wrappedVals, Number(v))
// 	}
// 	return wrappedVals
// }

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

func (it Iota) getMaxNum() int {
	return int(getValueAtIndex(it.Count()-1, it.start, it.step, it.Count()).(Number))
}

func (it Iota) inRange(n Value) bool {
	if !isValidRange(it.start, it.end, it.step) {
		return false
	}
	num := int(n.(Number))
	maxNum := it.getMaxNum()
	if it.step > 0 {
		return it.start <= num && num < maxNum
	}
	return it.start >= num && num < maxNum
}

func isInf(x int) bool {
	return math.IsInf(float64(x), 1) || math.IsInf(float64(x), -1)
}
