package rel

import (
	"fmt"
	"math"
	"strings"
)

type IotaExpr struct {
	dataRange     *RangeData
	resultWrapper func(...Value) Value
	strFormat     string
}

func NewIotaExpr(
	dataRange *RangeData,
	resultWrapper func(...Value) Value,
	stringFormat string,
) Expr {
	return &IotaExpr{dataRange, resultWrapper, stringFormat}
}

func (it *IotaExpr) Eval(local Scope) (Value, error) {
	data, err := it.dataRange.eval(local)
	if err != nil {
		return nil, err
	}
	values := generateSequence(data[0], data[1], int(data[2].(Number)), it.dataRange.isInclusive())
	return it.resultWrapper(values...), nil
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

func (it *IotaExpr) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf(it.strFormat, it.dataRange.string()))
	return str.String()
}
