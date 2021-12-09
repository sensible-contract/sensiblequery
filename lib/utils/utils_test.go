package utils

import (
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	addr := "1dayGM2EfpK6VeyzUZ9dYjjsiDEvdP5Ab"

	pkh, _ := DecodeAddress(addr)
	t.Logf("addr: %s", hex.EncodeToString(pkh))

	pkhHex := "a123a6fdc265e1bbcf1123458891bd7af1a1b5d9"
	pkh, _ = hex.DecodeString(pkhHex)

	t.Logf("addr: %s", EncodeAddress(pkh, PubKeyHashAddrIDMainNet))
}
