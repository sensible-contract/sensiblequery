package clickhouse

import (
	"database/sql"
)

const InitialCapacity = 256

type ScanRowFunc func(row *sql.Rows) (interface{}, error)
type ScanRowsFunc func(rows *sql.Rows) (interface{}, error)

type Operation interface {
	// 用户自定义解析过程
	Scan(psql string, srf ScanRowsFunc, args ...interface{}) (ret interface{}, err error)
	// 根据第一条数据反射结果, 要求首条数据结果不能为nil.
	ScanAll(psql string, srf ScanRowFunc, args ...interface{}) (ret interface{}, err error)
	ScanOne2(psql string, ret interface{}, args ...interface{}) (ok bool, err error)
	ScanOne(psql string, srf ScanRowFunc, args ...interface{}) (ret interface{}, err error)
	ScanRange(psql string, srf ScanRowFunc, offset int, limit int, args ...interface{}) (ret interface{}, err error)
	ScanPage(psql string, srf ScanRowFunc, offset int, limit int, sort string, desc bool, args ...interface{}) (tot int, ret interface{}, err error)
	scanPageTotal(psql string, meta *SqlMeta, args ...interface{}) (ret int, err error)

	Exec(psql string, args ...interface{}) (ret sql.Result, err error)
	ExecBatch(psql string, argsList ...interface{}) (retList []sql.Result, err error)
}

type Clickhouse interface {
	Operation
}

func Scan(psql string, srf ScanRowsFunc, args ...interface{}) (ret interface{}, err error) {
	return CK.Scan(psql, srf, args...)
}

func ScanAll(psql string, srf ScanRowFunc, args ...interface{}) (ret interface{}, err error) {
	return CK.ScanAll(psql, srf, args...)
}

func ScanOne2(psql string, ret interface{}, args ...interface{}) (ok bool, err error) {
	return CK.ScanOne2(psql, ret, args...)
}

func ScanOne(psql string, srf ScanRowFunc, args ...interface{}) (ret interface{}, err error) {
	return CK.ScanOne(psql, srf, args...)
}

func ScanRange(psql string, srf ScanRowFunc, offset int, limit int, args ...interface{}) (ret interface{}, err error) {
	return CK.ScanRange(psql, srf, offset, limit, args...)
}

func ScanPage(psql string, srf ScanRowFunc, offset int, limit int, sort string, desc bool, args ...interface{}) (tot int, ret interface{}, err error) {
	return CK.ScanPage(psql, srf, offset, limit, sort, desc, args...)
}

func Exec(psql string, args ...interface{}) (ret sql.Result, err error) {
	return CK.Exec(psql, args...)
}

func ExecBatch(psql string, argsList ...interface{}) (retList []sql.Result, err error) {
	return CK.ExecBatch(psql, argsList...)
}
