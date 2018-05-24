package frontend

import (
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
)

type ProdFrontPeerFetcher struct {
}

func NewPeerFetcher() *ProdFrontPeerFetcher {
	return new(ProdFrontPeerFetcher)
}

func (this ProdFrontPeerFetcher) FetchPeer() *Peer {
	// online
	return new(Peer).SetHost(cnst.FE_IP_ONLINE).SetPort(cnst.FE_PORT_ONLINE)
	// staging
	//return new(Peer).SetHost(cnst.FE_IP_STAGING).SetPort(cnst.FE_PORT_STAGING)
}
