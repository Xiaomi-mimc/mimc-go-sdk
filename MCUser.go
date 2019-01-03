package mimc

import (
	"container/list"
	"encoding/json"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"github.com/Xiaomi-mimc/mimc-go-sdk/frontend"
	"github.com/Xiaomi-mimc/mimc-go-sdk/message"
	"github.com/Xiaomi-mimc/mimc-go-sdk/packet"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/ims"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/mimc"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/byte"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/map"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/net"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/queue"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/string"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type UserStatus int

var logger = log.GetLogger()

const (
	Online UserStatus = iota
	Offline
)

func init() {

}

type MCUser struct {
	chid     int32
	uuid     int64
	resource string
	status   UserStatus

	isLogout bool

	clientAttrs string
	cloudAttrs  string

	appId      int64
	appAccount string
	appPackage string

	prefix  string
	indexer int

	securityKey string
	token       string
	tryLogin    bool

	feDomain  string
	feAddress []string

	sequenceReceived        map[uint32]interface{}
	conn                    *MIMCConnection
	lastLoginTimestamp      int64
	lastCreateConnTimestamp int64
	lastPingTimestamp       int64

	tokenDelegate  Token
	statusDelegate StatusDelegate
	msgDelegate    MessageHandlerDelegate

	messageToSend    *que.ConQueue
	messageToAck     *cmap.ConMap
	packetToCallback *que.ConQueue
}

func (this *MCUser) printUserInfo() {
	logger.Warn("\tappId:%v\n"+
		"\tappAccount:%v\n"+
		"\ttoken:%v\n"+
		"\tfeDomain:%v\n"+
		"\tfeAddress:%v\n"+
		"\tuuid:%v\n"+
		"\tsecurityKey:%v\n"+
		"\tresource:%v\n"+
		"\tchid:%v\n"+
		"\tappPackage:%v\n",
		this.appId, this.appAccount, this.token, this.feDomain, this.feAddress,
		this.uuid, this.securityKey, this.resource, this.chid, this.appPackage)
}

func NewUser(appId int64, appAccount string) *MCUser {
	this := NewMCUser()
	this.appAccount = appAccount
	this.appId = appId
	return this
}

func NewMCUser() *MCUser {
	this := new(MCUser)
	return this
}

func (this *MCUser) fetchFEAddr() []string {
	return this.feAddress
}

func (this *MCUser) RegisterTokenDelegate(tokenDelegate Token) *MCUser {
	this.tokenDelegate = tokenDelegate
	return this
}

func (this *MCUser) RegisterStatusDelegate(statusDelegate StatusDelegate) *MCUser {
	this.statusDelegate = statusDelegate
	return this
}

func (this *MCUser) RegisterMessageDelegate(msgDelegate MessageHandlerDelegate) *MCUser {
	this.msgDelegate = msgDelegate
	return this
}

func (this *MCUser) InitAndSetup() {
	void := ""
	this.status = Offline
	this.resource = "mimc_go_" + strutil.RandomStrWithLength(10)
	this.lastLoginTimestamp = 0
	this.lastCreateConnTimestamp = 0
	this.lastPingTimestamp = 0
	this.conn = NewConn().User(this)
	this.messageToSend = que.NewConQueue()
	this.messageToAck = cmap.NewConMap()
	this.packetToCallback = que.NewConQueue()
	this.appPackage = void
	this.chid = 0
	this.uuid = 0
	this.token = void
	this.securityKey = void
	this.clientAttrs = void
	this.cloudAttrs = void
	this.tryLogin = false
	this.feAddress = nil
	this.feDomain = void
	this.fetchUserInfo()
	go this.sendRoutine()
	go this.receiveRoutine()
	go this.triggerRoutine()
	go this.callBackRoutine()
}

func (this *MCUser) fetchUserInfo() {
	root, _ := exec.LookPath(os.Args[0])
	dir := cnst.CACHE_DIR
	key := strconv.FormatInt(this.appId, 10) + "_" + this.appAccount
	file := key + cnst.CACHE_FILE
	userInfo := map[string]interface{}{}

	if strutil.FetchUserInfo(&root, &dir, &file, userInfo) {
		this.token = userInfo["token"].(string)
		this.feDomain = userInfo["feDomain"].(string)
		this.feAddress = strutil.Transfer(userInfo["feAddress"].([]interface{}))
		resource := userInfo["resource"].(string)
		if strings.Compare(resource, "") != 0 {
			this.resource = resource
		}
		this.securityKey = userInfo["securityKey"].(string)
		chid, _ := userInfo["chid"].(json.Number).Int64()
		this.chid = int32(chid)
		this.appPackage = userInfo["appPackage"].(string)
		this.uuid, _ = userInfo["uuid"].(json.Number).Int64()

		if strings.Compare(this.feDomain, "") == 0 ||
			len(this.feAddress) == 0 ||
			this.uuid == 0 ||
			this.chid == 0 ||
			strings.Compare(this.token, "") == 0 {
			this.uuid = 0
			this.chid = 0
			this.token = ""
			this.feDomain = ""
			this.feAddress = []string{}
		}
	}
	//this.printUserInfo()
}

func (this *MCUser) flushUserInfo() bool {

	root, _ := exec.LookPath(os.Args[0])
	dir := cnst.CACHE_DIR
	key := strconv.FormatInt(this.appId, 10) + "_" + this.appAccount
	file := key + cnst.CACHE_FILE
	userInfo := map[string]interface{}{}
	userInfo["uuid"] = this.uuid
	userInfo["appAccount"] = this.appAccount
	userInfo["token"] = this.token
	userInfo["appId"] = this.appId
	userInfo["feDomain"] = this.feDomain
	userInfo["feAddress"] = this.feAddress
	userInfo["resource"] = this.resource
	userInfo["securityKey"] = this.securityKey
	userInfo["chid"] = this.chid
	userInfo["appPackage"] = this.appPackage
	return strutil.FlushUserInfo(&root, &dir, &file, userInfo)

}

/**
* when token is invalid, or conn get empyty feDomain/feAddress, it should refresh token
 */
func (this *MCUser) refreshToken() bool {
	logger.Info("[%v] refresh token", this.appAccount)
	if this.tokenDelegate == nil {
		logger.Error("%v Login fail, have to fetch token.", this.appAccount)
		return false
	}
	tokenJsonStr := this.tokenDelegate.FetchToken()
	if tokenJsonStr == nil {
		logger.Warn("%v Login fail, get nil token string.", this.appAccount)
		return false
	}
	var tokenMap map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(*tokenJsonStr))
	decoder.UseNumber()

	if err := decoder.Decode(&tokenMap); err == nil {
		code, _ := tokenMap["code"].(json.Number).Int64()
		if code != 200 {
			logger.Warn("%v fetch token fail, response code: %v", this.appAccount, tokenMap["message"].(string))
			return false
		}
		data := tokenMap["data"].(map[string]interface{})
		appId, _ := strconv.ParseInt(data["appId"].(string), 10, 64)
		if appId != this.appId {
			logger.Warn("appId:%v in token no match appId: %v.", appId, this.appId)
			return false
		}
		appAccount := data["appAccount"].(string)
		if appAccount != this.appAccount {
			logger.Warn("appAccount:%v in token not match appAccount: %v.", appAccount, this.appAccount)
			return false
		}
		this.appPackage = data["appPackage"].(string)
		chid, _ := data["miChid"].(json.Number).Int64()
		this.chid = int32(chid)
		uuid, err := strconv.ParseInt(data["miUserId"].(string), 10, 64)
		if err != nil {
			logger.Error("%v Login fail, can not parse token string.", this.appAccount)
			return false
		}
		this.uuid = uuid
		this.securityKey = data["miUserSecurityKey"].(string)
		token, ok := data["token"]
		this.feDomain = data["feDomainName"].(string)
		this.feAddress = httputil.GetFEAddress(cnst.ONLINE_RESOLVER_URL, this.feDomain)
		if ok {
			this.token = token.(string)
			//this.printUserInfo()
			return this.flushUserInfo()
		} else {
			logger.Warn("parse token failed")
			return false
		}
	} else {
		return false
	}
}

