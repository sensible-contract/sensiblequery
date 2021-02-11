package clickhouse

import (
	"database/sql"
	"fmt"
	"time"
)

type Type uint

const (
	Bool Type = iota
	Int
	Int32
	Int64
	Float32
	Float64
	String
	Timep
	Bytes
)

func Newp(v Type) interface{} {
	switch v {
	case Bool:
		var ret bool
		return &ret
	case Int:
		var ret int
		return &ret
	case Int32:
		var ret int32
		return &ret
	case Int64:
		var ret int64
		return &ret
	case Float32:
		var ret float32
		return &ret
	case Float64:
		var ret float64
		return &ret
	case String:
		var ret string
		return &ret
	case Timep:
		ret := (*time.Time)(nil)
		return &ret
	case Bytes:
		var ret []byte
		return &ret
	default:
		panic(fmt.Errorf("newp failed for: %#v", v))
	}
}

func Extv(v interface{}) interface{} {
	switch v := v.(type) {
	case *bool:
		return *v
	case *int:
		return *v
	case *int32:
		return *v
	case *int64:
		return *v
	case *float32:
		return *v
	case *float64:
		return *v
	case *string:
		return *v
	case **time.Time:
		return *v
	case *[]byte:
		return *v
	case *interface{}:
		return *v
	default:
		panic(fmt.Errorf("extv failed for: %#v", v))
	}
}

func BoolR(rows *sql.Rows) (interface{}, error) {
	var ret bool
	err := rows.Scan(&ret)
	return ret, err
}

func IntR(rows *sql.Rows) (interface{}, error) {
	var ret int
	err := rows.Scan(&ret)
	return ret, err
}

func Int32R(rows *sql.Rows) (interface{}, error) {
	var ret int32
	err := rows.Scan(&ret)
	return ret, err
}

func Int64R(rows *sql.Rows) (interface{}, error) {
	var ret int64
	err := rows.Scan(&ret)
	return ret, err
}

func Float32R(rows *sql.Rows) (interface{}, error) {
	var ret float32
	err := rows.Scan(&ret)
	return ret, err
}

func Float64R(rows *sql.Rows) (interface{}, error) {
	var ret float64
	err := rows.Scan(&ret)
	return ret, err
}

func StringR(rows *sql.Rows) (interface{}, error) {
	var ret string
	err := rows.Scan(&ret)
	return ret, err
}

func TimepR(rows *sql.Rows) (interface{}, error) {
	ret := new(time.Time)
	err := rows.Scan(&ret)
	return ret, err
}

func SliceR(ks ...Type) ScanRowFunc {
	// 下述在整个扫描
	ln := len(ks)
	return func(rows *sql.Rows) (interface{}, error) {
		ret := make([]interface{}, ln)
		for i, k := range ks {
			ret[i] = Newp(k)
		}
		err := rows.Scan(ret...)
		if err != nil {
			return nil, err
		}
		for i := 0; i < ln; i++ {
			ret[i] = Extv(ret[i])
		}
		return ret, nil
	}
}

/*name1,type1,name2,type2...*/
func MapR(pairs ...interface{}) ScanRowFunc {
	pln := len(pairs)
	len := pln / 2

	ks := make([]string, len)
	ts := make([]Type, len)

	idx := 0
	for i := 1; i < pln; i += 2 {
		ks[idx] = pairs[i-1].(string)
		ts[idx] = pairs[i].(Type)
		idx++
	}

	return func(rows *sql.Rows) (interface{}, error) {
		vs := make([]interface{}, len)
		for i, t := range ts {
			vs[i] = Newp(t)
		}
		err := rows.Scan(vs...)
		if err != nil {
			return nil, err
		}

		ret := make(map[string]interface{}, len)
		for i, k := range ks {
			ret[k] = Extv(vs[i])
		}
		return ret, nil
	}
}
