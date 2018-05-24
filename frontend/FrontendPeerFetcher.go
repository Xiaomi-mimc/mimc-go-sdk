package frontend

import (
	"strconv"
)

type IFrontendPeerFetcher interface {
	FetchPeer() *Peer
}

type Peer struct {
	host string
	port int
}

func (this *Peer) SetHost(host string) *Peer {
	this.host = host
	return this
}
func (this *Peer) SetPort(port int) *Peer {
	this.port = port
	return this
}
func (this *Peer) Host() string {
	return this.host
}
func (this *Peer) Port() int {
	return this.port
}

func (this *Peer) ToString() string {
	return this.host + ":" + strconv.Itoa(this.port)
}
