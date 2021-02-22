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
