package script

import "bytes"

var empty = make([]byte, 1)

func ExtractPkScriptGenesisIdAndAddressPkh(pkscript []byte) (genesisId, addressPkh []byte) {
	scriptLen := len(pkscript)
	if scriptLen < 1024 {
		return empty, empty
	}
	if !bytes.HasSuffix(pkscript, []byte("oraclesv")) {
		return empty, empty
	}

	genesisOffset := scriptLen - 8 - 4 - 20
	addressOffset := scriptLen - 8 - 4 - 20 - 8 - 20
	if pkscript[scriptLen-8-4] != 1 {
		return empty, empty
	}

	genesisId = make([]byte, 20)
	addressPkh = make([]byte, 20)
	copy(genesisId, pkscript[genesisOffset:genesisOffset+20])
	copy(addressPkh, pkscript[addressOffset:addressOffset+20])
	return genesisId, addressPkh
}