func (this *MCUser) Login() bool {
	this.tryLogin = true

	return true
}
func (this *MCUser) Logout() bool {
	if this.status == Offline {
		return false
	}
	v6PacketForUnbind := BuildUnBindPacket(this)
	unBindPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6PacketForUnbind)
	this.messageToSend.Push(unBindPacket)
	this.tryLogin = false
	return true
}

func (this *MCUser) SendMessage(toAppAccount string, payload []byte) string {
	if &toAppAccount == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[%v] [Send P2P Msg]%v -> %v: %v.", this.appAccount, this.appAccount, toAppAccount, string(payload))
	v6Packet, mimcPacket := BuildP2PMessagePacket(this, toAppAccount, payload, true, nil)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) SendMessageWithStore(toAppAccount string, payload []byte, isStore bool) string {
	if &toAppAccount == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[%v] [Send P2P Msg]%v -> %v: %v.", this.appAccount, this.appAccount, toAppAccount, string(payload))
	v6Packet, mimcPacket := BuildP2PMessagePacket(this, toAppAccount, payload, isStore, nil)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) SendMessageWithBizType(toAppAccount string, payload []byte, bizType *string) string {
	if &toAppAccount == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[%v] [Send P2P Msg]%v -> %v: %v.", this.appAccount, this.appAccount, toAppAccount, string(payload))
	v6Packet, mimcPacket := BuildP2PMessagePacket(this, toAppAccount, payload, true, bizType)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}
