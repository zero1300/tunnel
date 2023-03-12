package test

import (
	"io"
	"net"
	"testing"
)

const ServerAddr = "0.0.0.0:8899"
const TunnelAddr = "0.0.0.0:8900"
const BufSize = 10

// server
func TestServer(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", ServerAddr)
	if err != nil {
		t.Fatal(err.Error())
	}
	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err.Error())
	}
	var buf [BufSize]byte
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			t.Fatal(err.Error())
		}
		// 读数据
		b := make([]byte, 0)
		for {
			n, err := conn.Read(buf[:])
			if err != nil {
				t.Fatal(err.Error())
			}
			b = append(b, buf[:n]...)
			if n < BufSize {
				break
			}
		}
		t.Log(string(b))
		// 写数据
		_, _ = conn.Write([]byte("哔哔哔, " + string(b)))

		_ = conn.Close()

	}

}

// client
func TestClient(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", TunnelAddr)
	if err != nil {
		t.Fatal(err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	_, err = conn.Write([]byte("Hello there"))
	if err != nil {
		t.Log(err.Error())
	}
	b := make([]byte, 0)
	var buf [BufSize]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			t.Log(err)
		}
		b = append(b, buf[:n]...)
		if n < BufSize {
			break
		}
	}
	t.Log(string(b))
	_ = conn.Close()
}

// tunnel
func TestTunnel(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", TunnelAddr)
	if err != nil {
		t.Fatal(err.Error())
	}
	listener, err := net.ListenTCP("tcp", addr)

	for {
		clientConn, err := listener.AcceptTCP()

		if err != nil {
			t.Fatal(err)
		}

		SAddr, err := net.ResolveTCPAddr("tcp", ServerAddr)
		serverConn, err := net.DialTCP("tcp", nil, SAddr)
		go func() {
			_, err := io.Copy(serverConn, clientConn)
			if err != nil {
				t.Log(err)
			}
		}()
		go func() {
			_, err := io.Copy(clientConn, serverConn)
			if err != nil {
				t.Log(err)
			}
		}()

	}
}
