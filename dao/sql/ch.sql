
-- 不需要version版sign mergeTree
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


-- simple mergeTree
DROP TABLE utxo_simple;
CREATE TABLE IF NOT EXISTS utxo_simple (
	txid         FixedString(32),
	vout         UInt32,
	script_type  String
) engine=MergeTree()
ORDER BY (txid, vout)
SETTINGS storage_policy = 'prefer_ssd_policy';
-- 使用join初始化
-- 46172215
INSERT INTO utxo_simple
SELECT txout.txid, txout.vout, txout.script_type FROM txout
ANTI LEFT JOIN txin_spent
ON txout.txid = txin_spent.utxid AND
   txout.vout = txin_spent.vout
WHERE txout.satoshi > 0;
-- 可能需要删除op_return
ALTER TABLE utxo_simple DELETE
WHERE startsWith(script_type, char(0x00, 0x6a)) OR
      startsWith(script_type, char(0x6a));

ALTER TABLE utxo DELETE
WHERE startsWith(script_type, char(0x00, 0x6a)) OR
      startsWith(script_type, char(0x6a));


RENAME TABLE blktx_height_new TO blktx_height

RENAME TABLE utxo_genesis_ver TO utxo_genesis
RENAME TABLE utxo_address_ver TO utxo_address

INSERT INTO utxo_full
  SELECT txid, vout, address, genesis, satoshi, script_type, script, height, 1 FROM txout
SEMI LEFT JOIN utxo
ON txout.txid = utxo.txid AND
   txout.vout = utxo.vout



SELECT txid, idx, address, genesis, satoshi, script_type, height, io_type FROM
(
SELECT txid, vout AS idx, address, genesis, satoshi, script_type, height, 1 AS io_type FROM txout
WHERE (txid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, height, 0 AS io_type FROM txin_full
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
)
)
ORDER BY height DESC
LIMIT 32


    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('4b02ec0a55cf9008871b32a9f30e07dc46c470be')
    ORDER BY height DESC
    LIMIT 64

    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')

SELECT hex(reverse(txid)), idx, satoshi, height, io_type FROM

    SELECT hex(reverse(txid)), idx, height FROM txin_address_height
    WHERE address = unhex('4b02ec0a55cf9008871b32a9f30e07dc46c470be')



SELECT txid, idx, address, genesis, satoshi, script_type, height, 0 AS io_type FROM txin_full

SELECT hex(reverse(txid)), idx, satoshi, height FROM txin_full
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('4b02ec0a55cf9008871b32a9f30e07dc46c470be')
)
)



SELECT hex(reverse(txid)), idx, satoshi, height, io_type FROM
(
SELECT txid, vout AS idx, address, genesis, satoshi, script_type, height, 1 AS io_type FROM txout
WHERE (txid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('4b02ec0a55cf9008871b32a9f30e07dc46c470be')
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, height, 0 AS io_type FROM txin_full
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('4b02ec0a55cf9008871b32a9f30e07dc46c470be')
)
)
ORDER BY height DESC
LIMIT 32





SELECT hex(reverse(txid)), vout, satoshi, height FROM txout_address
SEMI LEFT JOIN
(
    SELECT txid, vout FROM utxo
    WHERE (txid, vout) IN (
          SELECT txid, vout FROM txout_address
          WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    )
) AS t
USING (txid, vout)
WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
LIMIT 32

SELECT hex(reverse(txid)), vout, satoshi, height FROM txout

SELECT hex(reverse(txid)), vout, satoshi, height FROM txout_address
SEMI LEFT JOIN
(
    SELECT txid, vout FROM utxo_address
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    ORDER BY height DESC
    LIMIT 32
) AS t
USING (txid, vout)
WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db') AND
height in (
    SELECT height FROM utxo_address
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    ORDER BY height DESC
    LIMIT 32
)



SELECT hex(reverse(utxid)), vout, satoshi, hex(script_type), height, hex(reverse(u.txid)), u.height FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = reverse(unhex('18470f1da114fa9a790866120cbebcfaa6fd91b2e27941ae9108841c615ef644')) AND
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = reverse(unhex('18470f1da114fa9a790866120cbebcfaa6fd91b2e27941ae9108841c615ef644'))
                    )
) AS u ON txout.utxid = u.utxid AND txout.vout = u.vout
WHERE txid = reverse(unhex('18470f1da114fa9a790866120cbebcfaa6fd91b2e27941ae9108841c615ef644')) AND
    height = 56511



SELECT hex(reverse(utxid)), vout, satoshi, hex(script_type), height FROM txout
WHERE utxid = reverse(unhex('18470f1da114fa9a790866120cbebcfaa6fd91b2e27941ae9108841c615ef644')) AND
     height = 56511




    SELECT hex(utxid), vout, height FROM txout_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    ORDER BY height DESC
    LIMIT 32



SELECT hex(reverse(txid)), vout, satoshi, height FROM txout
WHERE (txid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    ORDER BY height DESC
)
ORDER BY height DESC
LIMIT 10



 AS t
ON txout.txid = t.utxid AND
   txout.vout = t.vout AND
   txout.height = t.height

WHERE height in (
    SELECT height FROM txout_address_height
    WHERE address = unhex('64629ceff0c8da44ef193c6c34a0b5b9aa5f19db')
    ORDER BY height DESC
    LIMIT 32
)



INSERT INTO ta_full
   SELECT ta.k, ta.v,
          tb.k, tb.v FROM ta
   LEFT JOIN tb
   USING k



INSERT INTO ta_full
   SELECT ta.part, ta.k, ta.v,
          tb.part, tb.k, tb.v FROM ta
   LEFT JOIN (
       SELECT part, k, v FROM tb
       WHERE tb.part >= 0 AND tb.part < 10
   ) AS tb
   USING k
   WHERE ta.part >= 0 AND ta.part < 10