func (this *MCUser) SendMessageWithStoreAndBizType(toAppAccount string, payload []byte, bizType *string, isStore bool) string {
	if &toAppAccount == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[%v] [Send P2P Msg]%v -> %v: %v.", this.appAccount, this.appAccount, toAppAccount, string(payload))
	v6Packet, mimcPacket := BuildP2PMessagePacket(this, toAppAccount, payload, isStore, bizType)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) SendGroupMessage(topicId *int64, payload []byte) string {
	if &topicId == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[Send P2T Msg]%v send p2t msg to %v: %v.\n", this.appAccount, *topicId, string(payload))
	v6Packet, mimcPacket := BuildP2TMessagePacket(this, *topicId, payload, true, nil)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}
func (this *MCUser) SendGroupMessageWithStore(topicId *int64, payload []byte, isStore bool) string {
	if &topicId == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[Send P2T Msg]%v send p2t msg to %v: %v.\n", this.appAccount, *topicId, string(payload))
	v6Packet, mimcPacket := BuildP2TMessagePacket(this, *topicId, payload, isStore, nil)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) SendGroupMessageWithBizType(topicId *int64, payload []byte, bizType *string) string {
	if &topicId == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[Send P2T Msg]%v send p2t msg to %v: %v.\n", this.appAccount, *topicId, string(payload))
	v6Packet, mimcPacket := BuildP2TMessagePacket(this, *topicId, payload, true, bizType)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) SendGroupMessageWithStoreAndBizType(topicId *int64, payload []byte, bizType *string, isStore bool) string {
	if &topicId == nil || payload == nil || len(payload) == 0 {
		return ""
	}
	logger.Info("[Send P2T Msg]%v send p2t msg to %v: %v.\n", this.appAccount, *topicId, string(payload))
	v6Packet, mimcPacket := BuildP2TMessagePacket(this, *topicId, payload, isStore, bizType)
	timeoutPacket := packet.NewTimeoutPacket(CurrentTimeMillis(), mimcPacket)
	msgPacket := msg.NewMsgPacket(cnst.MIMC_C2S_DOUBLE_DIRECTION, v6Packet)
	this.messageToSend.Push(msgPacket)
	this.messageToAck.Push(*(mimcPacket.PacketId), timeoutPacket)
	return *(mimcPacket.PacketId)
}

