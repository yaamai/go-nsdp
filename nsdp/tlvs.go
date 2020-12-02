package nsdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

/*
from https://github.com/AlbanBedel/libnsdp/blob/master/nsdp_properties.h
from https://github.com/tabacha/ProSafeLinux/blob/master/psl_class.py
 NSDP_PROPERTY_MAC				0x0004
 NSDP_PROPERTY_IP				0x0006
 NSDP_PROPERTY_NETMASK			0x0007
 NSDP_PROPERTY_GATEWAY			0x0008
 NSDP_PROPERTY_PASSWORD			0x000A
 NSDP_PROPERTY_DHCP				0x000B
 NSDP_PROPERTY_FIRMWARE_VERSION	0x000D
 NSDP_PROPERTY_PORT_STATUS		0x0C00
 NSDP_PROPERTY_PORT_STATISTICS	0x1000
 NSDP_PROPERTY_VLAN_ENGINE		0x2000
 NSDP_PROPERTY_VLAN_MEMBERS		0x2800
 NSDP_PROPERTY_PORT_PVID		0x3000
 NSDP_PROPERTY_PORT_COUNT		0x6000

 CMD_MAC = psl_typ.PslTypMac(0x0004, "MAC")
 CMD_LOCATION = psl_typ.PslTypString(0x0005, "location")
 CMD_IP = psl_typ.PslTypIpv4(0x0006, "ip")
 CMD_NETMASK = psl_typ.PslTypIpv4(0x0007, "netmask")
 CMD_GATEWAY = psl_typ.PslTypIpv4(0x0008, "gateway")
 CMD_NEW_PASSWORD = psl_typ.PslTypPassword(0x0009, "new_password", True)
 CMD_PASSWORD = psl_typ.PslTypPassword(0x000a, "password", False)
 CMD_DHCP = psl_typ.PslTypDHCP(0x000b, "dhcp")
 CMD_FIXMEC = psl_typ.PslTypHex(0x000c, "fixmeC")
 CMD_FIRMWAREV = psl_typ.PslTypStringQueryOnly(0x000d, "firmwarever")
 CMD_FIRMWARE2V = psl_typ.PslTypStringQueryOnly(0x000e, "firmware2ver")
 CMD_FIRMWAREACTIVE = psl_typ.PslTypHex(0x000f, "firmware_active")
 CMD_REBOOT = psl_typ.PslTypAction(0x0013, "reboot")
 CMD_FACTORY_RESET = psl_typ.PslTypAction(0x0400, "factory_reset")
 CMD_SPEED_STAT = psl_typ.PslTypSpeedStat(0x0c00, "speed_stat")
 CMD_PORT_STAT = psl_typ.PslTypPortStat(0x1000, "port_stat")
 CMD_RESET_PORT_STAT = psl_typ.PslTypAction(0x1400, "reset_port_stat")
 CMD_TEST_CABLE = psl_typ.PslTypHexNoQuery(0x1800, "test_cable")
 CMD_TEST_CABLE_RESP = psl_typ.PslTypHexNoQuery(0x1c00, "test_cable_resp")
 CMD_VLAN_SUPPORT = psl_typ.PslTypVlanSupport(0x2000, "vlan_support")
 CMD_VLAN_ID = psl_typ.PslTypVlanId(0x2400, "vlan_id")
 CMD_VLAN802_ID = psl_typ.PslTypVlan802Id(0x2800, "vlan802_id")
 CMD_VLANPVID = psl_typ.PslTypVlanPVID(0x3000, "vlan_pvid")
 CMD_QUALITY_OF_SERVICE = psl_typ.PslTypQos(0x3400, "qos")
 CMD_PORT_BASED_QOS = psl_typ.PslTypPortBasedQOS(0x3800, "port_based_qos")
 CMD_BANDWIDTH_INCOMING_LIMIT = psl_typ.PslTypBandwidth(
                                           0x4c00, "bandwidth_in")
 CMD_BANDWIDTH_OUTGOING_LIMIT = psl_typ.PslTypBandwidth(
                                           0x5000, "bandwidth_out")
 CMD_FIXME5400 = psl_typ.PslTypHex(0x5400, "fixme5400")
 CMD_BROADCAST_BANDWIDTH = psl_typ.PslTypBandwidth(0x5800,
              "broadcast_bandwidth")
 CMD_PORT_MIRROR = psl_typ.PslTypPortMirror(0x5c00, "port_mirror")
 CMD_NUMBER_OF_PORTS = psl_typ.PslTypHex(0x6000, "number_of_ports")
 CMD_IGMP_SNOOPING = psl_typ.PslTypIGMPSnooping(0x6800, "igmp_snooping")
 CMD_BLOCK_UNKNOWN_MULTICAST = psl_typ.PslTypBoolean(
                                           0x6c00, "block_unknown_multicast")
 CMD_IGMP_HEADER_VALIDATION = psl_typ.PslTypBoolean(0x7000,
     "igmp_header_validation")
 CMD_FIXME7400 = psl_typ.PslTypHex(0x7400, "fixme7400")
*/
func ParseTLVs(tag uint16, length uint16, value []byte) TLV {
	log.Println("Value:", value)
	switch tag {
	case 0x0001:
		return &ModelName{StringValue{string(value)}}
	case 0x0003:
		return &HostName{StringValue{string(value)}}
	case 0x0c00:
		v := (&PortStatus{}).FromBytes(value)
		return v
	case 0x1000:
		return (&PortStat{}).FromBytes(value)
	case 0x2400:
		return (&PortVlanMembers{}).FromBytes(value)
	case 0x2800:
		return (&TagVlanMembers{}).FromBytes(value)
	}
	return nil
}

