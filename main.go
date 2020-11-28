package main

import (
	"log"
	"net"
)

func getSelfIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}

func main() {
	selfAddrStr, err := getSelfIp()
	if err != nil {
		log.Println(err)
	}
	selfAddr, err := net.ResolveUDPAddr("udp", selfAddrStr+":63321")
	if err != nil {
		log.Println(err)
	}
	sendToAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63322")
	if err != nil {
		log.Println(err)
	}
	sendConn, err := net.DialUDP("udp", selfAddr, sendToAddr)
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
		writeLen, err := sendConn.Write(queryModel)
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
