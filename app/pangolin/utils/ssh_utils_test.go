package utils

import (
	"golang.org/x/crypto/ssh"
	"net"
	"testing"
	"time"
)

func TestSSH(t *testing.T) {
	config := &ssh.ClientConfig{
		User:            "konghaishuo",
		Auth:            []ssh.AuthMethod{publicKeyAuthFunc("/Users/konghaishuo/.ssh/id_rsa")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	addr := "10.189.100.29:22"
	listener, err := net.Listen("tcp", "0.0.0.0:5432")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		localConn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			forward(localConn, addr, config)
		}()
	}

}

func forward(localConn net.Conn, addr string, config *ssh.ClientConfig) {
}

func TestSSHTunnel_Start(t *testing.T) {
	tunnel := &SSHTunnel{
		Local: &Endpoint{
			Host: "",
			Port: 45454,
		},
		Server: &Endpoint{
			Host: "10.189.100.29",
			Port: 22,
		},
		Remote: &Endpoint{
			Host: "10.9.173.80",
			Port: 5432,
		},
		Config: &ssh.ClientConfig{
			User:            "konghaishuo",
			Auth:            []ssh.AuthMethod{publicKeyAuthFunc("/Users/konghaishuo/.ssh/id_rsa")},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		},
	}
	tunnel.Start()
}
