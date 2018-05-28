package cnst

const (
	MIMC_CHID              int32  = 9
	MIMC_COUNTER_VALUE     uint64 = 1
	CONN_BIN_PROTO_VERSION uint32 = 106
	CONN_BIN_PROTO_SDK     uint   = 33

	V6_HEAD_LENGTH        byte = 8
	V6_BODY_HEADER_LENGTH byte = 8
	V6_MAGIC_OFFSET       int  = 0
	V6_VERSION_OFFSET     int  = 2
	V6_BODYLEN_OFFSET     int  = 4
	V6_PAYLOADTYPE_OFFSET int  = 0
	V6_HEADERLEN_OFFSET   int  = 2
	V6_PAYLOADLEN_OFFSET  int  = 4
	V6_CRC_LENGTH         int  = 4

	CMD_CONN   string = "CONN"
	CMD_BIND   string = "BIND"
	CMD_PING   string = "PING"
	CMD_UNBIND string = "UBND"
	CMD_SECMSG string = "SECMSG"
	CMD_KICK   string = "KICK"

	MIMC_SERVER string = "xiaomi.com"

	CIPHER_NONE int32 = 0
	CIPHER_RC4  int32 = 1
	CIPHER_AES  int32 = 2

	CRC_LEN      int    = 4
	RESOURCES    string = "Golang"
	MAGIC        uint16 = 0xc2fe
	V6_VERSION   uint16 = 0x0005
	HEADER_TYPE  int    = 3
	PAYLOAD_TYPE uint16 = 2

	MIMC_METHOD string = "XIAOMI-PASS"
	NO_KICK     string = "0"

	PING_TIMEVAL_MS                 int64 = 3000 //15s
	CONNECT_TIMEOUT                 int64 = 5000
	LOGIN_TIMEOUT                   int64 = 5000
	CHECK_TIMEOUT_TIMEVAL_MS        int64 = 10000
	RESET_SOCKET_TIMEOUT_TIMEVAL_MS int64 = 5000

	MIMC_TOKEN_EXPIRE string = "token-expired"

	MIMC_C2S_DOUBLE_DIRECTION string = "C2S_DOUBLE_DIRECTION"
	MIMC_C2S_SINGLE_DIRECTION string = "C2S_SINGLE_DIRECTION"

	FE_IP_ONLINE   string = "app.chat.xiaomi.net"
	FE_PORT_ONLINE int    = 80

	FE_IP_STAGING   string = "10.38.162.117"
	FE_PORT_STAGING int    = 5222

	TOKEN_IP_ONLINE  string = "mimc.chat.xiaomi.net/api/account/token"
	TOKEN_IP_STAGING string = "10.38.162.149"
)
