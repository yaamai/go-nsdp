package main

import (
	"bytes"
	"fmt"
	"net"
)

var (
	EmptyMac          = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	NSDPDefaultHeader = NSDPHeader{Version: 1, Op: 1, HostMac: EmptyMac, DeviceMac: EmptyMac, Seq: 1, Signature: [4]byte{0x4E, 0x53, 0x44, 0x50}}
	NSDPDefaultBody   = NSDPBody{}
	NSDPDefaultMarker = NSDPMarker{EndOfData: [4]byte{0xff, 0xff, 0x00, 0x00}}
	NSDPDefaultMsg    = NSDPMsg{NSDPDefaultHeader, NSDPDefaultBody, NSDPDefaultMarker}
)

func readInt8(b *bytes.Reader) int8 {
	v, _ := b.ReadByte()
	return int8(v)
}
func readInt16(b *bytes.Reader) int16 {
	v1, _ := b.ReadByte()
	v2, _ := b.ReadByte()

	var v int16
	v = int16(v1) << 8
	v = v | int16(v2)
	return v
}

type NSDPHeader struct {
	Version   int8
	Op        int8
	Result    int16
	Unknown1  [4]byte
	HostMac   net.HardwareAddr
	DeviceMac net.HardwareAddr
	Unknown2  [2]byte
	Seq       uint16
	Signature [4]byte
	Unknown3  [4]byte
}

func (h NSDPHeader) WriteToBuffer(b *bytes.Buffer) {
	b.WriteByte(byte(h.Version))
	b.WriteByte(byte(h.Op))
	b.WriteByte(byte(h.Result >> 8))
	b.WriteByte(byte(h.Result & 0xff))
	b.Write(h.Unknown1[:])
	b.Write(h.HostMac)
	b.Write(h.DeviceMac)
	b.Write(h.Unknown2[:])
	b.WriteByte(byte(h.Seq >> 8))
	b.WriteByte(byte(h.Seq & 0xff))
	b.Write(h.Signature[:])
	b.Write(h.Unknown3[:])
}

func (h *NSDPHeader) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 32 {
		return
	}
	h.Version = readInt8(b)
	h.Op = readInt8(b)
	h.Result = readInt16(b)
	b.Read(h.Unknown1[:])
	h.HostMac = make([]byte, 6)
	b.Read(h.HostMac[:])
	h.DeviceMac = make([]byte, 6)
	b.Read(h.DeviceMac[:])
	b.Read(h.Unknown2[:])
	h.Seq = uint16(readInt16(b))
	b.Read(h.Signature[:])
	b.Read(h.Unknown3[:])
}

func (h NSDPHeader) String() string {
	return fmt.Sprintf("V: %d, Op: %d(%d), HostMAC: %v, DevMAC: %v, Seq: %d", h.Version, h.Op, h.Result, h.HostMac, h.DeviceMac, h.Seq)
}

type NSDPTLV struct {
	Tag    int16
	Length int16
	Value  []byte
}

func (t NSDPTLV) WriteToBuffer(b *bytes.Buffer) {
	b.WriteByte(byte(t.Tag >> 8))
	b.WriteByte(byte(t.Tag & 0xff))
	b.WriteByte(byte(t.Length >> 8))
	b.WriteByte(byte(t.Length & 0xff))
	b.Write(t.Value)
}
func (t *NSDPTLV) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 4 {
		return
	}
	t.Tag = readInt16(b)
	t.Length = readInt16(b)

	if b.Len() < int(t.Length) {
		return
	}
	t.Value = make([]byte, t.Length)
	b.Read(t.Value)
}

type NSDPBody struct {
	Body []NSDPTLV
}

func (b NSDPBody) WriteToBuffer(buf *bytes.Buffer) {
	for idx, _ := range b.Body {
		b.Body[idx].WriteToBuffer(buf)
	}
}
func (b *NSDPBody) ReadFromBuffer(buf *bytes.Reader) {
	for buf.Len() > 4 {
		tlv := NSDPTLV{}
		tlv.ReadFromBuffer(buf)
		b.Body = append(b.Body, tlv)
	}
}

type NSDPMarker struct {
	EndOfData [4]byte
}

func (m NSDPMarker) WriteToBuffer(b *bytes.Buffer) {
	b.Write(m.EndOfData[:])
}

func (m *NSDPMarker) ReadFromBuffer(b *bytes.Reader) {
	if b.Len() < 4 {
		return
	}
	b.Read(m.EndOfData[:])
}

func (b NSDPMarker) String() string {
	if bytes.Compare(b.EndOfData[:], NSDPDefaultMarker.EndOfData[:]) == 0 {
		return "<MARK>"
	} else {
		return "<INVALID-MARK>"
	}
}

type NSDPMsg struct {
	NSDPHeader
	NSDPBody
	NSDPMarker
}

func ParseNSDPMsg(buf []byte) *NSDPMsg {
	m := &NSDPMsg{}
	r := bytes.NewReader(buf)
	m.NSDPHeader.ReadFromBuffer(r)
	m.NSDPBody.ReadFromBuffer(r)
	m.NSDPMarker.ReadFromBuffer(r)
	return m
}

func (m NSDPMsg) WriteToBuffer(b *bytes.Buffer) {
	m.NSDPHeader.WriteToBuffer(b)
	m.NSDPBody.WriteToBuffer(b)
	m.NSDPMarker.WriteToBuffer(b)
}

func (m NSDPMsg) Bytes() []byte {
	b := bytes.Buffer{}
	m.WriteToBuffer(&b)
	return b.Bytes()
}
