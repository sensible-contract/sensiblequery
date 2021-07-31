module satosensible

go 1.15

require (
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-contrib/cache v1.1.0
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-contrib/zap v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-redis/redis/v8 v8.8.3
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sensible-contract/sensible-script-decoder v1.9.1
	github.com/spf13/viper v1.7.1
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.6.7
	github.com/ybbus/jsonrpc/v2 v2.1.6
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/tools v0.1.0 // indirect
)

replace github.com/go-redis/redis/v8 v8.8.3 => github.com/sensible-contract/redis/v8 v8.8.3
