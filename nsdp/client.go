package nsdp

import (
	"errors"
	"math/rand"
	"net"
	"time"
)

const (
	DefaultDestAddr          = "255.255.255.255"
	DefaultRecvPort          = "63321"
	DefaultSendPort          = "63322"
	DefaultReceiveBufferSize = 0xffff
	DefaultTransmitInterval  = time.Millisecond * 300
	DefaultTransmitRetry     = 5
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

func (c *Client) SendRecvMsg(msg *Msg) (*Msg, error) {
	c.seq = c.seq + 1

	recvCh := make(chan bool, 1)
	buf := make([]byte, DefaultReceiveBufferSize)
	readLen := 0
	go func() {
		readLen, _, _ = c.conn.ReadFrom(buf)
		recvCh <- true
	}()

	retry := 0
	ticker := time.NewTicker(DefaultTransmitInterval)
	defer ticker.Stop()
	for retry < DefaultTransmitRetry {
		select {
		case <-recvCh:
			resp, err := NewMsgFromBinary(buf[:readLen])
			if resp == nil || err != nil {
				return resp, errors.New("Failed to respose message parse")
			}
			return resp, nil
		case <-ticker.C:
			b, err := msg.MarshalBinary()
			_, err = c.conn.WriteTo(b, c.targetAddr)
			if err != nil {
				return nil, err
			}
			retry += 1
		}
	}

	return nil, errors.New("Failed to wait response")
}

func (c Client) makeReadMsg(tlvs ...TLV) *Msg {
	m := Msg(DefaultMsg)
	m.Op = 1
	m.Seq = c.seq
	m.HostMac = c.sourceHwAddr
	m.Body = Body(tlvs)

	return &m
}

func (c *Client) Read(tlvs ...TLV) (*Msg, error) {
	return c.SendRecvMsg(c.makeReadMsg(tlvs...))
}

func (c Client) makeWriteMsg(tlvs ...TLV) *Msg {
	m := Msg(DefaultMsg)
	m.Op = 3
	m.Seq = c.seq
	m.HostMac = c.sourceHwAddr
	m.Body = Body(tlvs)
	return &m
}

func (c *Client) Write(tlvs ...TLV) (*Msg, error) {
	return c.SendRecvMsg(c.makeWriteMsg(tlvs...))
}

func (c *Client) WriteWithAuth(password string, tlvs ...TLV) (*Msg, error) {
	if len(password) == 0 {
		return nil, errors.New("maybe write operation need password")
	}

	resp, err := c.Read(AuthV2PasswordSalt{})
	if err != nil {
		return nil, err
	}

	mac := []byte(resp.Header.DeviceMac[:6]) // to format log clearly
	salt := resp.Body[0].(*AuthV2PasswordSalt).BytesValue
	encodedPassword := CalcAuthV2Password(password, mac, salt)
	auth := AuthV2Password{BytesValue: BytesValue(encodedPassword)}

	msg := c.makeWriteMsg()
	msg.DeviceMac = mac
	msg.Body = append(msg.Body, auth)
	msg.Body = append(msg.Body, tlvs...)
	return c.SendRecvMsg(msg)
}
