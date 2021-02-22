

DBHOST=192.168.31.236

# 创建：blk_height_new、blktx_height_new、txin_new、txout_new、txin_full_new。
# 新表结构如同：blk_height、blktx_height、txin、txout、txin_full。

clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/sync_init.sql

# 导入增量数据到新表

clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO blk_height_new FORMAT RowBinary"   < out/blk.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO blktx_height_new FORMAT RowBinary" < out/tx.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO txin_new FORMAT RowBinary"         < out/txin.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO txout_new FORMAT RowBinary"        < out/txout.ch


# 在更新之前，如果有上次已导入但是当前被孤立的块，需要先删除这些块的数据。直接从公有块高度（COMMON_HEIGHT）往上删除就可以了。
# 如果没有孤块，则无需处理
# clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/sync-delete_orphan.sql

# 针对已有的中间表，执行以下增量预处理语句：

clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/sync_process.sql
