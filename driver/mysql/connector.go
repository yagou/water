package mysql

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	errLog *log.Logger

	errBadConn = errors.New("driver: bad connection")
)

type mysql struct {
	buf      *buffer
	netConn  net.Conn
	sequence uint8
	flags    clientFlag
}

func GetMysqlConnector(address, username, password string, timeout time.Duration) (*mysql, error) {

	nd := net.Dialer{Timeout: timeout}
	conn, err := nd.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	my := &mysql{}

	my.buf = newBuffer(conn)
	cipher, err := my.readInitPacket()
	if err != nil {
		return nil, err
	}

	data := make([]byte, 1024)
	data[0] = 1
	n, err := my.netConn.Write(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(cipher))
	fmt.Println(n)
	return my, nil
}

func (m *mysql) Close() error {
	if m.netConn != nil {
		// TODO 此处应给mysql服务器发送断开命令

		m.netConn.Close()
		m.netConn = nil
	}
	m.buf = nil
	return nil
}
