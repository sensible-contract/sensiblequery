
-- block
DROP TABLE blk;
CREATE TABLE IF NOT EXISTS blk (
	height       UInt32,
	blkid        FixedString(32),
	previd       FixedString(32),
	ntx          UInt64
) engine=MergeTree()
ORDER BY blkid
PARTITION BY intDiv(height, 2100);


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

DROP TABLE blktx;
CREATE TABLE IF NOT EXISTS blktx (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	height       UInt32,
	blkid        FixedString(32),
	idx          UInt64
) engine=MergeTree()
ORDER BY blkid
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';


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



-- txin
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
	-- script_pk    String,
	script_type  String
) engine=MergeTree()
ORDER BY (txid, idx)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';


DROP TABLE txin_spent;
CREATE TABLE IF NOT EXISTS txin_spent (
	height       UInt32
	txid         FixedString(32),
	idx          UInt32,
	utxid        FixedString(32),
	vout         UInt32,
) engine=MergeTree()
ORDER BY (utxid, vout)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';


-- txout
DROP TABLE utxo;
CREATE TABLE IF NOT EXISTS utxo (
	height       UInt32,
	txid         FixedString(32),
	vout         UInt32,
   sign         Int8,
   version      UInt8
) engine=VersionedCollapsingMergeTree(sign, version)
ORDER BY (txid, vout)
PARTITION BY intDiv(height, 2100);

INSERT INTO utxo VALUES (1, 'a', 1, 1, 1)
INSERT INTO utxo VALUES (1, 'a', 1, 2, 1)


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


DROP TABLE txout_address;
CREATE TABLE IF NOT EXISTS txout_address (
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
ORDER BY (address, genesis)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';


DROP TABLE txout_genesis;
CREATE TABLE IF NOT EXISTS txout_genesis (
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
ORDER BY (genesis, address)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'multi_tiered_policy';




cat /data/blk.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO blk FORMAT RowBinary"
cat /data/blk.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO blk_height FORMAT RowBinary"

cat /data/tx.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO tx FORMAT RowBinary"
cat /data/tx.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO blktx FORMAT RowBinary"

cat /data/tx-in.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO txin FORMAT RowBinary"
cat /data/tx-in.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO txin_spent FORMAT RowBinary"

cat /data/tx-out.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO txout FORMAT RowBinary"
cat /data/tx-out.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO txout_address FORMAT RowBinary"
cat /data/tx-out.ch | clickhouse-client -h 192.168.31.236 --database="bsv" --query="INSERT INTO txout_genesis FORMAT RowBinary"
