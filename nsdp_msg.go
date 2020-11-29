package main

import (
	"bytes"
	"net"
)

var (
	EmptyMac          = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	NSDPDefaultHeader = NSDPHeader{Version: 1, Op: 1, HostMac: EmptyMac, DeviceMac: EmptyMac, Seq: 1, Signature: [4]byte{0x4E, 0x53, 0x44, 0x50}}
	NSDPDefaultBody   = NSDPBody{}
	NSDPDefaultMarker = NSDPMarker{EndOfData: [4]byte{0xff, 0xff, 0x00, 0x00}}
	NSDPDefaultMsg    = NSDPMsg{NSDPDefaultHeader, NSDPDefaultBody, NSDPDefaultMarker}
)

type NSDPHeader struct {
	Version   int8
	Op        int8
	Result    int16
	Unknown1  [4]byte
	HostMac   net.HardwareAddr
	DeviceMac net.HardwareAddr
	Unknown2  [2]byte
	Seq       int16
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

type NSDPBody struct {
	Body []NSDPTLV
}

func (b NSDPBody) WriteToBuffer(buf *bytes.Buffer) {
	for idx, _ := range b.Body {
		b.Body[idx].WriteToBuffer(buf)
	}
}

type NSDPMarker struct {
	EndOfData [4]byte
}

func (m NSDPMarker) WriteToBuffer(b *bytes.Buffer) {
	b.Write(m.EndOfData[:])
}

type NSDPMsg struct {
	NSDPHeader
	NSDPBody
	NSDPMarker
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
