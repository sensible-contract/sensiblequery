/*
预先创建并导入新表：blk_height_new、blktx_height_new、txin_new、txout_new、txin_full_new。新表只包含增量导入数据。

新表结构如同：blk_height、blktx_height、txin、txout、txin_full

cat /data/blk.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO blk_height_new FORMAT RowBinary"
cat /data/tx.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO blktx_height_new FORMAT RowBinary"
cat /data/txin.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO txin_new FORMAT RowBinary"
cat /data/txout.ch | clickhouse-client -h DBHOST --database="bsv" --query="INSERT INTO txout_new FORMAT RowBinary"

针对已有的中间表，执行以下预处理语句：
*/

-- 在更新之前，如果有上次已导出数据但是当前被孤立的块，需要先删除这些块的数据。直接从公有块高度（COMMON_HEIGHT）往上删除就可以了。
ALTER TABLE blk_height DELETE WHERE height > COMMON_HEIGHT
ALTER TABLE blk DELETE WHERE height > COMMON_HEIGHT

ALTER TABLE blktx_height DELETE WHERE height > COMMON_HEIGHT
ALTER TABLE tx DELETE WHERE height > COMMON_HEIGHT

ALTER TABLE txin DELETE WHERE height > COMMON_HEIGHT
ALTER TABLE txin_spent DELETE WHERE height > COMMON_HEIGHT

-- 回滚已被花费的utxo_address
INSERT INTO utxo
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height_txo, 1 FROM txin_full
  WHERE satoshi > 0 AND
      height > COMMON_HEIGHT;
-- 删除新添加的utxo_address˜
INSERT INTO utxo_address
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txout
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a)) AND
      height > COMMON_HEIGHT;


-- 回滚已被花费的utxo_genesis
INSERT INTO utxo_genesis
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height_txo, 1 FROM txin_full
  WHERE satoshi > 0 AND
      height > COMMON_HEIGHT;
-- 删除新添加的utxo_genesis
INSERT INTO utxo_genesis
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txout
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a)) AND
      height > COMMON_HEIGHT;


ALTER TABLE txin_full DELETE WHERE height > COMMON_HEIGHT

ALTER TABLE txout DELETE WHERE height > COMMON_HEIGHT



-- 更新现有基础数据表blk_height、blktx_height、txin、txout
INSERT INTO blk_height SELECT * FROM blk_height_new
INSERT INTO blktx_height SELECT * FROM blktx_height_new
INSERT INTO txin SELECT * FROM txin_new
INSERT INTO txout SELECT * FROM txout_new


-- 更新区块id索引
INSERT INTO blk SELECT * FROM blk_height_new;


-- 更新区块内tx索引
INSERT INTO tx SELECT * FROM blktx_height_new;
-- 更新tx到区块高度索引
INSERT INTO tx_height SELECT txid, height FROM blktx_height_new ORDER BY txid;


-- 更新txo被花费的tx索引
INSERT INTO txin_spent SELECT height, txid, idx, utxid, vout FROM txin_new;
-- 更新txo被花费的tx区块高度索引
INSERT INTO txout_spent_height SELECT height, utxid, vout FROM txin_new ORDER BY utxid;


-- 更新输入详情, 到新表txin_full_new
INSERT INTO txin_full_new
  SELECT height, txid, idx, script_sig, nsequence,
         txo.height, txo.utxid, txo.vout, txo.address, txo.genesis, txo.satoshi, txo.script_type, txo.script_pk FROM txin_new
  LEFT JOIN (
      SELECT height, utxid, vout, address, genesis, satoshi, script_type, script_pk FROM txout
      WHERE (height, utxid, vout) IN (
          SELECT height, txid, txin.vout FROM tx_height
          JOIN (
              SELECT utxid, vout FROM txin_new
          ) AS txin
          ON tx_height.txid = txin.utxid
          WHERE txid in (
              SELECT utxid FROM txin_new
          )
      )
  ) AS txo
  USING (utxid, vout)


INSERT INTO txin_full SELECT * FROM txin_full_new;


-- 更新地址参与的输出索引
INSERT INTO txout_address_height SELECT height, utxid, vout, address, genesis FROM txout_new ORDER BY address;
-- 更新溯源ID参与的输出索引
INSERT INTO txout_genesis_height SELECT height, utxid, vout, address, genesis FROM txout_new ORDER BY genesis;

-- 更新地址参与输入的相关tx区块高度索引
INSERT INTO txin_address_height SELECT height, txid, idx, address, genesis FROM txin_full_new ORDER BY address;
-- 更新溯源ID参与输入的相关tx区块高度索引
INSERT INTO txin_genesis_height SELECT height, txid, idx, address, genesis FROM txin_full_new ORDER BY genesis;


-- 更新地址相关的utxo索引
-- 增量添加utxo_address
INSERT INTO utxo_address
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM txout_new
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a));
-- 已花费txo标记清除
INSERT INTO utxo_address
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txin_new;
-- 如果一个satoshi=0的txo被花费(早期有这个现象)，就可能遗留一个sign=-1的数据，需要删除
ALTER TABLE utxo_address DELETE WHERE sign=-1;


-- 更新溯源ID相关的utxo索引
-- 增量添加utxo_genesis
INSERT INTO utxo_genesis
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM txout_new
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a));
-- 已花费txo标记清除
INSERT INTO utxo_genesis
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txin_new;
-- 如果一个satoshi=0的txo被花费(早期有这个现象)，就可能遗留一个sign=-1的数据，需要删除
ALTER TABLE utxo_genesis DELETE WHERE sign=-1;