type StringValue struct {
	string
}

func (t StringValue) Length() uint16 {
	return uint16(len(t.string))
}
func (t StringValue) Value() []byte {
	return []byte(t.string)
}
func (t StringValue) String() string {
	return t.string
}

type ModelName struct {
	StringValue
}

func (t ModelName) Tag() uint16 {
	return 0x0001
}

type HostName struct {
	StringValue
}

func (t HostName) Tag() uint16 {
	return 0x0003
}

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

type PortStatus struct {
	Port   int
	Speed  int
	Duplex PortDuplex
}

func (t PortStatus) Tag() uint16 {
	return 0x0c00
}
func (t PortStatus) Length() uint16 {
	return uint16(0)
}
func (t PortStatus) Value() []byte {
	return []byte{}
}
func (t *PortStatus) FromBytes(b []byte) *PortStatus {
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
func (t PortStatus) String() string {
	return fmt.Sprintf("%d:%dMbit/s(%v)", t.Port, t.Speed, t.Duplex)
}

type PortStat struct {
	Port      int
	Recv      uint64
	Send      uint64
	Pkt       uint64
	Broadcast uint64
	Multicast uint64
	Error     uint64
}

func (t PortStat) Tag() uint16 {
	return 0x1000
}
func (t PortStat) Length() uint16 {
	return uint16(0)
}
func (t PortStat) Value() []byte {
	return []byte{}
}

func (t *PortStat) FromBytes(b []byte) *PortStat {
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
func (t PortStat) String() string {
	return fmt.Sprintf("%d:Send=%d, Recv=%d, Pkts=%d, Broadcast=%d, Multicast=%d, Err=%d", t.Port, t.Recv, t.Send, t.Pkt, t.Broadcast, t.Multicast, t.Error)
}

func parsePortsBits(buf []byte) []int {
	ports := []int{}
	portNum := 1
	for _, b := range buf {
		for idx := 0x80; idx > 0; idx = idx >> 1 {
			if b&byte(idx) != 0 {
				ports = append(ports, portNum)
			}
			portNum += 1
		}
	}

	return ports
}

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
	return uint16(0)
}
func (t TagVlanMembers) Value() []byte {
	return []byte{}
}
func (t *TagVlanMembers) FromBytes(buf []byte) *TagVlanMembers {
	t.VlanID = int(buf[0])<<8 + int(buf[1])
	t.TaggedPorts = parsePortsBits(buf[2:3])
	t.UnTaggedPorts = parsePortsBits(buf[3:4])
	return t
}
func (t TagVlanMembers) String() string {
	return fmt.Sprintf("%d:%v,%v", t.VlanID, t.TaggedPorts, t.UnTaggedPorts)
}
