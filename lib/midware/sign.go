package midware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sensiblequery/dao/rdb"
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
	ats := []rune(ts)
	if len(ats) < 14 {
		return false
	}

	t, e := time.ParseInLocation("20060102150405.000", fmt.Sprint(string(ats[0:14])+"."+string(ats[14:])), time.UTC)
	if e != nil {
		return false
	}
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
		secretKey, err := rdb.UserClient.Get(ctx, "appid:"+appid).Result()
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
		c.Next()
	}
}
