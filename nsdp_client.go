package main

import (
	"errors"
	"gs308e/nsdp"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	DefaultDestAddr          = "255.255.255.255"
	DefaultRecvPort          = "63321"
	DefaultSendPort          = "63322"
	DefaultReceiveBufferSize = 0xffff
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

type Client struct {
	listenAddr   *net.UDPAddr
	targetAddr   *net.UDPAddr
	sourceHwAddr net.HardwareAddr
	password     []byte
	targetHwAddr net.HardwareAddr
	conn         *net.UDPConn
	seq          uint16
}

func NewDefaultClient() (*Client, error) {
	selfAddrStr, _, intfHwAddr, err := getSelfIntfAndIp()
	if err != nil {
		return nil, err
	}

	selfAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(selfAddrStr, DefaultRecvPort))
	if err != nil {
		return nil, err
	}

	anyAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:63322")
	if err != nil {
		return nil, err
	}

	return NewClient(selfAddr, anyAddr, net.HardwareAddr(intfHwAddr))
}

func NewClient(listenAddr, targetAddr *net.UDPAddr, sourceHwAddr net.HardwareAddr) (*Client, error) {
	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	// to avoid ignore msg, set random sequence number
	rand.Seed(time.Now().UnixNano())
	seq := uint16(rand.Intn(0xffff))

	return &Client{
		listenAddr:   listenAddr,
		targetAddr:   targetAddr,
		sourceHwAddr: sourceHwAddr,
		conn:         conn,
		seq:          seq,
	}, nil
}

func (c *Client) SendRecvMsg(msg nsdp.Msg) (*nsdp.Msg, error) {
	c.seq = c.seq + 1

	recvCh := make(chan bool, 1)
	buf := make([]byte, DefaultReceiveBufferSize)
	readLen := 0
	go func() {
		readLen, _, _ = c.conn.ReadFrom(buf)
		log.Println("recv", readLen, buf[:readLen])
		recvCh <- true
	}()

	retry := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for retry < 3 {
		select {
		case <-recvCh:
			resp := nsdp.ParseMsg(buf[:readLen])
			if resp == nil {
				return resp, errors.New("Failed to respose message parse")
			}
			return resp, nil
		case <-ticker.C:
			b := msg.Bytes()
			log.Println("send", b)
			_, err := c.conn.WriteTo(b, c.targetAddr)
			if err != nil {
				return nil, err
			}
			retry += 1
		}
	}

	return nil, errors.New("Failed to wait response")
}

func (c *Client) Login(password string) error {
	if password == "" {
		return errors.New("empty password")
	}

	resp, err := c.Read(nsdp.AuthV2PasswordSalt{})
	if err != nil {
		return err
	}

	mac := []byte(resp.Header.DeviceMac[:6]) // to format log clearly
	salt := resp.Body.Body[0].(*nsdp.AuthV2PasswordSalt).BytesValue
	encodedPassword := nsdp.CalcAuthV2Password(password, mac, salt)
	auth := nsdp.AuthV2Password{BytesValue: nsdp.BytesValue(encodedPassword)}

	log.Printf("%s, %x, %x", password, mac, salt)

	resp, err = c.Write(auth)
	if err != nil {
		return err
	}

	if resp.Result != 0 {
		return errors.New("switch returns login error")
	}

	c.password = encodedPassword
	c.targetHwAddr = mac

	return nil
}

func (c *Client) Read(msg ...nsdp.TLV) (*nsdp.Msg, error) {
	m := nsdp.Msg(nsdp.DefaultMsg)
	m.Op = 1
	m.Seq = c.seq
	m.HostMac = c.sourceHwAddr
	m.Body = nsdp.Body{Body: msg}

	return c.SendRecvMsg(m)
}

func (c *Client) WriteWithAuth(msg ...nsdp.TLV) (*nsdp.Msg, error) {
	if len(c.password) == 0 {
		return nil, errors.New("maybe write operation need password")
	}
	auth := nsdp.AuthV2Password{BytesValue: nsdp.BytesValue(c.password)}
	msg = append([]nsdp.TLV{auth}, msg...)

	return c.Write(msg...)
}

func (c *Client) Write(msg ...nsdp.TLV) (*nsdp.Msg, error) {
	m := nsdp.Msg(nsdp.DefaultMsg)
	m.Op = 3
	m.Seq = c.seq
	m.HostMac = c.sourceHwAddr
	m.Body = nsdp.Body{Body: msg}
	if len(c.targetHwAddr) != 0 {
		m.DeviceMac = c.targetHwAddr
	}

	log.Println(m)
	return c.SendRecvMsg(m)
}
