package script

import utils "satoblock/lib/blkparser"

func ExtractPkScriptForTxo(Pkscript, scriptType []byte) (isNFT bool, codeHash, genesisId, addressPkh []byte, value, decimal uint64) {
	if isPubkeyHash(scriptType) {
		addressPkh = make([]byte, 20)
		copy(addressPkh, Pkscript[3:23])
		return false, empty, empty, addressPkh, 0, 0
	}

	if isPayToScriptHash(scriptType) {
		addressPkh = utils.GetHash160(Pkscript[2 : len(Pkscript)-1])
		return false, empty, empty, addressPkh, 0, 0
	}

	if isPubkey(scriptType) {
		addressPkh = utils.GetHash160(Pkscript[1 : len(Pkscript)-1])
		return false, empty, empty, addressPkh, 0, 0
	}

	// if isMultiSig(scriptType) {
	// 	return Pkscript[:]
	// }

	return ExtractPkScriptGenesisIdAndAddressPkh(Pkscript)
}

func GetLockingScriptType(pkscript []byte) (scriptType []byte) {
	length := len(pkscript)
	if length == 0 {
		return
	}
	scriptType = make([]byte, 0)

	lenType := 0
	p := uint(0)
	e := uint(length)

	for p < e && lenType < 32 {
		c := pkscript[p]
		if 0 < c && c < 0x4f {
			cnt, cntsize := utils.DecodeVarIntForScript(pkscript[p:])
			p += cnt + cntsize
		} else {
			p += 1
		}
		scriptType = append(scriptType, c)
		lenType += 1
	}
	return
}

func IsLockingScriptOnlyEqual(pkscript []byte) bool {
	// test locking script
	// "0b 3c4b616e7965323032303e 87"

	length := len(pkscript)
	if length == 0 {
		return true
	}
	if pkscript[length-1] != 0x87 {
		return false
	}
	cnt, cntsize := utils.DecodeVarIntForScript(pkscript)
	if length == int(cnt+cntsize+1) {
		return true
	}
	return false
}