CREATE TABLE IF NOT EXISTS txin_test (
	height       UInt32,         --txo 花费的区块高度
	txid         FixedString(32),
	idx          UInt32,

	height_txo   UInt32,         --txo 产生的区块高度
	utxid        FixedString(32),
	vout         UInt32,
	address      FixedString(20),
	genesis      FixedString(20),
	satoshi      UInt64,
	script_type  String
) engine=MergeTree()
ORDER BY (txid, idx)
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'prefer_ssd_policy';



INSERT INTO txin_full
    SELECT txin.height, txin.txid, txin.idx, txin.script_sig, txin.nsequence,
       txout.height, txin.utxid, txin.vout, txout.address, txout.genesis, txout.satoshi, txout.script_type, txout.script_pk FROM txin
    LEFT JOIN txout
    USING (utxid, vout)

SETTINGS join_algorithm='partial_merge'


INSERT INTO txin_full
  SELECT height, txid, idx, script_sig, nsequence,
         txo.height, txo.utxid, txo.vout, txo.address, txo.genesis, txo.satoshi, txo.script_type, txo.script_pk FROM txin
  LEFT JOIN txout AS txo
  USING (utxid, vout)


INSERT INTO txin_full
  SELECT height, txid, idx,
         txo.height, txo.utxid, txo.vout, txo.address, txo.genesis, txo.satoshi, txo.script_type FROM txin
  LEFT JOIN txout AS txo
  USING (utxid, vout)




CREATE TABLE tin_test_v2 engine=MergeTree()
ORDER BY (txid, idx)
PARTITION BY intDiv(height, 2100)
 AS
  SELECT height, txid, idx,
         txo.height, txo.utxid, txo.vout, txo.address, txo.genesis, txo.satoshi, txo.script_type FROM txin
  LEFT JOIN txout AS txo
  USING (utxid, vout)
SETTINGS join_algorithm='partial_merge';

SETTINGS storage_policy = 'prefer_ssd_policy';

--  insert
INSERT INTO txin_full
SELECT txin.height, txin.txid, txin.idx, txin.script_sig, txin.nsequence,
       txout.height, txin.utxid, txin.vout, txout.address, txout.genesis, txout.satoshi, txout.script_type, txout.script_pk FROM txin
LEFT JOIN (
    SELECT height, utxid, vout, address, genesis, satoshi, script_type, script_pk FROM txout
    WHERE (height, utxid, vout) in (

        SELECT tx_height.height, tx_height.txid, txin.vout FROM tx_height
        JOIN (
            SELECT utxid, vout FROM txin WHERE height < 450000 and height >= 430000
        ) AS txin
        ON tx_height.txid = txin.utxid
        WHERE txid in (
            SELECT utxid FROM txin WHERE height < 450000 and height >= 430000
        )
    )
) AS txout
ON txout.utxid = txin.utxid AND
    txout.vout = txin.vout
WHERE height < 450000 and height >= 430000






INSERT INTO txin_full
SELECT txin.height, txin.txid, txin.idx, txin.script_sig, txin.nsequence,
       txout.height, txin.utxid, txin.vout, txout.address, txout.genesis, txout.satoshi, txout.script_type, txout.script_pk FROM txin
LEFT JOIN (
    SELECT height, utxid, vout, address, genesis, satoshi, script_type, script_pk FROM txout
    WHERE (height, utxid, vout) in (

        SELECT tx_height.height, tx_height.txid, txin.vout FROM tx_height
        JOIN (
            SELECT utxid, vout FROM txin WHERE height < 200000
        ) AS txin
        ON tx_height.txid = txin.utxid
        WHERE txid in (
            SELECT utxid FROM txin WHERE height < 200000
        )
    )
) AS txout
ON txout.utxid = txin.utxid AND
    txout.vout = txin.vout
WHERE height < 200000




INSERT INTO txin_full
SELECT txin.height, txin.txid, txin.idx, txin.script_sig, txin.nsequence,
       txout.height, txin.utxid, txin.vout, txout.address, txout.genesis, txout.satoshi, txout.script_type, txout.script_pk FROM txin
LEFT JOIN txout
ON txout.utxid = txin.utxid AND
    txout.vout = txin.vout
SETTINGS join_algorithm='partial_merge'


SELECT %s FROM blk_height ORDER BY height DESC LIMIT 1



CREATE TABLE IF NOT EXISTS tx (
	txid         FixedString(32),
	nin          UInt32,
	nout         UInt32,
	height       UInt32,
	blkid        FixedString(32),
	idx          UInt64
) engine=MergeTree()
ORDER BY height
PARTITION BY intDiv(height, 2100)
SETTINGS storage_policy = 'prefer_nvme_policy';



ALTER TABLE txin_address_height MODIFY COLUMN genesis String
ALTER TABLE txin_address_height MODIFY COLUMN address String

ALTER TABLE txin_genesis_height MODIFY COLUMN genesis String
ALTER TABLE txin_genesis_height MODIFY COLUMN address String


ALTER TABLE txout_address_height MODIFY COLUMN genesis String
ALTER TABLE txout_address_height MODIFY COLUMN address String

ALTER TABLE txout_genesis_height MODIFY COLUMN genesis String
ALTER TABLE txout_genesis_height MODIFY COLUMN address String


ALTER TABLE txin_address_height UPDATE genesis = unhex('00') WHERE genesis = ''


select count(1) from txin_address_height where genesis = ''
select count(1) from txin_address_height where genesis = unhex('0000000000000000000000000000000000000000')

ALTER TABLE txin_address_height UPDATE genesis = unhex('00') WHERE genesis = ''
