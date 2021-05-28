package script

import (
	"bytes"
	"encoding/binary"
	"satosensible/lib/blkparser"
)

var empty = make([]byte, 1)

func ExtractPkScriptGenesisIdAndAddressPkh(pkscript []byte) (isNFT bool, codeHash, genesisId, addressPkh, metaTxId []byte, name, symbol string, value, decimal uint64) {
	scriptLen := len(pkscript)
	if scriptLen < 1024 {
		return false, empty, empty, empty, empty, "", "", 0, 0
	}

	metaTxId = empty

	dataLen := 0
	genesisIdLen := 0
	genesisOffset := scriptLen - 8 - 4
	valueOffset := scriptLen - 8 - 4 - 8
	addressOffset := scriptLen - 8 - 4 - 8 - 20
	decimalOffset := scriptLen - 8 - 4 - 8 - 20 - 1
	symbolOffset := scriptLen - 8 - 4 - 8 - 20 - 1 - 1
	nameOffset := scriptLen - 8 - 4 - 8 - 20 - 1 - 1 - 10

	if (bytes.HasSuffix(pkscript, []byte("sensible")) || bytes.HasSuffix(pkscript, []byte("oraclesv"))) &&
		pkscript[scriptLen-8-4] == 1 { // PROTO_TYPE == 1

		if pkscript[scriptLen-72-36-1-1] == 0x4c && pkscript[scriptLen-72-36-1] == 108 {
			genesisIdLen = 36        // new ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-72-20-1-1] == 0x4c && pkscript[scriptLen-72-20-1] == 92 {
			genesisIdLen = 20        // old ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-50-36-1-1] == 0x4c && pkscript[scriptLen-50-36-1] == 86 {
			genesisIdLen = 36        // old ft
			dataLen = 1 + 1 + 1 + 50 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-92-20-1-1] == 0x4c && pkscript[scriptLen-92-20-1] == 112 {
			genesisIdLen = 20        // old ft
			dataLen = 1 + 1 + 1 + 92 // opreturn + 0x4c + pushdata + data
		} else {
			genesisIdLen = 40        // error ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		}

		genesisOffset -= genesisIdLen
		valueOffset -= genesisIdLen
		addressOffset -= genesisIdLen
		decimalOffset -= genesisIdLen

		decimal = uint64(pkscript[decimalOffset])
		name = string(pkscript[nameOffset : nameOffset+20])
		symbol = string(pkscript[symbolOffset : symbolOffset+10])
	} else if pkscript[scriptLen-1] < 2 && pkscript[scriptLen-37-1] == 37 && pkscript[scriptLen-37-1-40-1] == 40 && pkscript[scriptLen-37-1-40-1-1] == OP_RETURN {
		// nft issue
		isNFT = true
		genesisIdLen = 40
		genesisOffset = scriptLen - 37 - 1 - genesisIdLen
		valueOffset = scriptLen - 1 - 8
		addressOffset = scriptLen - 1 - 8 - 8 - 20

		dataLen = 1 + 1 + 1 + 37 // opreturn + pushdata + pushdata + data
	} else if pkscript[scriptLen-1] == 1 && pkscript[scriptLen-61-1] == 61 && pkscript[scriptLen-61-1-40-1] == 40 && pkscript[scriptLen-61-1-40-1-1] == OP_RETURN {
		// nft transfer
		isNFT = true
		genesisIdLen = 40
		genesisOffset = scriptLen - 61 - 1 - genesisIdLen
		metaTxIdOffset := scriptLen - 1 - 32
		valueOffset = scriptLen - 1 - 32 - 8
		addressOffset = scriptLen - 1 - 32 - 8 - 20

		dataLen = 1 + 1 + 1 + 61 // opreturn + pushdata + pushdata + data

		metaTxId = make([]byte, 32)
		copy(metaTxId, pkscript[metaTxIdOffset:metaTxIdOffset+32])
	} else {
		return false, empty, empty, empty, empty, "", "", 0, 0
	}

	genesisId = make([]byte, genesisIdLen)
	addressPkh = make([]byte, 20)
	copy(genesisId, pkscript[genesisOffset:genesisOffset+genesisIdLen])
	copy(addressPkh, pkscript[addressOffset:addressOffset+20])

	value = binary.LittleEndian.Uint64(pkscript[valueOffset : valueOffset+8])

	codeHash = blkparser.GetHash160(pkscript[:scriptLen-genesisIdLen-dataLen])

	return isNFT, codeHash, genesisId, addressPkh, metaTxId, name, symbol, value, decimal
}
