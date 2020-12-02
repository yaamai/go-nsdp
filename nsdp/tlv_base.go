package nsdp

type TLV interface {
	Tag() uint16
	Length() uint16
	Value() []byte
}
