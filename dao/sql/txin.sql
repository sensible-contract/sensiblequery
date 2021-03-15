-- txout在哪个高度被花费，按txid首字节分区，分区内按交易txid+idx排序、索引。按txid+idx查询时可确定分区 (快)
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从txid+idx确定花费所在区块height。配合其他表查询
DROP TABLE txout_spent_height
CREATE TABLE IF NOT EXISTS txout_spent_height (
	height       UInt32,
	utxid        FixedString(32),
	vout         UInt32
) engine=MergeTree()
ORDER BY (utxid, vout)
PARTITION BY substring(utxid, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 从txin_spent表创建txout_spent_height
-- INSERT INTO txout_spent_height SELECT height, utxid, vout FROM txin_spent;
-- 需验证排序是否可以加快导入速度
-- INSERT INTO txout_spent_height SELECT height, utxid, vout FROM txin_spent ORDER BY substring(utxid, 1, 1);


-- address在哪些高度的tx中出现，按address首字节分区，分区内按address+genesis+height排序，按address索引。按address查询时可确定分区 (快)
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从address确定所在区块height。配合txin_full源表查询
DROP TABLE txin_address_height;
CREATE TABLE IF NOT EXISTS txin_address_height (
	height       UInt32,
	txid         FixedString(32),
	idx          UInt32,
	address      String,
	genesis      String
) engine=MergeTree()
PRIMARY KEY address
ORDER BY (address, genesis, height)
PARTITION BY substring(address, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 添加
-- INSERT INTO txin_address_height SELECT height, txid, idx, address, genesis FROM txin_full


-- genesis在哪些高度的tx中出现，按genesis首字节分区，分区内按genesis+address+height排序，按genesis索引。按genesis查询时可确定分区 (快)
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从genesis确定所在区块height。配合txin_full源表查询
DROP TABLE txin_genesis_height;
CREATE TABLE IF NOT EXISTS txin_genesis_height (
	height       UInt32,
	txid         FixedString(32),
	idx          UInt32,
	address      String,
	genesis      String
) engine=MergeTree()
PRIMARY KEY genesis
ORDER BY (genesis, address, height)
PARTITION BY substring(genesis, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 添加
-- INSERT INTO txin_genesis_height SELECT height, txid, idx, address, genesis FROM txin_full
