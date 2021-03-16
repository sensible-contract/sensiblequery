-- ================ 如果没有孤块，则无需处理
ALTER TABLE blktx_height DELETE WHERE height > ;
ALTER TABLE blk_height DELETE WHERE height > ;
ALTER TABLE blk DELETE WHERE height > ;

ALTER TABLE blktx_height DELETE WHERE height > ;
ALTER TABLE tx DELETE WHERE height > ;

ALTER TABLE txin_spent DELETE WHERE height > ;

ALTER TABLE txin_full DELETE WHERE height > ;
ALTER TABLE txout DELETE WHERE height > ;
