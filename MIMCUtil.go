package mimc

import (
	"bytes"
	"container/list"
	"encoding/base64"
	"fmt"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"github.com/Xiaomi-mimc/mimc-go-sdk/packet"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/ims"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/mimc"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/id"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/string"
	"github.com/golang/protobuf/proto"
	"sort"
	"strconv"
	"strings"
	"time"
)

func CurrentTimeMillis() int64 {
	now := time.Now()
	return now.UnixNano() / 1e6
}

func Sleep(millis int64) {
	time.Sleep(time.Duration(millis) * time.Millisecond)
}

func Deserialize(data []byte, pb proto.Message) bool {
	err := proto.Unmarshal(data, pb)
	if err != nil {
		fmt.Printf("deserialize error: %v\n", err)
		return false
	}
	return true
}

// 并不是所有的包的header都一样的
func createClientHeader(mcUser *MCUser, cmd string, msgId *string, cipher int32) *ClientHeader {
	if mcUser == nil || len(cmd) == 0 {
		return nil
	}
	header := new(ClientHeader)
	header.Id = msgId
	chid := cnst.MIMC_CHID
	header.Chid = &chid
	uuid := mcUser.Uuid()
	header.Uuid = &uuid
	resource := mcUser.Resource()
	header.Resource = &(resource)
	header.Cmd = &cmd
	header.Cipher = &cipher
	server := cnst.MIMC_SERVER
	header.Server = &server
	csreq := ClientHeader_CS_REQ
	header.DirFlag = &csreq
	return header
}

func createXMMsgBind(mcUser *MCUser, header *ClientHeader) *XMMsgBind {
	bind := new(XMMsgBind)
	bind.Token = mcUser.Token()
	method := cnst.MIMC_METHOD
	bind.Method = &method
	clientAttrs := mcUser.ClientAttrs()
	bind.ClientAttrs = &clientAttrs
	cloudAttrs := mcUser.CloudAttrs()
	bind.CloudAttrs = &cloudAttrs
	nokick := cnst.NO_KICK
	bind.Kick = &nokick
	bind.Sig = generateSing(header, bind, mcUser.Conn().Challenge(), mcUser.SecKey())
	return bind
}

func BuildBindPacket(mcUser *MCUser) *packet.MIMCV6Packet {
	clientHeader := createClientHeader(mcUser, cnst.CMD_BIND, id.Generate(), cnst.CIPHER_NONE)
	xmmsgBind := createXMMsgBind(mcUser, clientHeader)
	v6Packet := packet.NewV6Packet()
	v6Packet.ClientHeader(clientHeader)
	payload, _ := proto.Marshal(xmmsgBind)
	v6Packet.Payload(payload)
	v6Packet.PayloadType(cnst.PAYLOAD_TYPE)
	return v6Packet
}

func BuildUnBindPacket(mcUser *MCUser) *packet.MIMCV6Packet {
	clientHeader := createClientHeader(mcUser, cnst.CMD_UNBIND, id.Generate(), cnst.CIPHER_NONE)
	v6Packet := packet.NewV6Packet()
	v6Packet.ClientHeader(clientHeader)
	v6Packet.Payload(nil)
	return v6Packet
}
func BuildConnectionPacket(udid string, mcUser *MCUser) *packet.MIMCV6Packet {
	clientHeader := createClientHeader(mcUser, cnst.CMD_CONN, id.Generate(), cnst.CIPHER_NONE)
	v6Packet := packet.NewV6Packet()
	xMMsgConn := new(XMMsgConn)
	os := "macOs"
	xMMsgConn.Os = &os
	xMMsgConn.Udid = &udid
	version := cnst.CONN_BIN_PROTO_VERSION
	xMMsgConn.Version = &version
	v6Packet.PayloadType(cnst.PAYLOAD_TYPE)
	v6Packet.ClientHeader(clientHeader)
	payload, _ := proto.Marshal(xMMsgConn)
	v6Packet.Payload(payload)
	return v6Packet
}

