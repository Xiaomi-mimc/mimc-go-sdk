package mimc

import (
	"container/list"
	"github.com/Xiaomi-mimc/mimc-go-sdk/cipher"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/frontend"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/string"
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

	if this.status == NOT_CONNECTED {
		return
	}
	if this.user != nil {
		this.user.status = Offline
	}

	if this.tcpConn != nil {
		this.tcpConn.Close()
	}
	this.user.lastCreateConnTimestamp = 0
	network_error := "NETWORK_ERROR"
	this.user.status = Offline
	this.user.statusDelegate.HandleChange(false, &network_error, &network_error, &network_error)
	this.init()

}

func (this *MIMCConnection) Connect() bool {
	feAddrs := this.user.FeAddress()
	if feAddrs == nil {
		if this.user.refreshToken() {
			feAddrs = this.user.FeAddress()
		} else {
			logger.Info("[%s] fresh token failed", this.user.appAccount)
		}
	}
	for _, addr := range feAddrs {
		conn, err := net.Dial("tcp", addr)
		logger.Info("[%v] connect to server:%v", this.user.appAccount, addr)
		if err == nil {
			this.tcpConn = conn
			return true
		}
	}
	return false
}

func (this *MIMCConnection) Readn(buf *[]byte, length int) int {
	if !this.check(buf, length) {
		logger.Warn("[%v] check: buf len %v != length %v", this.user.appAccount, len(*buf), length)
		return -1
	}
	left := length
	for left > 0 {
		tmpBuf := make([]byte, left)
		nread, err := this.tcpConn.Read(tmpBuf)
		if err != nil || nread < 0 {
			logger.Error("[%v] read error. err: %v, nread: %v, length: %v", this.user.appAccount, err, nread, length)
			return -1
		}
		if nread == 0 {
			break
		}
		for i := 0; i < nread; i++ {
			(*buf)[length-left+i] = tmpBuf[i]
		}
		left = left - nread
		if left < 0 {
			logger.Debug("[%v] nread: %v, left: %v, length: %v, lenbuf: %v", this.user.appAccount, nread, left, length, len(*buf))
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
	var tmpBuf []byte
	for left > 0 {
		tmpBuf = make([]byte, left)
		for i := 0; i < left; i++ {
			tmpBuf[i] = (*buf)[length-left+i]
		}
		nwrite, err := this.tcpConn.Write(tmpBuf)
		if err != nil || nwrite < 0 {
			logger.Error("write error.")
			return -1
		}
		if nwrite == 0 {
			break
		}
		left = left - nwrite
	}
	//logger.Info("writen: %v", (length - left))
	return length - left

}

func (this *MIMCConnection) check(buf *[]byte, length int) bool {
	if this.tcpConn == nil || buf == nil || len(*buf) < length {
		logger.Debug("[%v] tcpConn: %v, buf:%v", this.user.appAccount, this.tcpConn, buf)
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
	//this.peerFetcher = NewPeerFetcher()

}
