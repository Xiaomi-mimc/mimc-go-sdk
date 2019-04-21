package main

import (
	"fmt"
	"github.com/Xiaomi-mimc/mimc-go-sdk"
	"github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler"
	"github.com/Xiaomi-mimc/mimc-go-sdk/util/log"
)

/**
 * @Important:
 *   以下appId/appKey/appSecurity是小米MIMCDemo APP所有，会不定期更新
 *   所以，开发者应该将以下三个值替换为开发者拥有APP的appId/appKey/appSecurity
 * @Important:
 *   开发者访问小米开放平台(https://dev.mi.com/console/man/)，申请appId/appKey/appSecurity
 **/

var httpUrl string = "https://mimc.chat.xiaomi.net/api/account/token"
var appId int64 = int64(2882303761517613988)
var appKey string = "5361761377988"
var appSecurt string = "2SZbrJOAL1xHRKb7L9AiRQ=="
var appAccount1 string = "happycsy"
var acc1UUID = int64(10776577642332160)
var appAccount2 string = "happyysm"
var acc2UUID = int64(10778725662851072)

func init() {
	fmt.Println("call init method of Demo")
	log.SetLogLevel(log.DebugLevel)
	log.SetLogPath("./mimc_demo.log")
	//logger := log.GetLogger()
}

func main() {

	// 创建用户
	leijun := createUser(appAccount1)
	mifen := createUser(appAccount2)

	// 用户登录
	leijun.Login()
	mifen.Login()
	mimc.Sleep(3000)
	leijun.SendMessage(appAccount2, []byte("123"))
	leijun.SendMessage(appAccount2, []byte("789"))
	//mifen.SendMessage(appAccount1, []byte("456"))
	/*counter := 0
	now := time.Now().Unix()
	var interval int64 = 60 * 60 * 1 // 1h
	var interval int64 = 60 // 1 minute
	for nowing := now; nowing-now < interval; {
		// 互发消息
		leijun.SendMessage(appAccount2, []byte("123"))
		mifen.SendMessage(appAccount1, []byte("456"))
		nowing = time.Now().Unix()
		counter = counter + 1
	}
	log.GetLogger().Info("send %d times.", counter)*/
	mimc.Sleep(5000)

	// 用户退出
	leijun.Logout()
	mifen.Logout()

	mimc.Sleep(1000)

}

// 创建用户
func createUser(appAccount string) *mimc.MCUser {

	mcUser := mimc.NewUser(appId, appAccount)
	statusDelegate, tokenDelegate, msgDelegate := createDelegates(appAccount)
	mcUser.RegisterStatusDelegate(statusDelegate).RegisterTokenDelegate(tokenDelegate).RegisterMessageDelegate(msgDelegate).InitAndSetup()
	return mcUser
}

// 用户自定义消息、用户状态、Token的处理器
func createDelegates(appAccount string) (*handler.StatusHandler, *handler.TokenHandler, *handler.MsgHandler) {
	return handler.NewStatusHandler(appAccount), handler.NewTokenHandler(&httpUrl, &appKey, &appSecurt, &appAccount, &appId), handler.NewMsgHandler(appAccount)
}
