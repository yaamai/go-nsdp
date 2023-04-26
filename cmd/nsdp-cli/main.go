package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yaamai/go-nsdp/nsdp"
	"log"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: nsdp-cli [flags] query/set attr-name:param1:param2:param3, ...\n")
	fmt.Fprintf(os.Stderr, "attrs:\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] model-name\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] host-name\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] macaddr\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] ipaddr\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] netmask\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] gateway\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] port-link-stat\n")
	fmt.Fprintf(os.Stderr, "  [Q  ] port-stat\n")
	fmt.Fprintf(os.Stderr, "  [Q/S] tag-vlan-pvid\n")
	fmt.Fprintf(os.Stderr, "        param1: port-id\n")
	fmt.Fprintf(os.Stderr, "        param2: vlan-id\n")
	fmt.Fprintf(os.Stderr, "  [Q/S] tag-vlan\n")
	fmt.Fprintf(os.Stderr, "        param1: vlan-id\n")
	fmt.Fprintf(os.Stderr, "        param2: untagged port-id\n")
	fmt.Fprintf(os.Stderr, "        param3: tagged port-id\n")
	fmt.Fprintf(os.Stderr, "examples: nsdp-cli set tag-vlan:1000:1,6:1\n")
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
		case "macaddr":
			result = append(result, nsdp.MacAddress{})
		case "ipaddr":
			result = append(result, nsdp.HostIPAddress{})
		case "netmask":
			result = append(result, nsdp.Netmask{})
		case "gateway":
			result = append(result, nsdp.GatewayAddress{})
		case "port-link-stat":
			result = append(result, nsdp.PortLinkStatus{})
		case "port-stat":
			result = append(result, nsdp.PortStatistics{})
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
					tlv.UnTaggedPorts = SplitInts(params[2], ",")
					tlv.TaggedPorts = SplitInts(params[3], ",")
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
		netdevice = flag.String("device", "", "network interface to use (default first non loopback)")
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

	c, err := nsdp.NewDefaultClient(*netdevice)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		resp *nsdp.Msg
	)
	if action == "set" {
		if *password == "" {
			*password = os.Getenv("NSDP_PASSWORD")
		}
		resp, err = c.WriteWithAuth(*password, tlvs...)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		resp, err = c.Read(tlvs...)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// to pretty output TLVs, make map
	tlvMap := map[string]interface{}{}
	for _, tlv := range resp.Body {
		tagName := tlv.Tag().String()
		if _, ok := tlvMap[tagName]; ok {
			// if same Tag present, aggregate to array
			array, ok := tlvMap[tagName].([]nsdp.TLV)
			if !ok {
				prevValue := tlvMap[tagName]
				tlvMap[tagName] = []nsdp.TLV{prevValue.(nsdp.TLV), tlv}
			} else {
				tlvMap[tagName] = append(array, tlv)
			}
		} else {
			tlvMap[tagName] = tlv
		}
	}
	jsonBytes, err := json.MarshalIndent(tlvMap, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(jsonBytes))
}
