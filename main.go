package main

import (
	"flag"
	"gs308e/nsdp"
	"log"
)

func main() {
	var (
		password = flag.String("password", "", "switch password (required when setting values)")
	)
	flag.Parse()

	c, err := NewDefaultClient()
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Login(*password)
	if err != nil {
		log.Fatalln(err)
	}

	vlan := nsdp.TagVlanMembers{VlanID: 1000, TaggedPorts: []int{7, 8}, UnTaggedPorts: []int{}}
	resp, err := c.WriteWithAuth(vlan)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%v", resp)
}
