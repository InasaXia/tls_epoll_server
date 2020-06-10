// +build linux
package epoll

import (
	logger "tls_epoll_server/log"
	"crypto/tls"
	_ "log"
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type Epoll struct {
	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

func MkEpoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]net.Conn),
	}, nil
}

func (e *Epoll) Add(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.connections[fd] = conn
	//log.Printf("Epoll add success:%+v;fd:%d;connections:%d\n",conn.RemoteAddr(),fd,len(e.connections))
	logger.Logger.Infof("Epoll add success:%+v;fd:%d;connections:%d\n",conn.RemoteAddr(),fd,len(e.connections))
	return nil
}

func (e *Epoll) Remove(conn net.Conn) error {
	//log.Printf("Epoll remove ：%+v;%p\n", conn.RemoteAddr(),conn)
	logger.Logger.Infof("Epoll remove ：%+v;%p\n", conn.RemoteAddr(),conn)
	fd := socketFD(conn)
	//log.Printf("Epoll remove fd: %d\n", fd)
	logger.Logger.Infof("Epoll remove fd: %d\n", fd)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		//log.Printf("Epoll remove error :%+v,conn:%+v\n", err, conn.RemoteAddr())
		logger.Logger.Errorf("Epoll remove error :%+v,conn:%+v\n", err, conn.RemoteAddr())
	}else {
		e.lock.Lock()
		defer e.lock.Unlock()
		delete(e.connections, fd)
		//log.Printf("Epoll remove :%+v;connections:%d\n", conn.RemoteAddr(), len(e.connections))
		logger.Logger.Infof("Epoll remove :%+v;connections:%d\n", conn.RemoteAddr(), len(e.connections))
	}
	return nil
}

func (e *Epoll) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 1000000)
	n, err := unix.EpollWait(e.fd, events, 1000000)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := e.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func socketFD(conn net.Conn) int {
	switch conn.(type) {
	case *tls.Conn:
		return tlsConnFD(conn)
	case *net.TCPConn:
		return tcpConnFD(conn)
	default:
		return -1
	}
}
func tcpConnFD(conn net.Conn) int {
	//edit Conn.go in tls package
	//add this function
	//func (c *Conn) UnderlyingConn() net.Conn {
	//	return c.conn
	//}
	netFD := reflect.Indirect(reflect.Indirect(reflect.ValueOf(conn)).FieldByName("fd"))
	pfd := reflect.Indirect(netFD.FieldByName("pfd"))
	fd := int(pfd.FieldByName("Sysfd").Int())
	return fd
}
func tlsConnFD(conn net.Conn) int {
	tlsConn,ok := conn.(*tls.Conn)
	if ok {
		tcpConn := tlsConn.UnderlyingConn()
		return tcpConnFD(tcpConn)
	}else {
		return -1
	}
}
