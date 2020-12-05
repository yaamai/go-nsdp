package nsdp

import (
	"encoding/json"
	"net"
)

type StringValue string

func (t StringValue) Length() uint16 {
	return uint16(len(t))
}

func (t StringValue) Value() []byte {
	return []byte(t)
}

func (t StringValue) String() string {
	return string(t)
}

func (t StringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

type BytesValue []byte

func (t BytesValue) Length() uint16 {
	return uint16(len(t))
}

func (t BytesValue) Value() []byte {
	return t
}

func (t BytesValue) String() string {
	return string(t)
}

func (t BytesValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

type IPV4Address struct {
	net.IP
}

func (t IPV4Address) Length() uint16 {
	return 4
}

func (t IPV4Address) Value() []byte {
	return t.IP
}

func (t IPV4Address) String() string {
	return t.IP.String()
}

func (t IPV4Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
