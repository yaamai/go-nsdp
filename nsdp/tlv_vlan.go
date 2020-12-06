package nsdp

import (
	"fmt"
)

type PortVlanMembers struct {
	VlanID int   `json:"vlan_id"`
	Ports  []int `json:"ports"`
}

func (t PortVlanMembers) Tag() Tag {
	return TagPortVlanMembers
}
func (t PortVlanMembers) Length() uint16 {
	return uint16(0)
}
func (t PortVlanMembers) Value() []byte {
	return []byte{}
}
func NewPortVlanMembersFromBytes(buf []byte) *PortVlanMembers {
	t := &PortVlanMembers{}
	t.VlanID = int(buf[0])<<8 + int(buf[1])
	t.Ports = parsePortsBits(buf[2:])
	return t
}
func (t PortVlanMembers) String() string {
	return fmt.Sprintf("%d:%v", t.VlanID, t.Ports)
}

type TagVlanMembers struct {
	VlanID        int   `json:"vlan_id"`
	TaggedPorts   []int `json:"tagged_ports"`
	UnTaggedPorts []int `json:"untagged_ports"`
}

func (t TagVlanMembers) Tag() Tag {
	return TagTagVlanMembers
}
func (t TagVlanMembers) Length() uint16 {
	return uint16(len(t.Value()))
}
func (t TagVlanMembers) Value() []byte {
	tagged := combinePortsBits(t.TaggedPorts)
	untagged := combinePortsBits(t.UnTaggedPorts)
	vlanid := []byte{byte((t.VlanID >> 8) & 0xff), byte(t.VlanID & 0xff)}

	// maybe len(tagged) and len(untagged) must be equal. but not checked
	return append(vlanid, append(untagged, tagged...)...)
}
func NewTagVlanMembersFromBytes(buf []byte) *TagVlanMembers {
	t := &TagVlanMembers{}
	t.VlanID = int(buf[0])<<8 + int(buf[1])
	t.TaggedPorts = parsePortsBits(buf[2:3])
	t.UnTaggedPorts = parsePortsBits(buf[3:4])
	return t
}
func (t TagVlanMembers) String() string {
	return fmt.Sprintf("{%d:%v,%v}", t.VlanID, t.TaggedPorts, t.UnTaggedPorts)
}

type TagVlanPVID struct {
	PortID int `json:"port_id"`
	VlanID int `json:"vlan_id"`
}

func (t TagVlanPVID) Tag() Tag {
	return TagTagVlanPVID
}
func (t TagVlanPVID) Length() uint16 {
	return uint16(len(t.Value()))
}
func (t TagVlanPVID) Value() []byte {
	return []byte{
		byte(t.PortID),
		byte((t.VlanID >> 8) & 0xff), byte(t.VlanID & 0xff)}
}
func NewTagVlanPVIDFromBytes(buf []byte) *TagVlanPVID {
	t := &TagVlanPVID{}
	t.PortID = int(buf[0])
	t.VlanID = int(buf[1])<<8 + int(buf[2])
	return t
}
func (t TagVlanPVID) String() string {
	return fmt.Sprintf("{%d:%d}", t.PortID, t.VlanID)
}
