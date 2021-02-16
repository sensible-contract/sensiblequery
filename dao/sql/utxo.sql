-- sign mergeTree
DROP TABLE utxo_full;
CREATE TABLE IF NOT EXISTS utxo_full (
	txid         FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script       String,
	height       UInt32,
   sign         Int8
) engine=CollapsingMergeTree(sign)
ORDER BY (txid, vout)
SETTINGS storage_policy = 'prefer_ssd_policy';
-- sign 初始化
-- 46111030
-- 添加
INSERT INTO utxo_full
  SELECT txid, vout, address, genesis, satoshi, script_type, height, 1 FROM txout
  WHERE satoshi > 0;
-- 删除
INSERT INTO utxo_full
  SELECT utxid, vout,'', '', 0, '', 0, -1 FROM txin_spent
-- 多余
ALTER TABLE utxo_full DELETE WHERE sign=-1;
-- 可能需要删除op_return
ALTER TABLE utxo_full DELETE
WHERE startsWith(script_type, char(0x00, 0x6a)) OR
      startsWith(script_type, char(0x6a));



-- utxo备份
DROP TABLE utxo_673490;
CREATE TABLE IF NOT EXISTS utxo_673490 (
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
SETTINGS storage_policy = 'prefer_ssd_policy';
-- 创建备份
INSERT INTO utxo_673490 SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM utxo_full



-- utxo address
DROP TABLE utxo_address;
CREATE TABLE IF NOT EXISTS utxo_address (
	txid         FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script       String,
	height       UInt32
) engine=MergeTree()
PRIMARY KEY address
ORDER BY (address, genesis, height)
SETTINGS storage_policy = 'prefer_ssd_policy';
-- 从utxo_full创建address
INSERT INTO utxo_address SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM utxo_full
-- 通过address查询utxo
SELECT hex(reverse(txid)), vout, satoshi, height FROM utxo_address
WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
ORDER BY height DESC
LIMIT 32



-- utxo genesis
DROP TABLE utxo_genesis;
CREATE TABLE IF NOT EXISTS utxo_genesis (
	txid         FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script       String,
	height       UInt32
) engine=MergeTree()
PRIMARY KEY genesis
ORDER BY (genesis, address, height)
SETTINGS storage_policy = 'prefer_ssd_policy';
-- 从utxo_full创建address
INSERT INTO utxo_genesis SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM utxo_full
