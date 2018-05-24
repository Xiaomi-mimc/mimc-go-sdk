package mimc

import (
	"bytes"
	"fmt"
	"github.com/Xiaomi-mimc/mimc-go-sdk/cipher"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler"
	"github.com/Xiaomi-mimc/mimc-go-sdk/packet"
	. "github.com/Xiaomi-mimc/mimc-go-sdk/protobuf/mimc"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/byte"
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
)

// online
var httpUrl string = "https://mimc.chat.xiaomi.net/api/account/token"
var appId int64 = int64(2882303761517613988)
var appKey string = "5361761377988"
var appSecurt string = "2SZbrJOAL1xHRKb7L9AiRQ=="
var appAccount1 string = "Alice"
var appAccount2 string = "Bob"

/*// staging
var httpUrl string = "http://10.38.162.149/api/account/token"
var appId int64 = int64(2882303761517479657)
var appKey string = "5221747911657"
var appSecurt string = "PtfBeZyC+H8SIM/UXhZx1w=="
var appAccount1 string = "Alice"
var acc1UUID = int64(10776577642332160)
var appAccount2 string = "Bob"
var acc2UUID = int64(10778725662851072)*/

func createHandlers(appAccount string) (*handler.StatusHandler, *handler.TokenHandler, *handler.MsgHandler) {
	return handler.NewStatusHandler(), handler.NewTokenHandler(&httpUrl, &appKey, &appSecurt, &appAccount, &appId), handler.NewMsgHandler()
}

func TestBuildBindPacket(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser := NewUser(appAccount1)
	mcUser.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()

}
func TestLoginAndOut(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser := NewUser(appAccount1)
	mcUser.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	mcUser.Login()
	Sleep(3000)
	mcUser.Logout()
	Sleep(5000)
}
func TestSendMessage(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser1 := NewUser(appAccount1)
	mcUser1.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	mcUser1.Login()
	Sleep(2000)
	for i := 0; i < 5; i++ {
		str := strconv.FormatInt(int64(i), 10)
		mcUser1.SendMessage(appAccount2, []byte("hello world!"+str))
		Sleep(100)
	}
	Sleep(30000)
	/*for i := 0; i < 5; i++ {
		str := strconv.FormatInt(int64(i), 10)
		mcUser1.SendMessage(appAccount2, []byte("hello world!"+str))
		Sleep(100)
	}*/
	//Sleep(3000)
}

func TestSendP2TMessage(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser1 := NewUser(appAccount1)
	mcUser1.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	topicId := int64(10871081150185472)
	mcUser1.Login()
	Sleep(2000)
	mcUser1.SendGroupMessage(&topicId, []byte("hello everybody!"))
	Sleep(30000)
}

func TestRecvMessage(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount2)
	mcUser2 := NewUser(appAccount2)
	mcUser2.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	mcUser2.Login()
	Sleep(90000)
}

func TestSerialAndUnSerial(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser1 := NewUser(appAccount1)
	mcUser1.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	mcUser1.Login()
	Sleep(3000)
	v6Packet, _ := BuildP2PMessagePacket(mcUser1, appAccount2, []byte("hello world!"), true)

	mimcPacketPayload := v6Packet.GetPayload()

	fmt.Printf("[seria] mimcpacket: %v\n", mimcPacketPayload)

	mimcPacket := new(MIMCPacket)

	err := proto.Unmarshal(mimcPacketPayload, mimcPacket)
	if err != nil {
		fmt.Printf("unserialize mimcpacket error!\n")
	}
	p2pMsgPayload := mimcPacket.Payload

	p2pMsg := new(MIMCP2PMessage)

	err = proto.Unmarshal(p2pMsgPayload, p2pMsg)
	if err != nil {
		fmt.Printf("unserialize p2pmsg error!\n")
	}
	Sleep(5000)
	fmt.Printf("p2p: %v\nmimcPacket: %v\n", *mimcPacket, *p2pMsg)

	buffer := new(bytes.Buffer)
	buffer.Write(p2pMsg.Payload)
	fmt.Printf("msg: %v\n", buffer.String())
}

func TestEncryptAndUnEncrypt(t *testing.T) {
	statusHandler, tokenHandler, msgHandler := createHandlers(appAccount1)
	mcUser1 := NewUser(appAccount1)
	mcUser1.RegisterStatusDelegate(statusHandler).RegisterTokenDelegate(tokenHandler).RegisterMessageDelegate(msgHandler).InitAndSetup()
	mcUser1.Login()
	Sleep(3000)
}

