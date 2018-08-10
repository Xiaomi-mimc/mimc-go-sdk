package main

import (
	"github.com/Xiaomi-mimc/mimc-go-sdk"
	"github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler"
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
var appAccount1 string = "t11"
var acc1UUID = int64(10776577642332160)
var appAccount2 string = "t12"
var acc2UUID = int64(10778725662851072)

func main() {
	// 创建用户
	leijun := createUser(&appAccount1)
	mifen := createUser(&appAccount2)

	// 用户登录
	leijun.Login()
	mifen.Login()
	mimc.Sleep(3000)

	// 互发消息
	leijun.SendMessage(appAccount2, []byte("1"))
	mimc.Sleep(20000)
	mifen.SendMessage(appAccount1, []byte("2"))
	//leijun.SendMessage(appAccount2, []byte("3"))
	//mifen.SendMessage(appAccount1, []byte("4"))
	mimc.Sleep(30000)

	// 用户退出
	leijun.Logout()
	mifen.Logout()

	mimc.Sleep(10000)
}

// 创建用户
func createUser(appAccount *string) *mimc.MCUser {

	mcUser := mimc.NewUser(*appAccount)
	statusDelegate, tokenDelegate, msgDelegate := createDelegates(appAccount)
	mcUser.RegisterStatusDelegate(statusDelegate).RegisterTokenDelegate(tokenDelegate).RegisterMessageDelegate(msgDelegate).InitAndSetup()
	return mcUser
}

// 用户自定义消息、用户状态、Token的处理器
func createDelegates(appAccount *string) (*handler.StatusHandler, *handler.TokenHandler, *handler.MsgHandler) {
	return handler.NewStatusHandler(), handler.NewTokenHandler(&httpUrl, &appKey, &appSecurt, appAccount, &appId), handler.NewMsgHandler()
}
