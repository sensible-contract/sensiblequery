package midware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sensiblequery/dao/rdb"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SignSha256(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha256.New, keyForSign)
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func VerifyTsWithTs(ts string, expired time.Duration) bool {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return false
	}
	t := time.Unix(timestamp, 0)
	now := time.Now().UTC()
	if t.After(now.Add(expired)) || t.Before(now.Add(-expired)) {
		return false
	}
	return true
}

func VerifySignatureForHttpGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		params := c.Request.URL.Query()

		ts := params.Get("ts")
		if ok := VerifyTsWithTs(ts, time.Minute*5); !ok {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "request expired"})
			c.Abort()
			return
		}

		appid := params.Get("appid")
		secretKey, err := rdb.UserClient.Get(ctx, "secretkey:"+appid).Result()
		if err != nil {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "appid not valid"})
			c.Abort()
			return
		}

		sign := params.Get("sign")
		params.Del("sign")
		paramsStrIgnoreSign := params.Encode()

		path := c.Request.URL.Path
		signCorrect := SignSha256(path+"?"+paramsStrIgnoreSign, secretKey)
		if signCorrect != sign {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "signature not match"})
			c.Abort()
			return
		}

		quota, err := rdb.UserClient.Get(ctx, "quota:"+appid).Int()
		if err != nil {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "quota unavilable"})
			c.Abort()
			return
		}

		if quota <= 0 {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "quota exhausted"})
			c.Abort()
			return
		}

		rdb.UserClient.Decr(ctx, "quota:"+appid)
		rdb.UserClient.Incr(ctx, "visit:"+appid)
		c.Next()
	}
}