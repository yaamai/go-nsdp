package nsdp

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcAuthV2Password(t *testing.T) {
	salt, _ := hex.DecodeString("6ab70143")
	mac, _ := hex.DecodeString("44a56e111111")
	password := "testtest"

	out := CalcAuthV2Password(password, mac, salt)
	expected, _ := hex.DecodeString("94c413ddf6a160a9")

	assert.Equal(t, expected, out)
}
