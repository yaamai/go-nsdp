package nsdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PortDuplex int

const (
	UnknownDuplex PortDuplex = iota
	HalfDuplex
	FullDuplex
)

func (s PortDuplex) String() string {
	switch s {
	case UnknownDuplex:
		return "U"
	case HalfDuplex:
		return "H"
	case FullDuplex:
		return "F"
	default:
		return "U"
	}
}

type PortLinkStatus struct {
	Port   int
	Speed  int
	Duplex PortDuplex
}

func (t PortLinkStatus) Tag() uint16 {
	return 0x0c00
}
func (t PortLinkStatus) Length() uint16 {
	return uint16(0)
}
func (t PortLinkStatus) Value() []byte {
	return []byte{}
}
func NewPortLinkStatusFromBytes(b []byte) *PortLinkStatus {
	t := &PortLinkStatus{}
	t.Port = int(b[0])
	switch int(b[1]) {
	case 0:
		t.Duplex = UnknownDuplex
		t.Speed = 0
	case 1:
		t.Duplex = HalfDuplex
		t.Speed = 10
	case 2:
		t.Duplex = FullDuplex
		t.Speed = 10
	case 3:
		t.Duplex = HalfDuplex
		t.Speed = 100
	case 4:
		t.Duplex = FullDuplex
		t.Speed = 100
	case 5:
		t.Duplex = FullDuplex
		t.Speed = 1000
	}
	return t
}
func (t PortLinkStatus) String() string {
	return fmt.Sprintf("%d:%dMbit/s(%v)", t.Port, t.Speed, t.Duplex)
}

type PortStatistics struct {
	Port      int
	Recv      uint64
	Send      uint64
	Pkt       uint64
	Broadcast uint64
	Multicast uint64
	Error     uint64
}

func (t PortStatistics) Tag() uint16 {
	return 0x1000
}
func (t PortStatistics) Length() uint16 {
	return uint16(0)
}
func (t PortStatistics) Value() []byte {
	return []byte{}
}

func NewPortStatisticsFromBytes(b []byte) *PortStatistics {
	t := &PortStatistics{}
	t.Port = int(b[0])
	reader := bytes.NewReader(b[1:])

	binary.Read(reader, binary.BigEndian, &t.Recv)
	binary.Read(reader, binary.BigEndian, &t.Send)
	binary.Read(reader, binary.BigEndian, &t.Pkt)
	binary.Read(reader, binary.BigEndian, &t.Broadcast)
	binary.Read(reader, binary.BigEndian, &t.Multicast)
	binary.Read(reader, binary.BigEndian, &t.Error)
	return t
}
func (t PortStatistics) String() string {
	return fmt.Sprintf("%d:Send=%d, Recv=%d, Pkts=%d, Broadcast=%d, Multicast=%d, Err=%d", t.Port, t.Recv, t.Send, t.Pkt, t.Broadcast, t.Multicast, t.Error)
}
