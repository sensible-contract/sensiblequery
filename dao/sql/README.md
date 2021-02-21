区块浏览器包括2个组件，一个是区块数据导出工具(blkparser)，一个是数据API server(satoblock)。使用clickhouse作为数据计算存储引擎。

blkparser 通过访问全节点的区块文件夹来导出区块数据(默认在`~/.bitcoin/blocks/`)，导出的数据保存在clickhouse中，以供satoblock查询使用。

后续会增加监听zmq，不完全基于区块文件导出，并使用MongoDB/Redis辅助做mempool的数据同步。

目前我们先列出手动同步区块数据的方法。


# 手动同步数据的方法

首次导出区块时执行全量导出，后续则使用增量导出来补充新区块数据。

blkparser支持指定起始/终止高度来导出一定范围内的区块数据，导出高度包括start但不包括end。全量导出时起始高度参数为0：

```
./blkparser -start 0 -end 674936
./blkparser -start 674936 -end 674940
```

导出数据文件有4种：

* blk.ch 区块头原始数据
* tx.ch 区块包括的tx原始数据
* txin.ch tx输入信息原始数据
* txout.ch tx输出信息原始数据

这些数据可以直接全量或分批导入clickhouse数据库：

```
cat /data/blk.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO blk_height FORMAT RowBinary"
cat /data/tx.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO blktx_height FORMAT RowBinary"
cat /data/txin.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO txin FORMAT RowBinary"
cat /data/txout.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO txout FORMAT RowBinary"
```

将导入的基础数据执行一些SQL预处理，生成几个具有不同索引的中间数据表，即可提供给浏览器API实现业务查询使用。

详细表定义可见SQL文件。

### 全量导入后预处理

全量导入的SQL处理语句详见`prepare_all.sql`。

其中表txin_full的全量生成对数据库压力很大，可以直接使用blkparser导出全量基础数据：

    $ ./blkparser -start 0 -end 674936 -full

将导出的txin.ch直接导入到txin_full表中即可，不需要对txin/txout表进行join运算：

* txin.ch

```
cat /data/txin.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO txin_full FORMAT RowBinary"
```
导入之后，需要通过txin_full表来创建出txin表：
```
INSERT INTO txin SELECT txid, idx, utxid, vout, script_sig, nsequence, height FROM txin_full
```

### 增量导入后预处理

执行增量导入的SQL处理语句详见`prepare_height.sql`。
