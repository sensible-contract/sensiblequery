-- 基础区块数据，此数据保证和最长链一致

-- ================================================================
-- 区块头
CREATE TABLE blk (
	height       UInt32,
	blkid        FixedString(32),
	previd       FixedString(32),
	merkle       FixedString(32),
	ntx          UInt64,
	blocktime    UInt32,
	bits         UInt32,
	blocksize    UInt32
)

-- GetBlocksByHeightRange
SELECT %s FROM blk WHERE height >= %d AND height < %d
-- GetBlockByHeight
SELECT %s FROM blk WHERE height = %d LIMIT 1
-- GetBlockById
SELECT %s FROM blk WHERE blkid = %s LIMIT 1
-- GetBestBlock
SELECT %s FROM blk ORDER BY height DESC LIMIT 1

-- ================================================================
-- 区块包含的交易列表
CREATE TABLE tx (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	txsize       UInt32,
	locktime     UInt32,
	height       UInt32,
	blkid        FixedString(32),
	txidx        UInt64
)

-- GetBlockTxsByBlockHeight
SELECT %s FROM tx WHERE height = %d
-- GetBlockTxsByBlockId
SELECT %s FROM tx WHERE blkid = %s
-- GetTxById
SELECT %s FROM tx WHERE txid = %s LIMIT 1

-- ================================================================
-- 交易输出列表
CREATE TABLE txout (
	utxid        FixedString(32),
	vout         UInt32,
	address      String,
	genesis      String,
	satoshi      UInt64,
	script_type  String,
	script_pk    String,
	height       UInt32,         --txo 产生的区块高度
	txidx        UInt64
)

-- txout_full
INSERT INTO txout_full
   SELECT * FROM txout
   LEFT JOIN txin
   USING (utxid, vout)

-- GetTxOutputsByTxId
SELECT %s FROM txout_full WHERE utxid = %s
-- GetTxOutputByTxIdAndIdx
SELECT %s FROM txout WHERE utxid = %s AND vout = %d LIMIT 1


-- ================================================================
-- 交易输入列表
CREATE TABLE txin (
	txid         FixedString(32),
	idx          UInt32,
	utxid        FixedString(32),
	vout         UInt32,
	script_sig   String,
	nsequence    UInt32,
	height       UInt32,         --txo 花费的区块高度
	txidx        UInt64
)

-- txin_full
INSERT INTO txin_full
   SELECT * FROM txin
   LEFT JOIN txout
   USING (utxid, vout)

-- GetTxInputsByTxId
SELECT %s FROM txin_full WHERE txid = %s
-- GetTxInputByTxIdAndIdx
SELECT %s FROM txin_full WHERE txid = %s AND idx = %d LIMIT 1
-- GetTxOutputSpentStatusByTxIdAndIdx
SELECT %s FROM txin WHERE utxid = %s AND vout = %d LIMIT 1


-- history/utxo
-- ================================================================
-- GetHistoryByAddress
 SELECT %s FROM txout WHERE address = %s
+SELECT %s FROM txin_full WHERE address = %s
-- GetHistoryByGenesis
 SELECT %s FROM txout WHERE genesis = %s
+SELECT %s FROM txin_full WHERE genesis = %s

-- GetUTXOByAddress
 SELECT %s FROM txout WHERE address = %s
-SELECT %s FROM txin_full WHERE address = %s

-- GetUTXOByGenesis
 SELECT %s FROM txout WHERE genesis = %s
-SELECT %s FROM txin_full WHERE genesis = %s