func (this *MCUser) sendRoutine() {
	logger.Info("[%s] initate send goroutine.", this.appAccount)
	if this.conn == nil {
		return
	}
	msgType := cnst.MIMC_C2S_DOUBLE_DIRECTION

	for {
		var pkt *packet.MIMCV6Packet = nil
		if this.conn.Status() == NOT_CONNECTED {
			logger.Debug("the conn not connected.\n")
			currTimeMillis := CurrentTimeMillis()
			if currTimeMillis-this.lastCreateConnTimestamp <= cnst.CONNECT_TIMEOUT {
				Sleep(100)
				continue
			}

			this.lastCreateConnTimestamp = CurrentTimeMillis()
			if !this.conn.Connect() {
				logger.Warn("connet to MIMC Server fail.\n")
				continue
			}
			this.conn.Sock_Connected()
			this.lastCreateConnTimestamp = 0
			logger.Info("[%v] build conn packet.", this.appAccount)
			pkt = BuildConnectionPacket(this.conn.Udid(), this)
		}
		if this.conn.Status() == SOCK_CONNECTED {
			Sleep(100)
		}
		if this.conn.Status() == HANDSHAKE_CONNECTED {
			currTimeMillis := CurrentTimeMillis()
			if this.status == Offline && currTimeMillis-this.lastLoginTimestamp <= cnst.LOGIN_TIMEOUT {
				Sleep(100)
				continue
			}
			if this.tryLogin && this.status == Offline && currTimeMillis-this.lastLoginTimestamp > cnst.LOGIN_TIMEOUT {
				logger.Debug("%v: build bind packet.", this.appAccount)
				pkt = BuildBindPacket(this)
				if pkt == nil {
					Sleep(100)
					continue
				}
				this.lastLoginTimestamp = CurrentTimeMillis()
			}
		}
		if this.status == Online {

			msgPacketToSend := this.messageToSend.Pop()
			if msgPacketToSend == nil {
				dist := CurrentTimeMillis() - this.lastPingTimestamp
				isPing := dist-cnst.PING_TIMEVAL_MS > 0
				if isPing {
					pkt = BuildPingPacket(this)
					logger.Info("[%v] build ping packet.", this.appAccount)
				} else {
					Sleep(100)
					continue
				}
			} else {
				msgPacket := msgPacketToSend.(*msg.MsgPacket)
				msgType = msgPacket.MsgType()
				pkt = msgPacket.Packet()
				logger.Debug("%v: send msg packet.", this.appAccount)

			}

		} else {

		}
		if pkt == nil {
			Sleep(100)
			continue
		}
		if msgType == cnst.MIMC_C2S_DOUBLE_DIRECTION {
			this.conn.TrySetNextResetSockTs()
		}
		payloadKey := PayloadKey(this.securityKey, pkt.HeaderId())
		bodyKey := this.conn.Rc4Key()
		packetData := pkt.Bytes(bodyKey, payloadKey)
		this.lastPingTimestamp = CurrentTimeMillis()
		size := len(packetData)
		if this.Conn().Writen(&packetData, size) != size {
			logger.Error("write data error.")
			this.conn.Reset()
		} else {
			if pkt.GetHeader() != nil {
				logger.Debug("[send]: send packet: %v succ.\n", *(pkt.GetHeader().Id))
			} else {
				logger.Debug("[send]: send packet succ.\n")
			}

		}
	}
}
func (this *MCUser) PeerFetcher(fetcher frontend.ProdFrontPeerFetcher) {
	this.conn.PeerFetcher(fetcher)
}
func (this *MCUser) receiveRoutine() {
	logger.Info("[%s] initate receive goroutine.", this.appAccount)
	var counter int = 0
	if this.conn == nil {
		return
	}
	for {
		if this.conn.Status() == NOT_CONNECTED {
			Sleep(1000)
			continue
		}
		headerBins := make([]byte, cnst.V6_HEAD_LENGTH)
		length := this.conn.Readn(&headerBins, int(cnst.V6_HEAD_LENGTH))
		if length != int(cnst.V6_HEAD_LENGTH) {
			logger.Error("%v->[rcv]: error head. need length: %v, read length: %v\n", this.appAccount, cnst.V6_HEAD_LENGTH, length)
			this.conn.Reset()
			Sleep(1000)
			continue

		}
		magic := byteutil.GetUint16FromBytes(&headerBins, cnst.V6_MAGIC_OFFSET)
		if magic != cnst.MAGIC {
			logger.Error("%v->[rcv]: error magic: %v.", this.appAccount, magic)
			this.conn.Reset()
			continue
		}
		version := byteutil.GetUint16FromBytes(&headerBins, cnst.V6_VERSION_OFFSET)
		if version != cnst.V6_VERSION {
			logger.Error("%v->[rcv]: error version: %v.", this.appAccount, version)
			this.conn.Reset()
			continue
		}
		bodyLen := byteutil.GetIntFromBytes(&headerBins, cnst.V6_BODYLEN_OFFSET)
		if bodyLen < 0 {
			logger.Error("%v->[rcv]: error bodylen: %v.", this.appAccount, bodyLen)
			this.conn.Reset()
			continue
		}
		var bodyBins []byte
		if bodyLen != 0 {
			bodyBins = make([]byte, bodyLen)
			if bodyLen != 0 {
				length = this.conn.Readn(&bodyBins, bodyLen)
				if length != bodyLen {
					logger.Error("%v->[rcv]: error body.length: %v, bodyLen:%v", this.appAccount, length, bodyLen)
					this.conn.Reset()
					continue
				} else {
					//logger.Debug("[rcv]: read.length: %v, bodyLen:%v", length, bodyLen)
				}
			}
		}
		crcBins := make([]byte, cnst.V6_CRC_LENGTH)
		crclen := this.conn.Readn(&crcBins, cnst.V6_CRC_LENGTH)
		if crclen != cnst.V6_CRC_LENGTH {
			logger.Error("%v->[rcv]: error crc: %v.", this.appAccount, crclen)
			this.conn.Reset()
			continue
		}
		this.conn.ClearSockTimestamp()
		bodyKey := this.conn.Rc4Key()
		packetBytes := packet.NewPacketBytes(&headerBins, &bodyBins, &crcBins, &bodyKey, &(this.securityKey))
		counter += 1
		this.packetToCallback.Push(packetBytes)
	}
}
func (this *MCUser) triggerRoutine() {
	logger.Info("[%s] initiate trigger goroutine.", this.appAccount)
	if this.conn == nil {
		return
	}
	for {
		nowTimeMillis := CurrentTimeMillis()
		nextRestSockTimeMillis := this.conn.NextResetSockTimestamp()
		if nextRestSockTimeMillis > 0 && nowTimeMillis-nextRestSockTimeMillis > 0 {
			logger.Warn("[trigger] wait for response timeout.")
			this.conn.Reset()
		}
		Sleep(200)
		this.scanAndCallback()
	}
}

