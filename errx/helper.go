package errx

import (
	"github.com/zedisdog/brynn/i18n"
	"reflect"
	"strconv"
)

func PbMap2MapStrAny(pbMap *Map) (m map[string]any, err error) {
	m = make(map[string]any, len(pbMap.Fields))
	for key, item := range pbMap.Fields {
		m[key], err = PbValue2Any(item)
		if err != nil {
			return
		}
	}

	return
}

func PbValue2Any(pbValue *Value) (v any, err error) {
	switch x := pbValue.Kind.(type) {
	case *Value_Int32Val:
		v = x.Int32Val
	case *Value_Int64Val:
		v = x.Int64Val
	case *Value_MapVal:
		v, err = PbMap2MapStrAny(x.MapVal)
	case *Value_ListVal:
		v, err = PbList2SliceAny(x.ListVal)
	case *Value_StrVal:
		v = x.StrVal
	}
	return
}

func PbList2SliceAny(pbList *List) (s []any, err error) {
	tmp := make([]any, len(pbList.List))
	for _, item := range pbList.List {
		var res any
		res, err = PbValue2Any(item)
		if err != nil {
			break
		}
		tmp = append(tmp, res)
	}

	return
}

func Map2Pb(val any) (m map[string]*Value, err error) {
	v := reflect.ValueOf(val)
	m = make(map[string]*Value, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.Kind() != reflect.String {
			err = New(InternalError, i18n.Trans("map key can only be string"))
			return
		}
		m[key.String()], err = Any2PbValue(iter.Value().Interface())
		if err != nil {
			return
		}
	}

	return
}

func Any2PbValue(val any) (value *Value, err error) {
	value = &Value{}
	switch x := val.(type) {
	case int32:
		value.Kind = &Value_Int32Val{
			Int32Val: x,
		}
	case int64:
		value.Kind = &Value_Int64Val{
			Int64Val: x,
		}
	case string:
		value.Kind = &Value_StrVal{
			StrVal: x,
		}
	case int:
		switch strconv.IntSize {
		case 32:
			value.Kind = &Value_Int32Val{
				Int32Val: int32(x),
			}
		case 64:
			value.Kind = &Value_Int64Val{
				Int64Val: int64(x),
			}
		}
	default:
		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Map:
			var m map[string]*Value
			m, err = Map2Pb(v.Interface())
			if err != nil {
				return
			}
			value.Kind = &Value_MapVal{
				MapVal: &Map{
					Fields: m,
				},
			}
		case reflect.Slice:
			s := &Value_ListVal{
				ListVal: &List{
					List: make([]*Value, 0, v.Len()),
				},
			}
			for i := 0; i < v.Len(); i++ {
				item := v.Index(i)
				var pbValue *Value
				pbValue, err = Any2PbValue(item.Interface())
				if err != nil {
					return
				}
				s.ListVal.List = append(s.ListVal.List, pbValue)
			}
			value.Kind = s
		}
	}
	return
}
