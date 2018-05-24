package mimc

import (
	"container/list"
	"mimc-go-sdk/cipher"
	"mimc-go-sdk/common/constant"
	. "mimc-go-sdk/frontend"
	"mimc-go-sdk/util/string"
	"net"
)

type ConnStatus int

const (
	NOT_CONNECTED = iota
	SOCK_CONNECTED
	HANDSHAKE_CONNECTED
)

type MIMCConnection struct {
	tcpConn     net.Conn
	peer        *Peer
	peerFetcher IFrontendPeerFetcher
	status      ConnStatus

	rc4Key        []byte
	challenge     string
	packetsToSend *list.List

	connpt string
	model  string
	os     string
	udid   string
	sdk    string
	locale string
	andVer int

	user *MCUser

	nextResetSockTimestamp int64
	lastPingTimestamp      uint64
	tryCreateConnCount     uint8
	defaults               map[string]string
}

func (this *MIMCConnection) Rc4Key() []byte {
	return this.rc4Key
}
func (this *MIMCConnection) Status() ConnStatus {
	return this.status
}
func (this *MIMCConnection) Sock_Connected() {
	this.status = SOCK_CONNECTED
}
func (this *MIMCConnection) HandshakeConnected() *MIMCConnection {
	this.status = HANDSHAKE_CONNECTED
	return this
}

func (this *MIMCConnection) ClearSockTimestamp() *MIMCConnection {
	this.nextResetSockTimestamp = -1
	return this
}

func (this *MIMCConnection) TrySetNextResetSockTs() {
	if this.nextResetSockTimestamp > 0 {
		return
	}
	this.nextResetSockTimestamp = CurrentTimeMillis() + cnst.RESET_SOCKET_TIMEOUT_TIMEVAL_MS
}
func (this *MIMCConnection) NextResetSockTimestamp() int64 {
	return this.nextResetSockTimestamp
}

func (this *MIMCConnection) PeerFetcher(peerFetcher IFrontendPeerFetcher) *MIMCConnection {
	this.peerFetcher = peerFetcher
	return this
}

func (this *MIMCConnection) User(user *MCUser) *MIMCConnection {
	this.user = user
	return this
}

func (this *MIMCConnection) Challenge() string {
	return this.challenge
}
func (this *MIMCConnection) SetChallenge(challenge string) {
	this.challenge = challenge
}
func (this *MIMCConnection) Udid() string {
	return this.udid
}

func (this *MIMCConnection) Reset() {
	if this.user != nil {
		this.user.status = Offline
	}

	if this.tcpConn != nil {
		this.tcpConn.Close()
	}
	logger.Info("reset conn.")
	this.init()
	if !this.Connect() {
		logger.Debug("connet to MIMC Server fail.")
	} else {
		logger.Debug("connet to MIMC Server succ.")
	}

}

func (this *MIMCConnection) Connect() bool {
	if this.peerFetcher == nil {
		logger.Warn("peerFetcher is nil.")
		return false
	}
	this.peer = this.peerFetcher.FetchPeer()
	conn, err := net.Dial("tcp", this.peer.ToString())
	if err == nil {
		this.tcpConn = conn
		return true
	}
	return false
}

func (this *MIMCConnection) Readn(buf *[]byte, length int) int {
	if !this.check(buf, length) {
		logger.Warn("check: buf len %v != length %v", len(*buf), length)
		return -1
	}
	left := length
	for left > 0 {
		nread, err := this.tcpConn.Read(*buf)
		if err != nil || nread < 0 {
			logger.Error("read error.\n")
			return -1
		}
		if nread == 0 {
			break
		}
		left = left - nread
		if left < 0 {
			logger.Debug("nread: %v, left: %v, length: %v, lenbuf: %v", nread, left, length, len(*buf))
			return length
		}
	}
	return length - left
}
func (this *MIMCConnection) Writen(buf *[]byte, length int) int {
	if !this.check(buf, length) {
		return -1
	}

	left := length
	for left > 0 {
		nwrite, err := this.tcpConn.Write(*buf)
		if err != nil || nwrite < 0 {
			logger.Error("write error.")
			return -1
		}
		if nwrite == 0 {
			break
		}
		left = left - nwrite
	}
	return length - left

}

func (this *MIMCConnection) check(buf *[]byte, length int) bool {
	if this.tcpConn == nil || buf == nil || len(*buf) < length {
		logger.Debug("tcpConn: %v, buf:%v", this.tcpConn, buf)
		return false
	}
	return true
}

func (this *MIMCConnection) SetChallengeAndRc4Key(challenge string) {
	this.challenge = challenge
	halfUdid := strutil.Substring(&this.udid, len(this.udid)/2)
	halfChallenge := strutil.Substring(&this.challenge, len(this.challenge)/2)
	key := strutil.Concat(&halfChallenge, &halfUdid)
	this.rc4Key = cipher.Encrypt(strutil.Bytes(&this.challenge), strutil.Bytes(&key))
}

func NewConn() *MIMCConnection {
	conn := new(MIMCConnection)
	conn.init()
	return conn
}

func (this *MIMCConnection) init() {
	this.status = NOT_CONNECTED
	this.rc4Key = nil
	this.connpt = ""
	this.model = ""
	this.os = ""
	this.udid = ""
	this.sdk = ""
	this.locale = ""
	this.andVer = 0
	this.lastPingTimestamp = 0
	this.nextResetSockTimestamp = -1
	this.tryCreateConnCount = 0
	if this.user != nil {
		this.user = nil
	}
	this.peerFetcher = NewPeerFetcher()

}