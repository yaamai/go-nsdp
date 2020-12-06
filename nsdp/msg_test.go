package nsdp

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestNewTagVlanMsgParse(t *testing.T) {

	input, err := hex.DecodeString("0103000000000000111111111111111111111111000001124e53445000000000001a000800000000000000002800000403e88484ffff0000")
	assert.Nil(t, err)

	expected := &Msg{
		Header: Header{Version: 1, Op: 3, Result: 0, Unknown1: [4]uint8{0x0, 0x0, 0x0, 0x0},
			HostMac: net.HardwareAddr{0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, DeviceMac: net.HardwareAddr{0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
			Unknown2: [2]uint8{0x0, 0x0}, Seq: 274, Signature: [4]uint8{0x4e, 0x53, 0x44, 0x50}, Unknown3: [4]uint8{0x0, 0x0, 0x0, 0x0}},
		Body: Body([]TLV{
			&AuthV2Password{BytesValue([]byte{0, 0, 0, 0, 0, 0, 0, 0})},
			&TagVlanMembers{VlanID: 1000, TaggedPorts: []int{1, 6}, UnTaggedPorts: []int{1, 6}},
		}),
		Marker: Marker{EndOfData: [4]uint8{0xff, 0xff, 0x0, 0x0}},
	}

	msg := &Msg{}
	msg.UnmarshalBinary(input)

	assert.Equal(t, expected, msg)
}

func TestTagVlanMsgParse(t *testing.T) {

	input, err := hex.DecodeString("0103000000000000111111111111111111111111000000e14e53445000000000001a000800000000000000002800000403e8e721ffff0000")
	assert.Nil(t, err)

	expected := &Msg{
		Header: Header{Version: 1, Op: 3, Result: 0, Unknown1: [4]uint8{0x0, 0x0, 0x0, 0x0},
			HostMac: net.HardwareAddr{0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, DeviceMac: net.HardwareAddr{0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
			Unknown2: [2]uint8{0x0, 0x0}, Seq: 225, Signature: [4]uint8{0x4e, 0x53, 0x44, 0x50}, Unknown3: [4]uint8{0x0, 0x0, 0x0, 0x0}},
		Body: Body([]TLV{
			&AuthV2Password{BytesValue([]byte{0, 0, 0, 0, 0, 0, 0, 0})},
			&TagVlanMembers{VlanID: 1000, TaggedPorts: []int{1, 2, 3, 6, 7, 8}, UnTaggedPorts: []int{3, 8}},
		}),
		Marker: Marker{EndOfData: [4]uint8{0xff, 0xff, 0x0, 0x0}},
	}

	msg := &Msg{}
	msg.UnmarshalBinary(input)

	assert.Equal(t, expected, msg)
}

func TestEmptyMsg(t *testing.T) {
	buf, err := DefaultMsg.MarshalBinary()
	assert.Nil(t, err)

	empty := []byte{
		0x01,
		0x01,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x01,
		0x4E, 0x53, 0x44, 0x50,
		0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00,
		//		0x00, 0x00, 0xff, 0xff,
	}

	assert.Equal(t, empty, buf)
}

func TestMsgParse(t *testing.T) {

	input := []byte{
		1,
		2,
		0, 0,
		0, 0, 0, 0,
		17, 17, 17, 17, 17, 17,
		68, 165, 110, 17, 17, 17,
		0, 0,
		211, 43,
		78, 83, 68, 80,
		0, 0, 0, 0,
		0, 1, 0, 6, 71, 83, 51, 48, 56, 69,
		255, 255, 0, 0,
	}
	expected := &Msg{
		Header: Header{Version: 1, Op: 2, Result: 0, Unknown1: [4]uint8{0x0, 0x0, 0x0, 0x0},
			HostMac: net.HardwareAddr{0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, DeviceMac: net.HardwareAddr{0x44, 0xa5, 0x6e, 0x11, 0x11, 0x11},
			Unknown2: [2]uint8{0x0, 0x0}, Seq: 0xd32b, Signature: [4]uint8{0x4e, 0x53, 0x44, 0x50}, Unknown3: [4]uint8{0x0, 0x0, 0x0, 0x0}},
		Body:   Body([]TLV{&ModelName{StringValue("GS308E")}}),
		Marker: Marker{EndOfData: [4]uint8{0xff, 0xff, 0x0, 0x0}},
	}

	msg := &Msg{}
	msg.UnmarshalBinary(input)

	assert.Equal(t, expected, msg)
}
