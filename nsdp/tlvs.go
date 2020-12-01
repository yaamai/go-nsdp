package nsdp

import (
	"bytes"
)

/*
from https://github.com/AlbanBedel/libnsdp/blob/master/nsdp_properties.h
from https://github.com/tabacha/ProSafeLinux/blob/master/psl_class.py
 NSDP_PROPERTY_NONE				0x0000
 NSDP_PROPERTY_MODEL			0x0001
 NSDP_PROPERTY_HOSTNAME			0x0003
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

 CMD_MODEL = psl_typ.PslTypStringQueryOnly(0x0001, "model")
 CMD_FIMXE2 = psl_typ.PslTypHex(0x0002, "fixme2")
 CMD_NAME = psl_typ.PslTypString(0x0003, "name")
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
func ParseTLVs(base TLVBase) TLV {
	switch base.Tag {
	case 1:
		v := TLVModelName{}
		v.ReadFromBase(base)
		return &v
	default:
		return &base
	}
}

type TLVModelName struct {
	Name string
}

func (t TLVModelName) WriteToBuffer(b *bytes.Buffer) {
	TLVBase{Tag: 1, Length: uint16(len([]byte(t.Name))), Value: []byte(t.Name)}.WriteToBuffer(b)
}

func (t *TLVModelName) ReadFromBase(base TLVBase) {
	t.Name = string(base.Value)
}

func (t *TLVModelName) String() string {
	return t.Name
}

type TLVName struct {
	Name string
}

func (t TLVName) WriteToBuffer(b *bytes.Buffer) {
	TLVBase{Tag: 0x17, Length: uint16(len([]byte(t.Name))), Value: []byte(t.Name)}.WriteToBuffer(b)
}

func (t *TLVName) ReadFromBase(base TLVBase) {
	t.Name = string(base.Value)
}

func (t *TLVName) String() string {
	return t.Name
}
