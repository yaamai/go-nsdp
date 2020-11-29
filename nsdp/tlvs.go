package nsdp

import (
	"bytes"
)

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
