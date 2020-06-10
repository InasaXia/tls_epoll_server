package main

import (
	"tls_epoll_server/config"
	logger "tls_epoll_server/log"
	"tls_epoll_server/pool"
	"tls_epoll_server/socket"
	"tls_epoll_server/epoll"
	"flag"
	_ "log"
	"net"
)

var Epoller *epoll.Epoll
var workerPool *pool.Pool
var configFilePath string
var cfg *config.Config

func init()  {
	flag.StringVar(&configFilePath,"c","/etc/config.yaml","config file path")
	flag.Parse()
	cfg = new(config.Config)
	Epoller, _ = epoll.MkEpoll()
}
func startWorkPool(wp *config.WorkPool) {
	workerPool = pool.NewPool(&wp.Capacity,&wp.Queue,Epoller)
	workerPool.Start()

}
func main() {
	cfg.Init(configFilePath)
	ln := socket.StartTcpServer(cfg)
	startWorkPool(&cfg.WorkPool)
	go startEpollWait()
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				//log.Printf("accept temp err: %v", ne)
				logger.Logger.Errorf("accept temp err: %v", ne)
				continue
			}
			//log.Printf("accept err: %v", e)
			logger.Logger.Infof("accept err: %v", e)
			return
		}
		if err := Epoller.Add(conn); err != nil {
			//log.Printf("failed to add connection %v", err)
			logger.Logger.Errorf("failed to add connection %v", err)
			conn.Close()
		}
	}
	workerPool.Close()
}
func startEpollWait() {
	for {
		connections, err := Epoller.Wait()
		if err != nil {
			//log.Printf("failed to epoll wait %v", err)
			logger.Logger.Errorf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			workerPool.AddTask(conn)
		}
	}
}
