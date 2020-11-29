package nsdp

import (
	"bytes"
	"fmt"
	"net"
)

type Header struct {
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

func (h Header) WriteToBuffer(b *bytes.Buffer) {
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

func (h *Header) ReadFromBuffer(b *bytes.Reader) {
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

func (h Header) String() string {
	return fmt.Sprintf("V: %d, Op: %d(%d), HostMAC: %v, DevMAC: %v, Seq: %d", h.Version, h.Op, h.Result, h.HostMac, h.DeviceMac, h.Seq)
}
