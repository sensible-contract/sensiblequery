-- 基础区块数据按每2100个区块分区保存到clickhouse，此数据保证和最长链一致

-- block
-- ================================================================
-- 区块头，分区内按区块blkid排序、索引。按blkid查询时将遍历所有分区 (慢)
DROP TABLE blk;
CREATE TABLE IF NOT EXISTS blk (
	height       UInt32,
	blkid        FixedString(32),
	previd       FixedString(32),
	ntx          UInt64
) engine=MergeTree()
ORDER BY blkid
PARTITION BY intDiv(height, 2100);

-- 区块头，分区内按区块高度height排序、索引。按blk height查询时可确定分区 (快)
DROP TABLE blk_height;
CREATE TABLE IF NOT EXISTS blk_height (
	height       UInt32,
	blkid        FixedString(32),
	previd       FixedString(32),
	ntx          UInt64
) engine=MergeTree()
ORDER BY height
PARTITION BY intDiv(height, 2100);

-- tx list
-- ================================================================
-- 区块包含的交易列表，分区内按交易txid排序、索引。仅按txid查询时将遍历所有分区 (慢)
-- 查询需附带height。可配合tx_height表查询
DROP TABLE tx;
CREATE TABLE IF NOT EXISTS tx (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	height       UInt32,
	blkid        FixedString(32),
	idx          UInt64
) engine=MergeTree()
ORDER BY txid
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';

-- 区块包含的交易列表，分区内按区块高度height排序、索引。按blk height查询时可确定分区 (快)
DROP TABLE blktx_height;
CREATE TABLE IF NOT EXISTS blktx_height (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	height       UInt32,
	blkid        FixedString(32),
	idx          UInt64
) engine=MergeTree()
ORDER BY height
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';

-- txin
-- ================================================================
-- 交易输入列表，分区内按交易txid+idx排序、索引，单条记录包括输入的各种细节。仅按txid查询时将遍历所有分区（慢）
-- 查询需附带height。可配合tx_height表查询
DROP TABLE txin_full;
CREATE TABLE IF NOT EXISTS txin_full (
	height       UInt32,
	txid         FixedString(32),
	idx          UInt32,
	script_sig   String,
	height_txo   UInt32,
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String
) engine=MergeTree()
ORDER BY (txid, idx)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';

-- 交易输入的outpoint列表，分区内按outpoint txid+idx排序、索引。用于查询某txo被哪个tx花费，需遍历所有分区（慢）
-- 查询需附带height，需配合txout_spent_height表查询
DROP TABLE txin_spent;
CREATE TABLE IF NOT EXISTS txin_spent (
	height       UInt32,
	txid         FixedString(32),
	idx          UInt32,
	utxid        FixedString(32),
	vout         UInt32
) engine=MergeTree()
ORDER BY (utxid, vout)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';

-- txout
-- ================================================================
-- 交易输出列表，分区内按交易txid+idx排序、索引，单条记录包括输出的各种细节。仅按txid查询时将遍历所有分区（慢）
-- 查询需附带height，可配合tx_height表查询
DROP TABLE txout;
CREATE TABLE IF NOT EXISTS txout (
	txid         FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script       String,
	height       UInt32
) engine=MergeTree()
ORDER BY (txid, vout)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';
