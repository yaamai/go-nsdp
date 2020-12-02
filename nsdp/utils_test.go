package nsdp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombinePortsBits(t *testing.T) {

	empty := combinePortsBits([]int{})
	assert.Equal(t, []byte{0}, empty)

	withInOneByte := combinePortsBits([]int{2, 3, 8})
	// 0110 0001 = 0x61
	assert.Equal(t, []byte{0x61}, withInOneByte)

	twoBytes := combinePortsBits([]int{2, 3, 8, 9, 13})
	// 0110 0001 1000 1000 = 0x6188
	assert.Equal(t, []byte{0x61, 0x88}, twoBytes)
}
