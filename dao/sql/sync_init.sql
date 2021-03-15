
DROP TABLE blk_height_new;
DROP TABLE blktx_height_new;
DROP TABLE txout_new;
DROP TABLE txin_full_new;

CREATE TABLE IF NOT EXISTS blk_height_new AS blk_height;
CREATE TABLE IF NOT EXISTS blktx_height_new AS blktx_height;
CREATE TABLE IF NOT EXISTS txout_new AS txout;
CREATE TABLE IF NOT EXISTS txin_full_new AS txin_full;
