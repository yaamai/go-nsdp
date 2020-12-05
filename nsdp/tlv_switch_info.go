package nsdp

import (
	"encoding/json"
	"net"
)

type ModelName struct {
	StringValue
}

func (t ModelName) Tag() Tag {
	return TagModelName
}

type HostName struct {
	StringValue
}

func (t HostName) Tag() Tag {
	return TagHostName
}

type MacAddress struct {
	net.HardwareAddr
}

func (t MacAddress) Tag() Tag {
	return TagMAC
}

func (t MacAddress) Length() uint16 {
	return uint16(len(t.HardwareAddr))
}

func (t MacAddress) Value() []byte {
	return t.HardwareAddr
}

func (t MacAddress) String() string {
	return t.HardwareAddr.String()
}

func (t MacAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type HostIPAddress struct {
	IPV4Address
}

func (t HostIPAddress) Tag() Tag {
	return TagIP
}

type Netmask struct {
	IPV4Address
}

func (t Netmask) Tag() Tag {
	return TagNetmask
}

type GatewayAddress struct {
	IPV4Address
}

func (t GatewayAddress) Tag() Tag {
	return TagGateway
}
