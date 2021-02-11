package clickhouse

import (
	"bytes"
	"fmt"
	"sync"
	"unicode"
)

type Map struct {
	Data map[string]*SqlMeta
	Mutx *sync.RWMutex
}

var sqlMetaMap = &Map{
	Data: make(map[string]*SqlMeta),
	Mutx: new(sync.RWMutex),
}

var (
	SELECT   = []rune("SELECT")
	DISTINCT = []rune("DISTINCT")
	FROM     = []rune("FROM")
	WHERE    = []rune("WHERE")
	GROUP    = []rune("GROUP")
	ORDER    = []rune("ORDER")
	LIMIT    = []rune("LIMIT")

	LEN_SELECT    = len(SELECT)
	LEN_DINSTINCT = len(DISTINCT)
	LEN_FROM      = len(FROM)
	LEN_WHERE     = len(WHERE)
	LEN_GROUP     = len(GROUP)
	LEN_ORDER     = len(ORDER)
	LEN_LIMIT     = len(LIMIT)
)

const (
	SPACE     = '\x20'
	UPPE_DIFF = 'A' - 'a'
)

/*
清除psql中的空白
*/
func TWS(psql string) (ret string) {
	bd := new(bytes.Buffer)
	start := 0
	end := 0
	ps := []rune(psql)
	len := len(ps)
	for end < len {
		start = indexOf(IsNWS, ps, end, len)
		if start == -1 {
			break
		}
		end = indexOf(IsWS, ps, start, len)
		if end == -1 {
			end = len
		}
		if bd.Len() > 0 {
			bd.WriteByte(SPACE)
		}
		for i := start; i < end; i++ {
			bd.WriteRune(ps[i])
		}
	}
	ret = bd.String()
	return
}
func IsNWS(ch rune) bool {
	return !unicode.IsSpace(ch)
}

func IsWS(ch rune) bool {
	return unicode.IsSpace(ch)
}

