-- ================ 如果没有孤块，则无需处理
SET custom_height = 23432;


ALTER TABLE blk_height DELETE WHERE height > getSetting('custom_height');
ALTER TABLE blk DELETE WHERE height > getSetting('custom_height');

ALTER TABLE blktx_height DELETE WHERE height > getSetting('custom_height');
ALTER TABLE tx DELETE WHERE height > getSetting('custom_height');

ALTER TABLE txin DELETE WHERE height > getSetting('custom_height');
ALTER TABLE txin_spent DELETE WHERE height > getSetting('custom_height');

-- 回滚已被花费的utxo_address
INSERT INTO utxo
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height_txo, 1 FROM txin_full
  WHERE satoshi > 0 AND
      height > getSetting('custom_height');
-- 删除新添加的utxo_address˜
INSERT INTO utxo_address
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txout
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a)) AND
      height > getSetting('custom_height');

-- 回滚已被花费的utxo_genesis
INSERT INTO utxo_genesis
  SELECT utxid, vout, address, genesis, satoshi, script_type, script_pk, height_txo, 1 FROM txin_full
  WHERE satoshi > 0 AND
      height > getSetting('custom_height');
-- 删除新添加的utxo_genesis
INSERT INTO utxo_genesis
  SELECT utxid, vout,'', '', 0, '', '', 0, -1 FROM txout
  WHERE satoshi > 0 AND
      NOT startsWith(script_type, char(0x6a)) AND
      NOT startsWith(script_type, char(0x00, 0x6a)) AND
      height > getSetting('custom_height');

ALTER TABLE txin_full DELETE WHERE height > getSetting('custom_height');
ALTER TABLE txout DELETE WHERE height > getSetting('custom_height');
