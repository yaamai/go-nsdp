package nsdp

type TLV interface {
	Tag() uint16
	Length() uint16
	Value() []byte
}

func NewTLVFromBytes(tag uint16, length uint16, value []byte) TLV {
	switch tag {
	case 0x0001:
		return &ModelName{StringValue(value)}
	case 0x0003:
		return &HostName{StringValue(value)}
	case 0x0017:
		return &AuthV2PasswordSalt{BytesValue(value)}
	case 0x001a:
		return &AuthV2Password{BytesValue(value)}
	case 0x0c00:
		return NewPortLinkStatusFromBytes(value)
	case 0x1000:
		return NewPortStatisticsFromBytes(value)
	case 0x2400:
		return NewPortVlanMembersFromBytes(value)
	case 0x2800:
		return NewTagVlanMembersFromBytes(value)
	}
	return nil
}