func indexOf(mf func(ch rune) bool, ps []rune, start int, len int) int {
	for start < len {
		switch ch := ps[start]; ch {
		case '\'':
			start++
			for start < len {
				ch = ps[start]
				if ch == '\'' {
					nxt := start + 1
					if nxt < len && ps[nxt] == '\'' {
						start++
					} else {
						break
					}
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '"':
			start++
			for start < len {
				ch = ps[start]
				if ch == '"' {
					break
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '`':
			start++
			for start < len {
				ch = ps[start]
				if ch == '`' {
					nxt := start + 1
					if nxt < len && ps[nxt] == '`' {
						start++
					} else {
						break
					}
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '/':
			ch = ps[start+1]
			if ch == '*' {
				start += 3 //match "*/"
				for start < len {
					if ps[start] == '/' && ps[start-1] == '*' {
						break
					}
					start++
				}
			}
		case '#':
			start++
			for start < len {
				ch = ps[start]
				if ch == '\n' { //不需要/r
					break
				}
				start++
			}
		case '-':
			ch = ps[start+1]
			if ch == '-' {
				start++
				for start < len {
					ch = ps[start]
					if ch == '\n' {
						break
					}
					start++
				}
			}
		default:
			if mf(ch) {
				return start
			}
		}
		start++
	}
	return -1
}

func indexOfIncludeParent(mf func(ch rune) bool, ps []rune, start int, len int) int {
	left := 0
	for start < len {
		switch ch := ps[start]; ch {
		case '\'':
			start++
			for start < len {
				ch = ps[start]
				if ch == '\'' {
					nxt := start + 1
					if nxt < len && ps[nxt] == '\'' {
						start++
					} else {
						break
					}
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '"':
			start++
			for start < len {
				ch = ps[start]
				if ch == '"' {
					break
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '`':
			start++
			for start < len {
				ch = ps[start]
				if ch == '`' {
					nxt := start + 1
					if nxt < len && ps[nxt] == '`' {
						start++
					} else {
						break
					}
				} else if ch == '\\' {
					start++
				}
				start++
			}
		case '/':
			ch = ps[start+1]
			if ch == '*' {
				start += 3 //match "*/"
				for start < len {
					if ps[start] == '/' && ps[start-1] == '*' {
						break
					}
					start++
				}
			}
		case '#':
			start++
			for start < len {
				ch = ps[start]
				if ch == '\n' { //不需要/r
					break
				}
				start++
			}
		case '-':
			ch = ps[start+1]
			if ch == '-' {
				start++
				for start < len {
					ch = ps[start]
					if ch == '\n' {
						break
					}
					start++
				}
			}
		case '(':
			left++
		case ')':
			left--
		default:
			if left == 0 && mf(ch) {
				return start
			}
		}
		start++
	}
	return -1
}

type SqlMeta struct {
	Select   int
	Distinct int
	From     int
	Where    int
	Group    int
	Order    int
	Limit    int

	Mutx          *sync.Mutex
	LimitPsql     string   //带上limit的SQL缓存
	TotalPsql     string   //带上count(*)的SQL缓存
	LimitPsqlMeta *SqlMeta //LimitPsql的元数据,用于分页显示
}

func NewSqlMeta() (ret *SqlMeta) {
	return &SqlMeta{Select: -1, Distinct: -1, From: -1, Where: -1, Group: -1, Order: -1, Limit: -1, Mutx: new(sync.Mutex)}
}

func IsIdentifier(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '_'
}

/*
Mysql:
1. ', ", `, (, 需要查找配对
2. 在',",`内遇到\自动跳过下一个字符
*/
func ParseSqlMeta(psql string) (ret *SqlMeta) {
	ret = NewSqlMeta()

	start := 0
	end := 0
	ps := []rune(psql)
	len := len(ps)
	for end < len {
		start = indexOfIncludeParent(IsNWS, ps, end, len)
		if start == -1 {
			break
		}
		end = indexOfIncludeParent(IsWS, ps, start, len)
		if end == -1 {
			end = len
		}

		n := end - start
		switch {
		case n == LEN_SELECT && regionMatchesIgnoreCase(SELECT, 0, ps, start, n):
			ret.Select = start
		case n == LEN_DINSTINCT && regionMatchesIgnoreCase(DISTINCT, 0, ps, start, n):
			ret.Distinct = start
		case n == LEN_FROM && regionMatchesIgnoreCase(FROM, 0, ps, start, n):
			ret.From = start
		case n == LEN_WHERE && regionMatchesIgnoreCase(WHERE, 0, ps, start, n):
			ret.Where = start
		case n == LEN_GROUP && regionMatchesIgnoreCase(GROUP, 0, ps, start, n):
			ret.Group = start
		case n == LEN_ORDER && regionMatchesIgnoreCase(ORDER, 0, ps, start, n):
			ret.Order = start
		case n == LEN_LIMIT && regionMatchesIgnoreCase(LIMIT, 0, ps, start, n):
			ret.Limit = start
		}
	}
	return
}

func regionMatchesIgnoreCase(ps []rune, start1 int, ns []rune, start2 int, len int) bool {
	for i := 0; i < len; i++ {
		if diff := ps[start1] - ns[start2]; diff != 0 && diff != UPPE_DIFF {
			return false
		}
		start1++
		start2++
	}
	return true
}

func GetSqlMeta(psql string) (ret *SqlMeta) {
	sqlMetaMap.Mutx.RLock()
	vl, ok := sqlMetaMap.Data[psql]
	sqlMetaMap.Mutx.RUnlock()
	if !ok {
		vl = ParseSqlMeta(psql)
		sqlMetaMap.Mutx.Lock()
		sqlMetaMap.Data[psql] = vl
		sqlMetaMap.Mutx.Unlock()
	}
	ret = vl
	return
}

func GenLimitSql(psql string, meta *SqlMeta) {
	meta.Mutx.Lock()
	if meta.LimitPsql == "" {
		if meta.Limit > 0 {
			// 存在limit, 需用子查询select * from(...) limit ?,?
			meta.LimitPsql = fmt.Sprintf("SELECT * FROM ( %s ) LIMIT ?,?", psql)
		} else {
			// 不存limit, 直接后拼
			meta.LimitPsql = psql + ` LIMIT ?,?`
		}
	}
	meta.Mutx.Unlock()
}

func GenTotalSql(psql string, meta *SqlMeta) {
	meta.Mutx.Lock()
	if meta.TotalPsql == "" {
		bd := new(bytes.Buffer)
		// 如果存在DISTINCT, GROUP, 或LIMIT, 需用子查询 select count(*) from(...)
		// 否则直接select count(*) <from_clause> <where_clause> <group_clause> <limit_clause>, 不再需要order子句
		if meta.Distinct > 0 || meta.Group > 0 || meta.Limit > 0 {
			bd.WriteString(`SELECT COUNT(*) FROM ( `)
			if meta.Order > 0 {
				bd.WriteString(psql[0:meta.Order])
				if meta.Limit > 0 {
					bd.WriteString(psql[meta.Limit:])
				}
			} else {
				// 没有order
				bd.WriteString(psql)
			}
			bd.WriteString(` )`)
		} else {
			bd.WriteString(`SELECT COUNT(*) `)
			if meta.Order > 0 {
				bd.WriteString(psql[meta.From:meta.Order])
				if meta.Limit > 0 {
					bd.WriteString(psql[meta.Limit:])
				}
			} else {
				// 没有order
				bd.WriteString(psql[meta.From:])
			}
		}
		meta.TotalPsql = bd.String()
	}
	meta.Mutx.Unlock()
}

func GenDataSql(psql string, meta *SqlMeta, field string, desc bool) (ret string) {
	if meta.LimitPsql == "" {
		GenLimitSql(psql, meta)
	}

	// 不需要order by直接使用LimitSql即可, 否则动态LimitSql的order by
	if field == "" {
		ret = meta.LimitPsql
	} else {
		// 注意: 肯定会有limit,因为这是LimitPsql
		if meta.LimitPsqlMeta == nil {
			meta.LimitPsqlMeta = ParseSqlMeta(meta.LimitPsql)
		}

		bd := new(bytes.Buffer)
		if meta.LimitPsqlMeta.Order > 0 {
			bd.WriteString(meta.LimitPsql[0:meta.LimitPsqlMeta.Order])
		} else {
			bd.WriteString(meta.LimitPsql[0:meta.LimitPsqlMeta.Limit])
		}
		bd.WriteString("ORDER BY `")
		bd.WriteString(field)
		bd.WriteString("` ")
		if desc {
			bd.WriteString(`DESC `)
		}
		bd.WriteString(meta.LimitPsql[meta.LimitPsqlMeta.Limit:])
		ret = bd.String()
	}
	return
}

func append2(args []interface{}, offset interface{}, limit interface{}) []interface{} {
	ln := len(args)
	nargs := make([]interface{}, ln+2)
	if ln > 0 {
		copy(nargs, args)
	}
	nargs[ln] = offset
	nargs[ln+1] = limit
	return nargs
}