func BuildSequenceAckPacket(mcUser *MCUser, packetList *MIMCPacketList) *packet.MIMCV6Packet {
	clientHeader := createClientHeader(mcUser, cnst.CMD_SECMSG, id.Generate(), cnst.CIPHER_RC4)

	mimcPacket := new(MIMCPacket)
	mimcPacket.PacketId = id.Generate()
	pkg := mcUser.AppPackage()
	mimcPacket.Package = &pkg
	mimcPacket.Sequence = packetList.MaxSequence
	msgType := MIMC_MSG_TYPE_SEQUENCE_ACK
	mimcPacket.Type = &msgType

	seqAck := new(MIMCSequenceAck)
	seqAck.Uuid = packetList.Uuid
	seqAck.Resource = packetList.Resource
	seqAck.Sequence = packetList.MaxSequence

	seqAckBin, _ := proto.Marshal(seqAck)

	mimcPacket.Payload = seqAckBin

	mimcBins, _ := proto.Marshal(mimcPacket)

	v6Packet := packet.NewV6Packet()
	v6Packet.PayloadType(cnst.PAYLOAD_TYPE)
	v6Packet.ClientHeader(clientHeader)
	v6Packet.Payload(mimcBins)

	return v6Packet
}
func BuildP2TMessagePacket(mcUser *MCUser, appTopic int64, msg []byte, isStore bool) (*packet.MIMCV6Packet, *MIMCPacket) {
	clientHeader := createClientHeader(mcUser, cnst.CMD_SECMSG, id.Generate(), cnst.CIPHER_RC4)

	fromUser := buildMIMCUser(mcUser)
	toGroup := buildMIMCGroup(mcUser.AppId(), appTopic)

	p2tMsg := new(MIMCP2TMessage)
	p2tMsg.From = fromUser
	p2tMsg.To = toGroup
	p2tMsg.Payload = msg
	p2tMsg.IsStore = &isStore

	mimcPacket := new(MIMCPacket)
	mimcPacket.PacketId = clientHeader.Id
	pkg := mcUser.AppPackage()
	mimcPacket.Package = &pkg
	msgType := MIMC_MSG_TYPE_P2T_MESSAGE
	mimcPacket.Type = &msgType
	payload, err := proto.Marshal(p2tMsg)
	if err != nil {
		fmt.Printf("serialize P2P msg fail.\n")
	}
	mimcPacket.Payload = payload

	v6Packet := packet.NewV6Packet()
	v6Packet.PayloadType(cnst.PAYLOAD_TYPE)
	v6Packet.ClientHeader(clientHeader)
	payload, err = proto.Marshal(mimcPacket)
	if err != nil {
		fmt.Printf("serialize MIMCPackdet fail.\n")
	}
	v6Packet.Payload(payload)
	return v6Packet, mimcPacket
}
func BuildP2PMessagePacket(mcUser *MCUser, appAccount string, msg []byte, isStore bool) (*packet.MIMCV6Packet, *MIMCPacket) {
	clientHeader := createClientHeader(mcUser, cnst.CMD_SECMSG, id.Generate(), cnst.CIPHER_RC4)

	fromUser := buildMIMCUser(mcUser)
	toUser := buildMIMCUser(NewMCUser().SetAppAccount(appAccount).SetAppId(mcUser.AppId()))

	p2pMsg := new(MIMCP2PMessage)
	p2pMsg.From = fromUser
	p2pMsg.To = toUser
	p2pMsg.Payload = msg
	p2pMsg.IsStore = &isStore

	mimcPacket := new(MIMCPacket)
	mimcPacket.PacketId = clientHeader.Id
	pkg := mcUser.AppPackage()
	mimcPacket.Package = &pkg
	msgType := MIMC_MSG_TYPE_P2P_MESSAGE
	mimcPacket.Type = &msgType
	payload, err := proto.Marshal(p2pMsg)
	if err != nil {
		fmt.Printf("serialize P2P msg fail.\n")
	}
	mimcPacket.Payload = payload

	v6Packet := packet.NewV6Packet()
	v6Packet.PayloadType(cnst.PAYLOAD_TYPE)
	v6Packet.ClientHeader(clientHeader)
	payload, err = proto.Marshal(mimcPacket)
	if err != nil {
		fmt.Printf("serialize MIMCPackdet fail.\n")
	}
	v6Packet.Payload(payload)
	return v6Packet, mimcPacket
}

func buildMIMCUser(mcUser *MCUser) *MIMCUser {
	mimcUser := new(MIMCUser)
	appId := mcUser.AppId()
	mimcUser.AppId = &appId
	account := mcUser.AppAccount()
	mimcUser.AppAccount = &account
	uuid := mcUser.Uuid()
	mimcUser.Uuid = &uuid
	resource := mcUser.Resource()
	mimcUser.Resource = &resource
	return mimcUser
}

func buildMIMCGroup(appId, topicId int64) *MIMCGroup {
	mimcGroup := new(MIMCGroup)
	mimcGroup.AppId = &appId
	mimcGroup.TopicId = &topicId
	return mimcGroup
}

func BuildPingPacket(mcUser *MCUser) *packet.MIMCV6Packet {
	v6Packet := packet.NewV6Packet()
	return v6Packet
}

func PayloadKey(key string, value []byte) []byte {
	if len(key) == 0 || len(value) == 0 {
		return nil
	}
	keyBytes, _ := base64.StdEncoding.DecodeString(key)
	buffer := new(bytes.Buffer)
	buffer.Write(keyBytes)
	buffer.WriteByte('_')
	buffer.Write(value)
	return buffer.Bytes()
}

func generateSing(header *ClientHeader, xmmsgBind *XMMsgBind, challenge string, secKey string) *string {
	params := make(map[string]string)
	params["challenge"] = challenge
	params["token"] = *(xmmsgBind.Token)
	params["chid"] = strconv.FormatInt(int64(*(header.Chid)), 10)
	server := "@xiaomi.com/"
	uuid := strconv.FormatInt(*(header.Uuid), 10)
	params["from"] = strutil.ConcatStrs(&uuid, &server, header.Resource)

	params["id"] = *(header.Id)
	params["to"] = *(header.Server)
	params["kick"] = *(xmmsgBind.Kick)
	params["client_attrs"] = *(xmmsgBind.ClientAttrs)
	params["cloud_attrs"] = *(xmmsgBind.CloudAttrs)
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	exps := list.New()
	exps.PushBack(strings.ToUpper(*(xmmsgBind.Method)))

	for _, key := range keys {
		equal := "="
		val := params[key]
		exps.PushBack(strutil.ConcatStrs(&key, &equal, &val))
	}
	exps.PushBack(secKey)
	and := "&"
	expsStr := strutil.ConcatStrsByStr(exps, &and)
	return strutil.Sha1(&expsStr)

}
