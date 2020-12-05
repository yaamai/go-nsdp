package nsdp

import (
	"net"
)

type Tag uint16

const (
	TagUnknown            Tag = iota
	TagModelName              = 0x0001
	TagHostName               = 0x0003
	TagMAC                    = 0x0004
	TagIP                     = 0x0006
	TagNetmask                = 0x0007
	TagGateway                = 0x0008
	TagAuthV2PasswordSalt     = 0x0017
	TagAuthV2Password         = 0x001a
	TagPortLinkStatus         = 0x0c00
	TagPortStatistics         = 0x1000
	TagPortVlanMembers        = 0x2400
	TagTagVlanMembers         = 0x2800
	TagTagVlanPVID            = 0x3000
)

func (t Tag) String() string {
	switch t {
	case TagModelName:
		return "model_name"
	case TagHostName:
		return "host_name"
	case TagMAC:
		return "mac"
	case TagIP:
		return "ip"
	case TagNetmask:
		return "netmask"
	case TagGateway:
		return "gateway"
	case TagPortLinkStatus:
		return "port_link_status"
	case TagPortStatistics:
		return "port_statistics"
	case TagPortVlanMembers:
		return "port_vlans"
	case TagTagVlanMembers:
		return "tag_vlans"
	case TagTagVlanPVID:
		return "tag_vlan_pvids"
	default:
		return ""
	}
}

type TLV interface {
	Tag() Tag
	Length() uint16
	Value() []byte
}

func NewTLVFromBytes(tag uint16, length uint16, value []byte) TLV {
	switch tag {
	case TagModelName:
		return &ModelName{StringValue(value)}
	case TagHostName:
		return &HostName{StringValue(value)}
	case TagMAC:
		return &MacAddress{net.HardwareAddr(value)}
	case TagIP:
		return &HostIPAddress{IPV4Address{net.IP(value)}}
	case TagNetmask:
		return &Netmask{IPV4Address{net.IP(value)}}
	case TagGateway:
		return &GatewayAddress{IPV4Address{net.IP(value)}}
	case TagAuthV2PasswordSalt:
		return &AuthV2PasswordSalt{BytesValue(value)}
	case TagAuthV2Password:
		return &AuthV2Password{BytesValue(value)}
	case TagPortLinkStatus:
		return NewPortLinkStatusFromBytes(value)
	case TagPortStatistics:
		return NewPortStatisticsFromBytes(value)
	case TagPortVlanMembers:
		return NewPortVlanMembersFromBytes(value)
	case TagTagVlanMembers:
		return NewTagVlanMembersFromBytes(value)
	case TagTagVlanPVID:
		return NewTagVlanPVIDFromBytes(value)
	}
	return nil
}
