package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHTunnel struct {
	Local    *Endpoint
	Server   *Endpoint
	Remote   *Endpoint
	Config   *ssh.ClientConfig
	errChan  chan error
	shutdown *bool
}

// 通过密钥连接：
func PrivateKeyFile(file string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(signer)
}

func PrivateKeyString(key string) (ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (tunnel *SSHTunnel) Shutdown() {
	boolValue := true
	tunnel.shutdown = &boolValue
}

func (tunnel *SSHTunnel) Start() {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go tunnel.forward(conn)
		if tunnel.shutdown != nil {
			break
		}
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		return
	}
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
	for {
		select {
		case <-time.Tick(1 * time.Second):
			if tunnel.shutdown != nil {
				serverConn.Close()
				remoteConn.Close()
			}
		}
	}
}
