package nsdp

type NewPassword struct {
	StringValue
}

func (t NewPassword) Tag() uint16 {
	return 0x0009
}

type Password struct {
	StringValue
}

func (t Password) Tag() uint16 {
	return 0x000a
}

type AuthV2PasswordSalt struct {
	BytesValue
}

func (t AuthV2PasswordSalt) Tag() uint16 {
	return 0x0017
}

type AuthV2Password struct {
	BytesValue
}

func (t AuthV2Password) Tag() uint16 {
	return 0x001a
}
