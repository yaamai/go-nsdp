package nsdp

import (
	"fmt"
)

type PortVlanMembers struct {
	VlanID int
	Ports  []int
}

func (t PortVlanMembers) Tag() uint16 {
	return 0x2400
}
func (t PortVlanMembers) Length() uint16 {
	return uint16(0)
}
func (t PortVlanMembers) Value() []byte {
	return []byte{}
}
func (t *PortVlanMembers) FromBytes(buf []byte) *PortVlanMembers {
	t.VlanID = int(buf[0])<<8 + int(buf[1])
	t.Ports = parsePortsBits(buf[2:])
	return t
}
func (t PortVlanMembers) String() string {
	return fmt.Sprintf("%d:%v", t.VlanID, t.Ports)
}

type TagVlanMembers struct {
	VlanID        int
	TaggedPorts   []int
	UnTaggedPorts []int
}

func (t TagVlanMembers) Tag() uint16 {
	return 0x2800
}
func (t TagVlanMembers) Length() uint16 {
	return uint16(len(t.Value()))
}
func (t TagVlanMembers) Value() []byte {
	tagged := combinePortsBits(t.TaggedPorts)
	untagged := combinePortsBits(t.UnTaggedPorts)
	vlanid := []byte{byte((t.VlanID >> 8) & 0xff), byte(t.VlanID & 0xff)}

	// maybe len(tagged) and len(untagged) must be equal. but not checked
	return append(vlanid, append(tagged, untagged...)...)
}
func (t *TagVlanMembers) FromBytes(buf []byte) *TagVlanMembers {
	t.VlanID = int(buf[0])<<8 + int(buf[1])
	t.TaggedPorts = parsePortsBits(buf[2:3])
	t.UnTaggedPorts = parsePortsBits(buf[3:4])
	return t
}
func (t TagVlanMembers) String() string {
	return fmt.Sprintf("{%d:%v,%v}", t.VlanID, t.TaggedPorts, t.UnTaggedPorts)
}