func TestUnmarshalMIMCPacket(t *testing.T) {
	bytes := []byte{10, 12, 48, 57, 100, 49, 49, 55, 48, 54, 97, 55, 48, 97, 18, 24, 99, 111, 109, 46, 120, 105, 97, 111, 109, 105, 46, 105, 109, 99, 46, 116, 101, 115, 116, 95, 97, 112, 112, 49, 24, 129, 212, 248, 138, 250, 180, 142, 27, 32, 4, 42, 97, 10, 17, 69, 115, 114, 87, 122, 119, 90, 98, 119, 68, 100, 74, 101, 114}
	mimcPacket := new(MIMCPacket)
	proto.Unmarshal(bytes, mimcPacket)

	packetAck := new(MIMCPacketAck)

	err := proto.Unmarshal(mimcPacket.Payload, packetAck)
	if err != nil {
		return
	}
	fmt.Printf("%v\n", *(mimcPacket.Type))
	fmt.Printf("%v\n", *packetAck)
}
func TestUnmarshalMIMCPacketWithEcnrypt(t *testing.T) {
	packets := []byte{194, 254, 0, 5, 0, 0, 1, 51, 213, 21, 24, 153, 188, 73, 243, 90, 60, 167, 223, 190, 39, 185, 52, 189, 184, 73, 34, 85, 140, 173, 58, 72, 171, 189, 143, 4, 171, 146, 61, 108, 63, 71, 25, 93, 129, 92, 8, 133, 153, 216, 160, 150, 238, 156, 51, 232, 226, 205, 222, 175, 244, 133, 252, 244, 83, 38, 221, 236, 18, 56, 80, 187, 208, 72, 226, 178, 7, 178, 55, 148, 88, 151, 250, 137, 199, 140, 44, 18, 130, 154, 220, 120, 97, 130, 14, 42, 95, 206, 20, 176, 94, 102, 138, 241, 147, 110, 103, 114, 137, 98, 164, 153, 151, 69, 174, 193, 104, 219, 191, 127, 130, 196, 148, 4, 224, 16, 227, 59, 45, 123, 221, 124, 181, 110, 156, 180, 96, 93, 153, 72, 123, 120, 246, 195, 176, 194, 7, 121, 246, 68, 167, 212, 252, 92, 156, 92, 231, 135, 66, 39, 92, 116, 181, 112, 152, 211, 149, 76, 201, 87, 170, 218, 232, 245, 205, 48, 27, 176, 54, 141, 231, 147, 227, 6, 218, 68, 161, 194, 38, 104, 253, 7, 65, 194, 111, 135, 152, 143, 165, 181, 106, 171, 69, 91, 201, 92, 226, 103, 20, 139, 124, 215, 202, 74, 140, 164, 9, 83, 164, 141, 240, 152, 231, 186, 249, 201, 189, 220, 8, 243, 121, 242, 177, 165, 40, 40, 70, 48, 74, 95, 61, 109, 0, 187, 126, 237, 243, 11, 178, 158, 53, 37, 19, 218, 20, 55, 73, 44, 240, 153, 75, 148, 123, 144, 180, 117, 82, 80, 40, 145, 21, 182, 45, 218, 235, 5, 173, 55, 205, 45, 21, 150, 135, 218, 101, 89, 25, 208, 153, 11, 12, 174, 176, 39, 213, 205, 202, 129, 40, 237, 223, 30, 208, 12, 82, 217, 234, 222, 34, 107, 116, 126, 19, 15, 46, 196, 86, 159, 129}
	rc4Key := []byte{11, 182, 82, 247, 122}
	secKey := "XOF8dTB2dOVVkEdZB9O21w=="

	headerBins := byteutil.Copy(&packets, 0, int(cnst.V6_HEAD_LENGTH))
	magic := byteutil.GetUint16FromBytes(&headerBins, cnst.V6_MAGIC_OFFSET)
	fmt.Printf("magic: %v\n", magic)
	version := byteutil.GetUint16FromBytes(&headerBins, cnst.V6_VERSION_OFFSET)
	fmt.Printf("version: %v\n", version)
	bodyLen := byteutil.GetIntFromBytes(&headerBins, cnst.V6_BODYLEN_OFFSET)
	fmt.Printf("packet len: %v\n", bodyLen)
	bodyBins := byteutil.Copy(&packets, int(cnst.V6_HEAD_LENGTH), bodyLen)
	crcBins := byteutil.Copy(&packets, int(cnst.V6_HEAD_LENGTH)+bodyLen, cnst.V6_CRC_LENGTH)

	v6Packdet := packet.ParseBytesToPacket(&headerBins, &bodyBins, &crcBins, rc4Key, secKey)

	if v6Packdet == nil {
		fmt.Printf("v6packet is nil.\n")
		return
	}

	mimcPacket := new(MIMCPacket)

	err := proto.Unmarshal(v6Packdet.GetPayload(), mimcPacket)
	if err != nil {
		fmt.Printf("parse MIMCPacket fail. err: %v\n", err)
		return
	} else {
		fmt.Printf("MIMCPacket: %v\n", mimcPacket)
	}

}

func TestPayloadKey(t *testing.T) {
	key := "XOF8dTB2dOVVkEdZB9O21w=="
	val := "3"
	value := []byte(val)
	enVal := PayloadKey(key, value)
	enVal1 := cipher.GenerateKeyForRC4(key, val)

	fmt.Printf("key: %v\nvalue: %v\nenVal: %v\nenVal: %v\n", []byte(key), value, enVal, enVal1)
}
