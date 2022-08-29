package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

const tcpNetwork = "tcp"

var (
	ErrNoConnection     = errors.New("no connection was made")
	ErrConnectionClosed = errors.New("connection closed by peer")
)

type TelnetClient struct {
	network, addr string
	timeout       time.Duration
	in            io.ReadCloser
	out           io.Writer
	conn          net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return TelnetClient{
		network: tcpNetwork,
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout(t.network, t.addr, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *TelnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return ErrNoConnection
}

func (t *TelnetClient) Send() error {
	if t.conn == nil {
		return ErrNoConnection
	}
	buf := make([]byte, 4096)
	reader := bufio.NewReader(t.in)
	bytesRead, err := reader.Read(buf)
	if err != nil {
		return err
	}
	if bytesRead > 0 {
		_, err := t.conn.Write(buf[:bytesRead])
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TelnetClient) Receive() error {
	if t.conn == nil {
		return ErrNoConnection
	}
	buf := make([]byte, 4096)
	reader := bufio.NewReader(t.conn)
	bytesRead, err := reader.Read(buf)
	if err != nil {
		return ErrConnectionClosed
	}
	_, err = t.out.Write(buf[:bytesRead])
	if err != nil {
		return err
	}

	return nil
}
