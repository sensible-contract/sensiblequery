
# 区块数据服务 Bitcoin SV blockchain API service

我们部署了一个浏览器Demo [BSV Browser](https://sensiblequery.com/#/blocks) ，可测试查看Blockchain的数据。

Api Endpoint: `https://api.sensiblequery.com`

支持的API见：`https://api.sensiblequery.com/swagger/index.html`

Testnet Api Endpoint: `https://api.sensiblequery.com/test/`

Testnet 支持的API见：`https://api.sensiblequery.com/test/swagger/index.html`

区块数据服务包括3个组件，使用clickhouse作为数据计算存储引擎，redis作为每个地址的UTXO集合的数据存储。

## 1. 节点区块同步程序：satoblock

satoblock 通过访问全节点的区块文件夹来同步区块数据(默认在`~/.bitcoin/blocks/`)，同步的数据保存在clickhouse中，UTXO信息保存在redis中。可支持已确认区块数据查询。

## 2. 节点mempool实时同步程序：satomempool

satomempool 通过监听节点zmq，实时获取tx内容并更新到redis、clickhouse中。可支持tx、余额、UTXO数据的实时查询。

## 3. 数据API server：satosensible

查询redis、clickhouse中的数据，以对外提供数据API服务。

## satosensible 运行依赖

1. 需要节点提供rpc服务。以封装pushtx接口。
2. 需要与satoblock服务使用同一个redis实例，同一个clickhouse实例。以便查询数据。

## 配置文件

在conf目录有程序运行需要的多个配置文件。

* db.yaml

clickhouse数据库配置，主要包括address、database等。

* chain.yaml

节点配置，rpc地址。

* redis.yaml

redis配置，主要包括address、database等。

需要占用2个database号，database_block存放UTXO原始script，database存放UTXO集合key。需要和satomempool配置保持一致。


## 运行方式

需要设置环境变量LISTEN，以配置API服务的监听端口。然后直接启动程序即可。此时日志会直接输出到终端。

    $ LISTEN=:5555 ./satosensible

可使用nohup或其他技术将程序放置到后台运行。
