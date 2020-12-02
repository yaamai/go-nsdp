package main

import (
	"fmt"
	"gs308e/nsdp"
	"log"
)

func CalcAuthV2Password(password string, mac, salt []byte) []byte {
	buf := make([]byte, 0xff)
	copy(buf[0x38:], []byte(password))

	var (
		al, bl, dl uint8
	)
	dl = mac[5]
	al = mac[1]
	al ^= dl
	dl = uint8(salt[3]&0xff) ^ salt[2]
	buf[0x13] = al

	dl ^= al
	dl ^= buf[0x3A]
	al = mac[0]
	dl ^= buf[0x39]
	dl ^= buf[0x38]
	buf[0x4C] = dl

	bl = salt[3] ^ salt[1]
	dl = mac[4]
	bl ^= dl
	bl ^= al
	bl ^= buf[0x3C]
	al = mac[3]
	bl ^= buf[0x3B]
	al ^= mac[2]
	bl ^= buf[0x3D]
	dl ^= mac[5]
	buf[0x4D] = bl

	bl = salt[0] ^ salt[2]
	buf[0x18] = dl

	bl ^= al
	bl ^= buf[0x40]
	bl ^= buf[0x3E]
	bl ^= buf[0x3F]
	buf[0x4E] = bl

	bl = salt[0] ^ salt[1]
	bl ^= dl
	bl ^= buf[0x43]
	bl ^= buf[0x42]
	bl ^= buf[0x41]
	buf[0x4F] = bl

	dl = salt[3] ^ salt[2]
	dl ^= buf[0x13]
	dl ^= buf[0x44]
	dl ^= buf[0x46]
	dl ^= buf[0x45]
	buf[0x50] = dl

	dl = salt[3] ^ salt[1]
	dl ^= mac[4]
	dl ^= mac[0]
	dl ^= buf[0x49]
	dl ^= buf[0x48]
	dl ^= buf[0x47]
	buf[0x51] = dl

	dl = salt[0] ^ salt[2]
	dl ^= al
	dl ^= buf[0x4B]
	dl ^= buf[0x4A]
	dl ^= buf[0x38]
	buf[0x52] = dl

	al = salt[0] ^ salt[1]
	al ^= buf[0x18]
	al ^= buf[0x3B]
	al ^= buf[0x3D]
	al ^= buf[0x39]
	buf[0x53] = al

	return buf[0x4c:0x54]
}
func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	resp := c.Read(nsdp.AuthV2PasswordSalt{})
	log.Printf("%v", resp)

	password := ""
	mac := []byte(resp.Header.DeviceMac[:6]) // to format clearly
	salt := resp.Body.Body[0].(*nsdp.AuthV2PasswordSalt).BytesValue
	log.Printf("%s, %x, %x", password, mac, salt)
	encodedPassword := CalcAuthV2Password(password, mac, salt)
	fmt.Printf("%x", encodedPassword)

	auth := nsdp.AuthV2Password{BytesValue: nsdp.BytesValue(encodedPassword)}
	vlan := nsdp.TagVlanMembers{VlanID: 1000, TaggedPorts: []int{8}, UnTaggedPorts: []int{}}
	// pass := nsdp.NewPassword{StringValue: nsdp.StringValue("Test1234")}
	resp = c.Write(auth, vlan)
	log.Printf("%v", resp)
}
