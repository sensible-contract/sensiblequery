package utils

import (
	"errors"
	"satoblock/lib/base58"

	"golang.org/x/crypto/ripemd160"
)

func ReverseBytes(data []byte) (result []byte) {
	for _, b := range data {
		result = append([]byte{b}, result...)
	}
	return result
}

var (
	PubKeyHashAddrIDMainNet = byte(0x00) // starts with 1
	PubKeyHashAddrIDTestNet = byte(0x6f) // starts with m or n
	ErrChecksumMismatch     = errors.New("checksum mismatch")
)

// encodeAddress returns a human-readable payment address given a ripemd160 hash
// and netID which encodes the bitcoin network and address type.  It is used
// in both pay-to-pubkey-hash (P2PKH) and pay-to-script-hash (P2SH) address
// encoding.
func EncodeAddress(hash160 []byte, netID byte) string {
	// Format is 1 byte for a network and address class (i.e. P2PKH vs
	// P2SH), 20 bytes for a RIPEMD160 hash, and 4 bytes of checksum.
	return base58.CheckEncode(hash160[:ripemd160.Size], netID)
}

func DecodeAddress(addr string) (decoded []byte, err error) {
	// Switch on decoded length to determine the type.
	decoded, _, err = base58.CheckDecode(addr)
	if err != nil {
		if err == base58.ErrChecksum {
			return nil, ErrChecksumMismatch
		}
		return nil, errors.New("decoded address is of unknown format")
	}
	return
}
