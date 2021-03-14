-- mempool数据

-- ================================================================
-- mempool包含的交易列表
CREATE TABLE tx (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	txsize       UInt32,
	locktime     UInt32,
	txidx        UInt64
)

-- GetMempoolTxsByRange
SELECT %s FROM tx WHERE idx >= %d AND idx < %d
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
-SELECT %s FROM txin WHERE address = %s

-- GetUTXOByGenesis
 SELECT %s FROM txout WHERE genesis = %s
-SELECT %s FROM txin WHERE genesis = %s
