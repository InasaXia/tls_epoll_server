# tls_epoll_server
a simple tls epoll server 

Since tls.Conn do not allow access to file descriptor(https://github.com/golang/go/issues/29257), I added a method in the tls package

```go
// Get Underlying Connection
func (c *Conn) UnderlyingConn() net.Conn {
	return c.conn
}
```
