package blkparser

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/ripemd160"
)

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

func HashString(data []byte) (res string) {
	n := len(data) // 32
	reverseData := make([]byte, n)

	// need reverse
	for i := 0; i < n; i++ {
		reverseData[i] = data[n-i-1]
	}
	return hex.EncodeToString(reverseData)
}
