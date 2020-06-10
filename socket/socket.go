package socket

import (
	"tls_epoll_server/config"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net"
	"strconv"
)

func StartTcpServer(c *config.Config) net.Listener {
	tlsCfg := c.TLS
	if tlsCfg.Ca==""&&tlsCfg.Crt==""&&tlsCfg.Key=="" {
		return tcpNoTLS(&c.Server)
	}else {
		return tcpWithTLS(&c.Server,&c.TLS)
	}
}

func tcpNoTLS(s *config.Server) net.Listener {
	listen,err := net.Listen("tcp",s.Address+":"+strconv.Itoa(s.Port))
	if err!=nil {
		panic(err)
	}
	return listen
}
func tcpWithTLS(s *config.Server,t *config.TLS) net.Listener {
	if t.Phrase=="" {
		return tcpWithTLSNoPharase(s,t)
	}else {
		return tcpWithTLSAndPhrase(s,t)
	}
}
func tcpWithTLSAndPhrase(s *config.Server,t *config.TLS) net.Listener {
	pool := x509.NewCertPool()
	caCrt,err := ioutil.ReadFile(t.Ca)
	if err!=nil {
		panic(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	keyByte,err := ioutil.ReadFile(t.Key)
	certS,err := ioutil.ReadFile(t.Crt)
	keyBlock,_ := pem.Decode(keyByte)
	keyDER,err := x509.DecryptPEMBlock(keyBlock,[]byte(t.Phrase))
	if err!=nil {
		panic(err)
	}
	keyBlock.Bytes=keyDER
	keyBlock.Headers=nil
	keyPem := pem.EncodeToMemory(keyBlock)
	certificate,err := tls.X509KeyPair(certS,keyPem)
	if err!=nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          pool,
		InsecureSkipVerify: false,
	}
	listen,err := tls.Listen("tcp",s.Address+":"+strconv.Itoa(s.Port),tlsConfig)
	if err!=nil {
		panic(err)
	}
	return listen
}
func tcpWithTLSNoPharase(s *config.Server,t *config.TLS) net.Listener {
	pool := x509.NewCertPool()
	caCrt,err := ioutil.ReadFile(t.Ca)
	if err!=nil {
		panic(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	certificate,err := tls.LoadX509KeyPair(t.Crt,t.Key)
	if err!=nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          pool,
		InsecureSkipVerify: false,
	}
	listen,err := tls.Listen("tcp",s.Address+":"+strconv.Itoa(s.Port),tlsConfig)
	if err!=nil {
		panic(err)
	}
	return listen
}


