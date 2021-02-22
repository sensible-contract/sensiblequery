-- 生成区块id索引
INSERT INTO blk SELECT * FROM blk_height;


-- 生成区块内tx索引
INSERT INTO tx SELECT * FROM blktx_height;
-- 生成tx到区块高度索引
INSERT INTO tx_height SELECT txid, height FROM tx;


-- 生成txo被花费的tx索引
INSERT INTO txin_spent SELECT height, txid, idx, utxid, vout FROM txin;
-- 生成txo被花费的tx区块高度索引
INSERT INTO txout_spent_height SELECT height, utxid, vout FROM txin_spent;


-- 生成地址参与的输出索引
INSERT INTO txout_address_height SELECT height, utxid, vout, address, genesis FROM txout;
-- 生成溯源ID参与的输出索引
INSERT INTO txout_genesis_height SELECT height, utxid, vout, address, genesis FROM txout;

-- 生成地址参与输入的相关tx区块高度索引
INSERT INTO txin_address_height SELECT height, txid, idx, address, genesis FROM txin_full;
-- 生成溯源ID参与输入的相关tx区块高度索引
INSERT INTO txin_genesis_height SELECT height, txid, idx, address, genesis FROM txin_full;


-- 不执行
-- 生成txin_full的JOIN语句在大数据量时无法执行，需要直接导入txin_full表的数据
-- INSERT INTO txin_full
--   SELECT height, txid, idx, script_sig, nsequence,
--          txo.height, txo.utxid, txo.vout, txo.address, txo.genesis, txo.satoshi, txo.script_type, txo.script_pk FROM txin
--   LEFT JOIN txout AS txo
--   USING (utxid, vout)

-- 不执行
-- 如果txin_full表通过直接导入，则txin表可使用txin_full来生成；或者直接导入txin表的数据
-- INSERT INTO txin SELECT txid, idx, utxid, vout, script_sig, nsequence, height FROM txin_full


-- 全量生成utxo
INSERT INTO utxo
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM txout
  ANTI LEFT JOIN txin_spent
  USING (utxid, vout)
  WHERE txout.satoshi > 0 AND
        NOT startsWith(script_type, char(0x6a)) AND
        NOT startsWith(script_type, char(0x00, 0x6a));

-- 生成地址相关的utxo索引
INSERT INTO utxo_address SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM utxo;
-- 生成溯源ID相关的utxo索引
INSERT INTO utxo_genesis SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height, 1 FROM utxo;
