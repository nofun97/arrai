package rel

import "fmt"

type IotaExpr struct {
	dataRange           *RangeData
	evalWrapper         func(...Value) Set
	iotaType, strFormat string
}

func NewIotaExpr(dataRange *RangeData, iotaType string) IotaExpr {
	var evalWrapper func(...Value) Set
	var strFormat string
	switch iotaType {
	case "iota_set":
		evalWrapper = NewSet
		strFormat = "{%s}"
	case "iota_array":
		evalWrapper = NewArray
		strFormat = "{%s}"
	default:
		panic("parsing iota of unknown type")
	}
	return IotaExpr{dataRange, evalWrapper, strFormat, iotaType}
}

func (ie IotaExpr) Eval(local Scope) (Value, error) {
	data, err := ie.dataRange.eval(local)
	if err != nil {
		return nil, err
	}
	return NewIota(data, ie.evalWrapper, ie.strFormat), nil
}

func (ie IotaExpr) String() string {
	return fmt.Sprintf(ie.strFormat, ie.dataRange.string())
}