func (this *MCUser) callBackRoutine() {
	logger.Info("[%s] initiate callback goroutine.", this.appAccount)
	if this.conn == nil {
		return
	}
	for {
		//logger.Info("%v size: %v", this.appAccount, this.packetToCallback.Size())
		pktByts := this.packetToCallback.Pop()
		if pktByts != nil {
			packetBytes := pktByts.(*packet.PacketBytes)
			v6Packet := packet.ParseBytesToPacket(packetBytes.HeaderBins, packetBytes.BodyBins, packetBytes.CrcBins, packetBytes.BodyKey, packetBytes.SecKey)
			if v6Packet == nil {
				logger.Error("[rcv]: parse into v6Packet fail.")
				this.conn.Reset()
				continue
			}
			//logger.Info("%v size: %v", this.appAccount, this.packetToCallback.Size())
			this.handleResponse(v6Packet)
		} else {
			Sleep(100)
		}
	}
}

func (this *MCUser) scanAndCallback() {
	if this.msgDelegate == nil {
		logger.Warn("%v need to handle Message for timeout.", this.appAccount)
		return
	}
	this.messageToAck.Lock()
	defer this.messageToAck.Unlock()
	kvs := this.messageToAck.KVs()
	timeoutKeys := list.New()
	for key := range kvs {
		timeoutPacket := kvs[key].(*packet.MIMCTimeoutPacket)
		if CurrentTimeMillis()-timeoutPacket.Timestamp() < cnst.CHECK_TIMEOUT_TIMEVAL_MS {
			continue
		}
		mimcPacket := timeoutPacket.Packet()
		if *(mimcPacket.Type) == MIMC_MSG_TYPE_P2P_MESSAGE {
			p2pMessage := new(MIMCP2PMessage)
			err := Deserialize(mimcPacket.Payload, p2pMessage)
			if !err {
				return
			}
			p2pMsg := msg.NewP2pMsg(mimcPacket.PacketId, p2pMessage.From.AppAccount, p2pMessage.To.AppAccount, mimcPacket.Sequence, mimcPacket.Timestamp, p2pMessage.BizType, p2pMessage.Payload)
			this.msgDelegate.HandleSendMessageTimeout(p2pMsg)
		} else if *(mimcPacket.Type) == MIMC_MSG_TYPE_P2T_MESSAGE {
			p2tMessage := new(MIMCP2TMessage)
			err := Deserialize(mimcPacket.Payload, p2tMessage)
			if !err {
				return
			}
			p2tMsg := msg.NewP2tMsg(mimcPacket.PacketId, p2tMessage.From.AppAccount, mimcPacket.Sequence, mimcPacket.Timestamp, p2tMessage.To.TopicId, p2tMessage.BizType, p2tMessage.Payload)
			this.msgDelegate.HandleSendGroupMessageTimeout(p2tMsg)
		}
		timeoutKeys.PushBack(key)
	}
	for ele := timeoutKeys.Front(); ele != nil; ele = ele.Next() {
		packet := this.messageToAck.Pop(ele.Value.(string))
		if packet == nil {
			logger.Warn("pop message fails. packetId: %v", ele.Value.(string))
		}
	}
}

