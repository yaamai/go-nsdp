package nsdp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagVlanMember(t *testing.T) {
	m := TagVlanMembers{VlanID: 1000, TaggedPorts: []int{1}, UnTaggedPorts: []int{1, 2, 3}}
	assert.Equal(t, []byte{0x03, 0xe8, 0x80, 0xe0}, m.Value())
}
