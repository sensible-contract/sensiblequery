package service

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	viper.SetConfigFile("conf/redis.yaml")
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
	poolSize := viper.GetInt("poolSize")
	rdb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        addrs,
		Password:     password,
		DB:           database,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
	})
}