func (this *MCUser) handleResponse(v6Packet *packet.MIMCV6Packet) {
	if v6Packet.GetHeader() == nil {
		logger.Info("[%v] [handle packet]get a pong packet.", this.appAccount)
		return
	}
	cmd := v6Packet.GetHeader().Cmd
	if cnst.CMD_SECMSG == *cmd {
		this.handleSecMsg(v6Packet)
	} else if cnst.CMD_CONN == *cmd {
		logger.Debug("[handle packet] conn response.")
		connResp := new(XMMsgConnResp)
		err := Deserialize(v6Packet.GetPayload(), connResp)
		if !err {
			logger.Error("[handle packet] parse connResp fail.")
			this.conn.Reset()
			return
		}
		this.conn.HandshakeConnected()
		logger.Debug("[handle packet] handshake succ.")
		this.conn.SetChallenge(*(connResp.Challenge))
		this.conn.SetChallengeAndRc4Key(*(connResp.Challenge))
	} else if cnst.CMD_BIND == *cmd {
		bindResp := new(XMMsgBindResp)
		err := Deserialize(v6Packet.GetPayload(), bindResp)
		if err {
			if *bindResp.Result {
				this.status = Online
				this.lastLoginTimestamp = 0
				logger.Debug("[handle packet] login succ.")
			} else {
				if strings.Compare(bindResp.GetErrorType(), cnst.TOKEN_EXPIRED) == 0 ||
					(strings.Compare(bindResp.GetErrorType(), "auto") == 0 &&
						strings.Compare(bindResp.GetErrorReason(), cnst.TOKEN_INVALID) == 0) {
					this.refreshToken()
				} else {
					this.status = Offline
					logger.Warn("[%v] [handle packet] login fail.", this.appAccount)
				}
			}
			if this.statusDelegate == nil {
				logger.Warn("[%v] status changed, you need to handle this.", this.appAccount)
			} else {
				this.statusDelegate.HandleChange(*(bindResp.Result), bindResp.ErrorType, bindResp.ErrorReason, bindResp.ErrorDesc)
			}
		}
	} else if cnst.CMD_KICK == *cmd {
		this.status = Offline
		kick := "kick"
		logger.Debug("[handle] logout succ.")
		if this.statusDelegate == nil {
			logger.Warn("[%v] status changed, you need to handle this.", this.appAccount)
		} else {
			this.statusDelegate.HandleChange(false, &kick, &kick, &kick)
		}
	} else {
		logger.Debug("cmd: %v", *cmd)
		return
	}
}

