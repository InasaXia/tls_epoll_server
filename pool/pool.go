package pool

import (
	"tls_epoll_server/epoll"
	logger "tls_epoll_server/log"
	"io"
	"net"
	"sync"
)

type Pool struct {
	workers   int
	maxTasks  int
	taskQueue chan net.Conn
	mu     *sync.Mutex
	closed bool
	done   chan struct{}
	epoller *epoll.Epoll
}

func NewPool(w *int, t *int,epoll2 *epoll.Epoll) *Pool {
	return &Pool{
		workers:   *w,
		maxTasks:  *t,
		taskQueue: make(chan net.Conn, *t),
		done:      make(chan struct{}),
		mu: 	   &sync.Mutex{},
		epoller:   epoll2,
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *Pool) AddTask(conn net.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()
	p.taskQueue <- conn
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *Pool) startWorker() {
	buf := make([]byte,4096)
	for {
		select {
		case <-p.done:
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				n,err := conn.Read(buf)
				if err==io.EOF {
					p.epoller.Remove(conn)
					//log.Printf("Client exit ...\n")
					logger.Logger.Infof("Client exit ...\n")
					continue
				}else {
					tmp := string(buf[:n])
					//log.Printf("Recv : %+v",tmp)
					logger.Logger.Infof("Recv : %+v",tmp)
				}
			}
		}
	}
}


