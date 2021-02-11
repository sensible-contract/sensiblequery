package clickhouse

import (
	"database/sql"
	"net/url"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

var (
	CK *clickhImpl
)

//     # 地址(必需). 多值用逗号分隔
//     address: "10.11.165.44:9000"
//     # DB名字(必需)
//     database: "demo"
//     # 用户名(可选)
//     username: ""
//     # 密码(可选)
//     password: ""
//     # 最大空闲数量(可选)
//     maxIdleConns:
//     # 最大打开数量(可选)
//     maxOpenConns:
//     # 最大lifetime(可选)
//     connMaxLifetime: "0s"
//     # 读超时(秒)
//     read_timeout: 10
//     # 写超时(秒)
//     write_timeout : 10
//     # 无延迟
//     no_delay: true
//     # connection_open_strategy 连接打开策略
//     connection_open_strategy: "random"
//     # 读写块大小
//     block_size: 1000000
//     # 连接池大小
//     pool_size: 100
//     # 是否debug
//     debug: false

func init() {
	config := map[string]string{
		"username":                 "",
		"password":                 "",
		"database":                 "bsv",
		"read_timeout":             "10",
		"write_timeout":            "10",
		"no_delay":                 "true",
		"connection_open_strategy": "random",
		"block_size":               "1000000",
		"pool_size":                "10",
		"debug":                    "false",
	}
	maxIdleConns := 10
	maxOpenConns := 10
	connMaxLifetime := 10 * time.Minute

	address := "192.168.31.236:9000"

	sb := new(strings.Builder)
	for key, value := range config {
		addit(sb, key, value)
	}
	db, err := sql.Open("clickhouse", "tcp://"+address+sb.String())
	if err != nil {
		panic(err)
	}
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	}
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	}
	if connMaxLifetime > 0 {
		db.SetConnMaxLifetime(connMaxLifetime)
	}

	CK = &clickhImpl{DB: db}
}

func addit(sb *strings.Builder, key, val string) {
	if strings.TrimSpace(val) == "" {
		return
	}

	if sb.Len() == 0 {
		sb.WriteByte('?')
	} else {
		sb.WriteByte('&')
	}
	sb.WriteString(key)
	sb.WriteByte('=')
	sb.WriteString(url.QueryEscape(val))
}
