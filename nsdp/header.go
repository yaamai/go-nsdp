package nsdp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

const (
	HeaderLength = 32
)

type ResultCode int

const (
	ResultSuccess     ResultCode = iota
	ResultInvalidAuth            = 0x0013
)

func (c ResultCode) String() string {
	switch c {
	case ResultSuccess:
		return "Success"
	case ResultInvalidAuth:
		return "InvalidAuth"
	default:
		return "UnknownError"
	}
}

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

func (h Header) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}
	err := h.MarshalBinaryBuffer(&buf)
	return buf.Bytes(), err
}
func (h *Header) UnmarshalBinary(buf []byte) error {
	r := bytes.NewReader(buf)
	return h.UnmarshalBinaryBuffer(r)
}

func (h Header) MarshalBinaryBuffer(buf *bytes.Buffer) error {

	buf.WriteByte(byte(h.Version))
	buf.WriteByte(byte(h.Op))
	buf.WriteByte(byte(h.Result >> 8))
	buf.WriteByte(byte(h.Result & 0xff))
	buf.Write(h.Unknown1[:])
	buf.Write(h.HostMac)
	buf.Write(h.DeviceMac)
	buf.Write(h.Unknown2[:])
	buf.WriteByte(byte(h.Seq >> 8))
	buf.WriteByte(byte(h.Seq & 0xff))
	buf.Write(h.Signature[:])
	buf.Write(h.Unknown3[:])

	return nil
}

func (h *Header) UnmarshalBinaryBuffer(r *bytes.Reader) error {
	if r.Len() < 32 {
		return errors.New("too short header length")
	}
	h.Version = readInt8(r)
	h.Op = readInt8(r)
	h.Result = readInt16(r)
	r.Read(h.Unknown1[:])
	h.HostMac = make([]byte, 6)
	r.Read(h.HostMac[:])
	h.DeviceMac = make([]byte, 6)
	r.Read(h.DeviceMac[:])
	r.Read(h.Unknown2[:])
	h.Seq = readUint16(r)
	r.Read(h.Signature[:])
	r.Read(h.Unknown3[:])

	return nil
}

func (h Header) String() string {
	return fmt.Sprintf("V: %d, Op: %d(%v), HostMAC: %v, DevMAC: %v, Seq: %d", h.Version, h.Op, h.Result, h.HostMac, h.DeviceMac, h.Seq)
}
