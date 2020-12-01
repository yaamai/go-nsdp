package main

import (
	"gs308e/nsdp"
	"log"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	resp := c.Read(&nsdp.TLVBase{Tag: uint16(0x7400), Length: 0})
	log.Printf("%v", resp)
}
