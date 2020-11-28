package main

import (
	"log"
	"net"
)

// get first non-loopback address, intf-name and mac
func getSelfIntfAndIp() (string, string, []byte, error) {
	intfs, err := net.Interfaces()
	if err != nil {
		return "", "", nil, err
	}

	for _, intf := range intfs {
		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), intf.Name, intf.HardwareAddr, nil
				}
			}
		}
	}

	return "", "", nil, nil
}

func main() {
	selfAddrStr, intfName, intfHwAddr, err := getSelfIntfAndIp()
	if err != nil {
		log.Println(err)
	}
	log.Println("using", intfName, selfAddrStr)

	// NSDP require two socket
	// 192.168.0.xxx:63321 -> 255.255.255.255:63322
	// 192.168.0.yyy:63322 -> 192.168.0.xxx:63321

	selfAddr, err := net.ResolveUDPAddr("udp", selfAddrStr+":63321")
	if err != nil {
		log.Println(err)
	}
	sendToAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63322")
	if err != nil {
		log.Println(err)
	}
	anyAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:63322")
	if err != nil {
		log.Println(err)
	}
	sendConn, err := net.DialUDP("udp", selfAddr, anyAddr)
	if err != nil {
		log.Println(err)
	}

	queryModel := []byte{
		0x01,
		0x01,
		0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
		0x00, 0x01,
		0x4E, 0x53, 0x44, 0x50,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00,
		0xff, 0xff, 0x00, 0x00,
		//		0x00, 0x00, 0xff, 0xff,
	}
	// set mac addr
	for idx, b := range intfHwAddr {
		queryModel[8+idx] = b
	}
	//
	// 01
	// 02
	// 0700
	// 0000 0000
	// 0000 0000 00
	// 0044 a56e 49
	// c7b2
	// 0000
	// 0001
	// 4e53 4450
	// 0000 0000
	// ffff 0000
	for range []int{1, 2, 3} {
		writeLen, err := sendConn.WriteTo(queryModel, sendToAddr)
		log.Println(writeLen, err)
	}

	go func() {
		for {
			buf := make([]byte, 65535)
			readLen, err := sendConn.Read(buf)
			log.Println(readLen, err)

		}
	}()

	for {
	}
}