func (this *MCUser) handleSecMsg(v6Packet *packet.MIMCV6Packet) {
	if this.msgDelegate == nil {
		logger.Warn("[%v] need to regist mssage handler for received messages.", this.appAccount)
		return
	}
	mimcPacket := new(MIMCPacket)
	err := Deserialize(v6Packet.GetPayload(), mimcPacket)
	if !err {
		logger.Warn("[%v] [handleSecMsg] unserialize mimcPacket fails.%v", this.appAccount, err)
		return
	} else {
		switch *(mimcPacket.Type) {
		case MIMC_MSG_TYPE_PACKET_ACK:
			logger.Debug("handle Sec Msg] packet Ack.")
			packetAck := new(MIMCPacketAck)
			err := Deserialize(mimcPacket.Payload, packetAck)
			if !err {
				return
			}
			this.msgDelegate.HandleServerAck(packetAck.PacketId, packetAck.Sequence, packetAck.Timestamp, packetAck.ErrorMsg)
			packet := this.messageToAck.Pop(*(packetAck.PacketId))
			if packet == nil {
				logger.Warn("[%v] pop message fails. packetId: %v", this.appAccount, *(packetAck.PacketId))
			}
			break
		case MIMC_MSG_TYPE_COMPOUND:
			packetList := new(MIMCPacketList)
			err := Deserialize(mimcPacket.Payload, packetList)
			if !err {
				return
			}
			if this.resource != *(packetList.Resource) {
				logger.Warn("[%v] Handle SecMsg MIMCPacketList resource:, current resource:", this.appAccount, *(packetList.Resource), this.resource)
				return
			}
			seqAckPacket := BuildSequenceAckPacket(this, packetList)
			pktToSend := msg.NewMsgPacket(cnst.MIMC_C2S_SINGLE_DIRECTION, seqAckPacket)
			this.messageToSend.Push(pktToSend)
			pktNum := len(packetList.Packets)
			p2pMsgList := list.New()
			p2tMsgList := list.New()
			for i := 0; i < pktNum; i++ {
				packet := packetList.Packets[i]
				if packet == nil {
					continue
				}
				if *(packet.Type) == MIMC_MSG_TYPE_P2P_MESSAGE {
					p2pMessage := new(MIMCP2PMessage)
					err := Deserialize(packet.Payload, p2pMessage)
					if !err {
						continue
					}
					p2pMsgList.PushBack(msg.NewP2pMsg(packet.PacketId, p2pMessage.From.AppAccount, p2pMessage.To.AppAccount, packet.Sequence, packet.Timestamp, p2pMessage.BizType, p2pMessage.Payload))
					continue
				} else if *(packet.Type) == MIMC_MSG_TYPE_P2T_MESSAGE {
					p2tMessage := new(MIMCP2TMessage)
					err := Deserialize(packet.Payload, p2tMessage)

					if !err {
						continue
					}
					p2tMsgList.PushBack(msg.NewP2tMsg(packet.PacketId, p2tMessage.From.AppAccount, packet.Sequence, packet.Timestamp, p2tMessage.To.TopicId, p2tMessage.BizType, p2tMessage.Payload))
					continue
				}
			}
			if p2pMsgList.Len() > 0 {
				logger.Info("[%v] recv %v p2p msg", this.appAccount, p2pMsgList.Len())
				this.msgDelegate.HandleMessage(p2pMsgList)
			}
			if p2tMsgList.Len() > 0 {
				logger.Info("[%v] recv %v p2t msg", this.appAccount, p2tMsgList.Len())
				this.msgDelegate.HandleGroupMessage(p2tMsgList)
			}
			break
		default:
			break
		}
	}
}
func (this *MCUser) handleToken() {
	this.token = *(this.tokenDelegate.FetchToken())
}

func (this *MCUser) SetResource(resource string) *MCUser {
	this.resource = resource
	return this
}
func (this *MCUser) SetUuid(uuid int64) *MCUser {
	this.uuid = uuid
	return this
}
func (this *MCUser) SetChid(chid int32) *MCUser {
	this.chid = chid
	return this
}
func (this *MCUser) SetConn(conn *MIMCConnection) *MCUser {
	this.conn = conn
	return this
}
func (this *MCUser) SetToken(token string) *MCUser {
	this.token = token
	return this
}
func (this *MCUser) SetSecKey(secKey string) *MCUser {
	this.securityKey = secKey
	return this
}
func (this *MCUser) SetAppPackage(appPackage string) *MCUser {
	this.appPackage = appPackage
	return this
}
func (this *MCUser) SetAppAccount(appAccount string) *MCUser {
	this.appAccount = appAccount
	return this
}
func (this *MCUser) SetAppId(appId int64) *MCUser {
	this.appId = appId
	return this
}

func (this *MCUser) AppAccount() string {
	return this.appAccount
}
func (this *MCUser) AppId() int64 {
	return this.appId
}
func (this *MCUser) Conn() *MIMCConnection {
	return this.conn
}

func (this *MCUser) FeDomain() *string {
	return &(this.feDomain)
}

func (this *MCUser) FeAddress() []string {
	return this.feAddress
}

func (this *MCUser) Uuid() int64 {
	return this.uuid
}
func (this *MCUser) Chid() int32 {
	return this.chid
}
func (this *MCUser) Resource() string {
	return this.resource
}
func (this *MCUser) SecKey() string {
	return this.securityKey
}
func (this *MCUser) Token() *string {
	return &(this.token)
}
func (this *MCUser) ClientAttrs() string {
	return this.clientAttrs
}
func (this *MCUser) CloudAttrs() string {
	return this.cloudAttrs
}
func (this *MCUser) AppPackage() string {
	return this.appPackage
}

func (this *MCUser) Status() UserStatus {
	return this.status
}
