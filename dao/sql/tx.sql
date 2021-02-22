-- tx在哪个高度被打包，按txid首字节分区，分区内按交易txid排序、索引。按txid查询时可确定分区（快）
-- 此数据表不能保证和最长链一致，而是包括所有已打包tx的height信息，其中可能存在已被孤立的块高度
-- 主要用于从txid确定所在区块height。配合其他表查询
CREATE TABLE IF NOT EXISTS tx_height (
	txid         FixedString(32),
	height       UInt32
) engine=MergeTree()
ORDER BY txid
PARTITION BY substring(txid, 1, 1)
SETTINGS storage_policy = 'prefer_nvme_policy';

-- 从tx表创建tx_height
-- INSERT INTO tx_height SELECT txid, height FROM tx;

-- 查询例子
-- SELECT hex(txid), height from tx_height where txid = unhex('24915ea87ae4f2ebbb30179b3ae35e5538c066550b717953d778e6d272088965');

-- 查询例子
-- SELECT height, hex(txid), idx,
--        height_txo, hex(utxid), vout, satoshi FROM txin_full
-- WHERE txid = reverse(unhex('24915ea87ae4f2ebbb30179b3ae35e5538c066550b717953d778e6d272088965')) AND
-- height IN (
--     SELECT height FROM tx_height
--     WHERE txid = reverse(unhex('24915ea87ae4f2ebbb30179b3ae35e5538c066550b717953d778e6d272088965'))
-- )
