package midware

import (
	"bytes"
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	isMetricsOn = true
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Tag  string `json:"tag,omitempty"`
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func getUrlPattern(urlPattern string, params gin.Params) string {
	numParams := len(params)
	if numParams == 0 {
      return urlPattern
   }
   for i := numParams; i > 0; i-- {
      p := params[i-1]

      idx := strings.LastIndex(urlPattern, "/"+p.Value)
      if idx < 0 {
         continue
      }
      urlPattern = urlPattern[:idx] + strings.Replace(urlPattern[idx:], "/"+p.Value, "/:"+p.Key, 1)
   }
	return urlPattern
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health/") || strings.HasPrefix(path, "/metrics") {
			// 不记录统计数据
			c.Next() // Process request
			return
		}

		start := time.Now() // Start timer
		wrapWriter := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = wrapWriter // duplicate response body

		ServiceMetrics.Inc("flight", "in") // 统计正在运行的接口
		defer ServiceMetrics.Dec("flight", "in")

		c.Next() // Process request

		// get response info
		var bizCode int
		var responseObj Response
		responseBodyJson := wrapWriter.body.Bytes()
		if e := json.Unmarshal(responseBodyJson, &responseObj); e == nil {
			bizCode = responseObj.Code
		}

		latency := time.Now().Sub(start)
		statusCode := c.Writer.Status()

		// 跟新监控
		if isMetricsOn {
			ServiceMetrics.Observe("all", "visit", latency)                          // 对所有进来的请求统计计数
			ServiceMetrics.Observe(strconv.Itoa(statusCode), "httpcode", latency)  // 对所有进来的请求按 Http Code 进行统计计数
			ServiceMetrics.Observe(strconv.Itoa(bizCode), "bizcode", latency) // 对所有进来的请求按 Biz Code 进行统计计数
			if statusCode != 404 {
				urlPattern := getUrlPattern(path, c.Params)
				ServiceMetrics.Observe(urlPattern, "url", latency)                    // 对所有进来的请求按 Http URL 进行统计计数
				ServiceMetrics.Observe(strconv.Itoa(bizCode)+urlPattern, "bizcode-url", latency)
			}
		}
	}
}


func CreateMetricsEndpoint(adminGinWeb gin.IRouter) {
	adminGinWeb.GET("/metrics", fetchMetricsSummary)
	adminGinWeb.GET("/metrics/:time/:type/:stage", fetchMetrics)
}

func fetchMetricsSummary(c *gin.Context) {
	timeFilter := "minute"
	typeFilter := "url"
	stageFilter := "past"

	metrics := ServiceMetrics.Dump(timeFilter, typeFilter, stageFilter)
	c.String(http.StatusOK, metrics)
}

func fetchMetrics(c *gin.Context) {
	timeFilter := c.Param("time")
	typeFilter := c.Param("type")
	stageFilter := c.Param("stage")

	metrics := ServiceMetrics.Dump(timeFilter, typeFilter, stageFilter)
	c.String(http.StatusOK, metrics)
}
