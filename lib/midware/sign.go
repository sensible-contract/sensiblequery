package midware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sensiblequery/dao/rdb"

	"github.com/gin-gonic/gin"
)

func SignSha256(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha256.New, keyForSign)
	h.Write([]byte(input))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func VerifySignatureForHttpGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		path := c.Request.URL.Path

		params := c.Request.URL.Query()
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

		signCorrect := SignSha256(path+"?"+paramsStrIgnoreSign, secretKey)
		if signCorrect != sign {
			c.JSON(http.StatusForbidden, &Response{Code: -1, Msg: "signature not match"})
			c.Abort()
			return
		}
		c.Next()
	}
}
