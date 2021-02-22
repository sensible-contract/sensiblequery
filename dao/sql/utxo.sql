-- sign mergeTree
DROP TABLE utxo;
CREATE TABLE IF NOT EXISTS utxo (
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script_pk    String,
	height       UInt32,
   sign         Int8
) engine=CollapsingMergeTree(sign)
ORDER BY (utxid, vout)
SETTINGS storage_policy = 'prefer_nvme_policy';
-- sign 初始化
-- 46111030
-- 添加
INSERT INTO utxo
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM txout
  WHERE satoshi > 0;
-- 删除
INSERT INTO utxo
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txin_spent
-- 多余
ALTER TABLE utxo DELETE WHERE sign=-1;
-- 可能需要删除op_return
ALTER TABLE utxo DELETE
WHERE startsWith(script_type, char(0x00, 0x6a)) OR
      startsWith(script_type, char(0x6a));

-- 使用ANTI JOIN添加
INSERT INTO utxo
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM txout
ANTI LEFT JOIN txin_spent
ON txout.utxid = txin_spent.utxid AND
    txout.vout = txin_spent.vout
WHERE txout.satoshi > 0;



-- utxo备份
DROP TABLE utxo_674936;
CREATE TABLE IF NOT EXISTS utxo_674936 (
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script_pk    String,
	height       UInt32,
   sign         Int8
) engine=CollapsingMergeTree(sign)
ORDER BY (utxid, vout)
SETTINGS storage_policy = 'prefer_ssd_policy';
-- 创建备份
INSERT INTO utxo_674936 SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height FROM utxo



-- utxo address
DROP TABLE utxo_address;
CREATE TABLE IF NOT EXISTS utxo_address (
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script_pk    String,
	height       UInt32,
   sign         Int8
) engine=CollapsingMergeTree(sign)
PRIMARY KEY address
ORDER BY (address, genesis, height)
SETTINGS storage_policy = 'prefer_nvme_policy';
-- 从utxo创建address
INSERT INTO utxo_address SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM utxo
-- 通过address查询utxo
SELECT hex(reverse(utxid)), vout, satoshi, height FROM utxo_address
WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
ORDER BY height DESC
LIMIT 32



-- utxo genesis
DROP TABLE utxo_genesis;
CREATE TABLE IF NOT EXISTS utxo_genesis (
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String,
	script_pk    String,
	height       UInt32,
   sign         Int8
) engine=CollapsingMergeTree(sign)
PRIMARY KEY genesis
ORDER BY (genesis, address, height)
SETTINGS storage_policy = 'prefer_nvme_policy';
-- 从utxo创建address
INSERT INTO utxo_genesis SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM utxo
