
# 区块数据服务 Bitcoin SV blockchain API service

我们部署了一个浏览器Demo [BSV Browser](https://sensiblequery.com/#/blocks) ，可测试查看Blockchain的数据。

Api Endpoint: `https://api.sensiblequery.com`

支持的API见：`https://api.sensiblequery.com/swagger/index.html`

Testnet Api Endpoint: `https://api.sensiblequery.com/test/`

Testnet 支持的API见：`https://api.sensiblequery.com/test/swagger/index.html`

区块数据服务包括2个组件，使用clickhouse作为数据计算存储引擎，redis作为每个地址的UTXO集合的数据存储。

### 1. 节点数据同步程序：sensibled

sensibled 通过访问全节点的区块文件夹来同步区块数据(默认在`~/.bitcoin/blocks/`)，同步的数据保存在clickhouse中，UTXO信息保存在redis中。以便支持已确认区块数据查询。

同时通过监听节点zmq，实时同步获取tx内容并更新到redis、clickhouse中。以便支持tx、余额、UTXO数据的实时查询。

### 2. 数据API server：sensiblequery

查询redis、clickhouse中的数据，以对外提供数据API服务。

## sensiblequery 运行依赖

1. 需要节点提供rpc服务。以封装pushtx接口。
2. 需要与sensibled服务使用同一个redis实例，同一个clickhouse实例。以便获取数据。

## 配置文件

在conf目录有程序运行需要的多个配置文件。

* db.yaml

clickhouse数据库配置，主要包括address、database等。

* chain.yaml

节点配置，rpc地址。

* redis.yaml

redis配置，主要包括addrs、database等。

目前同时兼容redis cluster和single-node。addrs配置单个地址将视为single-node。

## 使用Docker运行

使用docker-compose可以比较方便运行sensiblequery。首先设置好db/redis/node配置，然后运行：

	$ docker-compose up -d

停止请执行：

	$ docker-compose stop


## 使用主机运行

需要设置环境变量LISTEN，以配置API服务的监听端口。然后直接启动程序即可。此时日志会直接输出到终端。

    $ LISTEN=:5555 ./sensiblequery

可使用nohup或其他技术将程序放置到后台运行。

sensiblequery服务可以随时重启，除了会中断用户访问，不会造成任何最终数据问题。


## 部署资源需求

| 部署                 | DISK   | MEM   |
|----------------------|--------|-------|
| sensiblequery        | 20 GB  | 4 GB  |
| bsv-node + sensibled | 500 GB | 16 GB |
| clickhouse           | 1.5 TB | 16 GB |
| redis x 1            | 50GB   | 32GB  |
| redis-cluster x 6    | 50GB   | 16GB  |

其中sensiblequery用来对外提供API服务，可以部署多实例。sensibled是单实例运行。redis可以部署单节点或集群。
