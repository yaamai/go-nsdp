package nsdp

type AuthV2PasswordSalt struct {
	BytesValue
}

func (t AuthV2PasswordSalt) Tag() Tag {
	return TagAuthV2PasswordSalt
}

type AuthV2Password struct {
	BytesValue
}

func (t AuthV2Password) Tag() Tag {
	return TagAuthV2Password
}
