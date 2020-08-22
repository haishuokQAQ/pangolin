package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const (
	StatisticTypeUp   = "up"
	StatisticTypeDown = "down"
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
	Local         *Endpoint
	Server        *Endpoint
	Remote        *Endpoint
	Config        *ssh.ClientConfig
	errChan       chan error
	shutdown      *bool
	flowStatistic map[string]*FlowStatistic
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
	tunnel.flowStatistic = map[string]*FlowStatistic{}
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
	copyConn := func(writer, reader net.Conn, statisticKey string) {
		_, err := tunnel.innerCopy(writer, reader, statisticKey)
		if err != nil {
		}
	}
	go copyConn(localConn, remoteConn, StatisticTypeDown)
	go copyConn(remoteConn, localConn, StatisticTypeUp)
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

func (tunnel *SSHTunnel) innerCopy(dst io.Writer, src io.Reader, statisticKey string) (written int64, err error) {
	statistic, ok := tunnel.flowStatistic[statisticKey]
	if !ok {
		statistic = &FlowStatistic{}
		tunnel.flowStatistic[statisticKey] = statistic
	}
	var buf []byte
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		statisticAddress := &statistic.Bytes

		atomic.AddUint64(statisticAddress, uint64(nr))
	}
	return written, err
}

func (tunnel *SSHTunnel) GetStatistic() map[string]*FlowStatistic {
	return tunnel.flowStatistic
}

type FlowStatistic struct {
	Bytes uint64 `json:"bytes"`
}
