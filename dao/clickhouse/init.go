package clickhouse

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/spf13/viper"
)

var (
	CK *clickhImpl
)

func init() {
	viper.SetConfigFile("conf/db.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	address := viper.GetString("address")

	config := map[string]string{
		"username":                 viper.GetString("username"),
		"password":                 viper.GetString("password"),
		"database":                 viper.GetString("database"),
		"read_timeout":             strconv.Itoa(viper.GetInt("read_timeout")),
		"write_timeout":            strconv.Itoa(viper.GetInt("write_timeout")),
		"no_delay":                 fmt.Sprintf("%t", viper.GetBool("no_delay")),
		"connection_open_strategy": viper.GetString("connection_open_strategy"),
		"block_size":               strconv.Itoa(viper.GetInt("block_size")),
		"pool_size":                strconv.Itoa(viper.GetInt("pool_size")),
		"debug":                    fmt.Sprintf("%t", viper.GetBool("debug")),
	}
	maxIdleConns := viper.GetInt("maxIdleConns")
	maxOpenConns := viper.GetInt("maxOpenConns")
	connMaxLifetime := viper.GetDuration("connMaxLifetime")

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
