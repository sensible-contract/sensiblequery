

DBHOST=192.168.31.236

# 预先创建17张表，见sql定义文件。

clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/main.sql
clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/tx.sql
clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/txin.sql
clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/txout.sql
clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/utxo.sql

# 并导入相应全量数据到：blk_height、blktx_height、txin、txout、txin_full。

clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO blk_height FORMAT RowBinary"   < out/blk.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO blktx_height FORMAT RowBinary" < out/tx.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO txin FORMAT RowBinary"         < out/txin.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO txout FORMAT RowBinary"        < out/txout.ch
clickhouse-client -h $DBHOST --database="bsv" --query="INSERT INTO txin_full FORMAT RowBinary"    < out/txin_full.ch


# 再执行以下预处理语句：

clickhouse-client -h $DBHOST --database="bsv" --multiquery < sql/all_process.sql
