package main

import (
	"flag"
	"fmt"
	"github.com/yaamai/go-nsdp/nsdp"
	"log"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: nsdp-cli [flags] action action-args...\n")
	fmt.Fprintf(os.Stderr, "examples: nsdp-cli query modelname hostname\n")
	fmt.Fprintf(os.Stderr, "          nsdp-cli set tag-vlan:1000:6,7,8:1\n")
	fmt.Fprintf(os.Stderr, "\nflags:\n")
	flag.PrintDefaults()
}

func SplitInts(s, sep string) []int {
	ret := []int{}
	ary := strings.Split(s, sep)
	for _, elm := range ary {
		i, err := strconv.Atoi(elm)
		if err != nil {
			continue
		}
		ret = append(ret, i)
	}

	return ret
}

func ConvertCmdsToTLVs(cmds []string) []nsdp.TLV {
	var result []nsdp.TLV
	for _, cmd := range cmds {
		params := strings.Split(cmd, ":")
		cmdType := params[0]

		switch cmdType {
		case "model-name":
			result = append(result, nsdp.ModelName{})
		case "host-name":
			result = append(result, nsdp.HostName{})
		case "tag-vlan-pvid":
			tlv := nsdp.TagVlanPVID{}
			if len(params) > 2 {
				portId, err1 := strconv.Atoi(params[1])
				vlanId, err2 := strconv.Atoi(params[2])
				if err1 == nil && err2 == nil {
					tlv.PortID = portId
					tlv.VlanID = vlanId
				}
			}
			result = append(result, tlv)
		case "tag-vlan":
			tlv := nsdp.TagVlanMembers{}
			if len(params) > 3 {
				vlanId, err := strconv.Atoi(params[1])
				if err == nil {
					tlv.VlanID = vlanId
					tlv.TaggedPorts = SplitInts(params[2], ",")
					tlv.UnTaggedPorts = SplitInts(params[3], ",")
				}
			}
			result = append(result, tlv)
		}
	}

	return result
}

func main() {
	var (
		password = flag.String("password", "", "switch password (required when setting values)")
	)
	flag.Usage = usage
	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) < 2 {
		usage()
		os.Exit(1)
	}
	action := positionalArgs[0]
	tlvs := ConvertCmdsToTLVs(positionalArgs[1:])

	c, err := nsdp.NewDefaultClient()
	if err != nil {
		log.Fatalln(err)
	}

	if action == "set" {
		resp, err := c.WriteWithAuth(*password, tlvs...)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%v", resp)
	} else {
		resp, err := c.Read(tlvs...)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%v", resp)
	}
}
