package blkparser

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/ripemd160"
)

func DecodeVarIntForTx(raw []byte) (cnt uint, cnt_size uint) {
	if len(raw) < 1 {
		return 0, 0
	}
	if raw[0] < 0xfd {
		return uint(raw[0]), 1
	}

	if raw[0] == 0xfd {
		if len(raw) < 3 {
			return 0, 0
		}
		return uint(binary.LittleEndian.Uint16(raw[1:3])), 3

	} else if raw[0] == 0xfe {
		if len(raw) < 5 {
			return 0, 0
		}
		return uint(binary.LittleEndian.Uint32(raw[1:5])), 5
	}

	if len(raw) < 9 {
		return 0, 0
	}
	return uint(binary.LittleEndian.Uint64(raw[1:9])), 9
}

func DecodeVarIntForScript(raw []byte) (cnt uint, cnt_size uint) {
	if len(raw) < 1 {
		return 0, 0
	}
	if raw[0] < 0x4c {
		return uint(raw[0]), 1
	}

	if raw[0] == 0x4c {
		if len(raw) < 2 {
			return 0, 0
		}
		return uint(raw[1]), 2

	} else if raw[0] == 0x4d {
		if len(raw) < 3 {
			return 0, 0
		}
		return uint(binary.LittleEndian.Uint16(raw[1:3])), 3

	} else if raw[0] == 0x4e {
		if len(raw) < 5 {
			return 0, 0
		}
		return uint(binary.LittleEndian.Uint32(raw[1:5])), 5
	}

	return 0, 0
}

func GetHash256(data []byte) (hash []byte) {
	sha := sha256.New()
	sha.Write(data[:])
	tmp := sha.Sum(nil)
	sha.Reset()
	sha.Write(tmp)
	hash = sha.Sum(nil)
	return
}

func GetHash160(data []byte) (hash []byte) {
	sha := sha256.New()
	sha.Write(data[:])
	tmp := sha.Sum(nil)

	rp := ripemd160.New()
	rp.Write(tmp)
	hash = rp.Sum(nil)
	return
}

func HashString(data []byte) (res string) {
	reverseData := make([]byte, 32)

	// need reverse
	for i := 0; i < 32; i++ {
		reverseData[i] = data[32-i-1]
	}
	return hex.EncodeToString(reverseData)
}
