package rdb

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	CacheClient      *redis.Client
	UserClient       redis.UniversalClient
	BizClient        redis.UniversalClient
	RdbUtxoClient    redis.UniversalClient
	RdbAddressClient redis.UniversalClient

	ctx = context.Background()
)

func init() {
	CacheClient = InitClient("conf/cache.yaml")

	BizClient = Init("conf/redis.yaml")
	RdbUtxoClient = Init("conf/rdb_utxo.yaml")
	RdbAddressClient = Init("conf/rdb_address.yaml")
	UserClient = Init("conf/user.yaml")
}

func InitClient(filename string) (rds *redis.Client) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	addr := viper.GetString("addr")
	password := viper.GetString("password")
	database := viper.GetInt("database")
	dialTimeout := viper.GetDuration("dialTimeout")
	readTimeout := viper.GetDuration("readTimeout")
	writeTimeout := viper.GetDuration("writeTimeout")
	idleTimeout := viper.GetDuration("idleTimeout")
	idleCheckFrequency := viper.GetDuration("idleCheckFrequency")
	poolSize := viper.GetInt("poolSize")
	rds = redis.NewClient(&redis.Options{
		Addr:               addr,
		Password:           password,
		DB:                 database,
		DialTimeout:        dialTimeout,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		PoolSize:           poolSize,
		IdleTimeout:        idleTimeout,
		IdleCheckFrequency: idleCheckFrequency,
	})
	return rds
}

func Init(filename string) (rds redis.UniversalClient) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	addrs := viper.GetStringSlice("addrs")
	password := viper.GetString("password")
	database := viper.GetInt("database")
	dialTimeout := viper.GetDuration("dialTimeout")
	readTimeout := viper.GetDuration("readTimeout")
	writeTimeout := viper.GetDuration("writeTimeout")
	idleTimeout := viper.GetDuration("idleTimeout")
	idleCheckFrequency := viper.GetDuration("idleCheckFrequency")
	poolSize := viper.GetInt("poolSize")
	rds = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              addrs,
		Password:           password,
		DB:                 database,
		DialTimeout:        dialTimeout,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		PoolSize:           poolSize,
		IdleTimeout:        idleTimeout,
		IdleCheckFrequency: idleCheckFrequency,
	})
	return rds
}
