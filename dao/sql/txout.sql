-- address在哪些高度的tx中出现，按address首字节分区，分区内按address+genesis+height排序，按address索引。按address查询时可确定分区 (快)
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从address确定所在区块height。配合txout源表查询
DROP TABLE txout_address_height;
CREATE TABLE IF NOT EXISTS txout_address_height (
	height       UInt32,
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20)
) engine=MergeTree()
PRIMARY KEY address
ORDER BY (address, genesis, height)
PARTITION BY substring(address, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 添加
INSERT INTO txout_address_height
  SELECT height, utxid, vout, address, genesis FROM txout



-- genesis在哪些高度的tx中出现，按genesis首字节分区，分区内按genesis+address+height排序，按genesis索引。按genesis查询时可确定分区 (快)
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从genesis确定所在区块height。配合txout源表查询
DROP TABLE txout_genesis_height;
CREATE TABLE IF NOT EXISTS txout_genesis_height (
	height       UInt32,
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20)
) engine=MergeTree()
PRIMARY KEY genesis
ORDER BY (genesis, address, height)
PARTITION BY substring(genesis, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 添加
INSERT INTO txout_genesis_height
  SELECT height, utxid, vout, address, genesis FROM txout
