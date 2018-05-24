## 开发者文档
* [如何引入](#如何引入)
> * [下载源码](#下载源码)
> * [下载依赖库](#下载依赖库)
> * [引入包](#引入包)
* [用户创建及初始化](#用户创建及初始化)
> * [创建用户](#创建用户)
> * [安全认证](#安全认证)
> * [用户在线状态回调](#用户在线状态回调)
> * [初始化](#初始化)
* [登录](#登录)
* [发送单聊消息](#发送单聊消息)
* [发送群聊消息](#发送群聊消息)
* [接收消息](#接收消息)
> * [单聊消息回调](#单聊消息回调)
> * [群聊消息回调](#群聊消息回调)
> * [单聊消息超时回调](#单聊消息超时回调)
> * [群聊消息超时回调](#群聊消息超时回调)
> * [消息送达服务器响应回调](#消息送达服务器响应回调)
* [退出](#退出)
* [示例代码](#示例代码)

## 如何引入
MIMC-Go-SDK的源码已经上传到github中，开发者可以下载至本地GOPATH下安装。
### 下载源码
```sbtshell
    go get github.com/Xiaomi-mimc/mimc-go-sdk
    cd $GOPATH/github.com/Xiaomi-mimc/mimc-go-sdk
    go build
    go install
```
> 注：go get会从github上fetch所有的分支/tag，默认分支为master分支。

## 下载依赖库
MIMC-Go-SDK依赖proto buffer进行序列化与反序列化数据，在使用时，确保已经install了proto buffer。如未安装，可参考如下进行安装操作。
```sbtshell
    go get github.com/golang/protobuf/proto
    // 进入下载目录
    cd $GOPATH/src/github.com/golang/protobuf/proto 
    // 编译安装
    go build
    go install
 ```
 ## 引入包
 在使用MIMC-Go-SDK时，仅需引入一个包即可。
 ```golang
     // 引入包路径, 该路径下引入了mimc包
     import "github.com/Xiaomi-mimc/mimc-go-sdk"
 ```

## 用户创建及初始化
用户登录前需要：创建用户->安全认证->用户在线状态回调->接收消息回调->初始化。
### 创建用户
```golang
    import "github.com/XiaoMi-mimc/mimc-go-sdk"
    var appAccount string = "leijun"
    mcUser := mimc.NewUser(appAccount)
```
### 安全认证
参考[安全认证](03-auth.md)文档实现：

```golang
    /**
     * 同步访问代理认证服务(appProxyService)，
     * 从代理认证服务返回结果中解析[小米认证服务下发的原始数据]并返回
     **/
    type Token interface {
        FetchToken() *string
    }
```
实现该接口后，通过MCUser.registerTokenDelegate(tokenDeleget Token)方法实现**token回调**的注册。
Token接口实现示例如下：
```golang
    type TokenHandler struct {
        httpUrl    string                      
        AppId      int64  `json:"appId"`
        AppKey     string `json:"appKey"`
        AppSecret  string `json:"appSecret"`
        AppAccount string `json:"appAccount"`
    }
    func (this *TokenHandler) FetchToken() *string {
        jsonBytes, err := json.Marshal(*this)
        if err != nil {
            return nil
        }
        requestJsonBody := bytes.NewBuffer(jsonBytes).String()
        request, err := http.Post(this.httpUrl, "application/json", strings.NewReader(requestJsonBody))
        if err != nil {
            return nil
        }
        defer request.Body.Close()
        body, err := ioutil.ReadAll(request.Body)
        if err != nil {
            return nil
        }
        token := string(body)
        return &token
    }
```
> 注：示例中的实现是**FetchToken--请求-->MIMC服务**，推荐开发者实现方式：**FetchToken--请求-->appProxyService--请求-->MIMC服务**。
#### 用户在线状态回调
用户状态发生改变时，SDK通过**用户在线状态回调**来通知开发者进行处理，接口如下：
```golang
    type StatusDelegate interface {
        /**
        * @param[isOnline bool] true: 在线，false：离线
        * @param[errType *string] 登录失败类型
        * @param[errReason *string] 登录失败原因
        * @param[errDescription *string] 登录失败原因描述
        */
        HandleChange(isOnline bool, errType, errReason, errDescription *string)
    }
```

开发者在该接口中实现用户状态变化的业务逻辑，通过MCUser.registerStatusDelegate(statusDelegate StatusDelegate)方法实现**用户在线状态回调**的注册。

**StatusDelegate**接口实现示例如下：
```golang
    type StatusHandler struct {
    }
    func (this StatusHandler) HandleChange(isOnline bool, errType, errReason, errDescription *string) {
        if isOnline {
            logger.Info("status changed: online.")
            // to do something for online.
        } else {
            logger.Info("status changed: offline.")
            // to do something for offline.
        }
    }
```

### 初始化
MCUser初始化以及启动读、写协程。
```golang
    mcUser.InitAndSetup()
```
至此，用户就可以进行登录操作了。

## 登录
```golang
    mcUser.Login()
```
## 发送单聊消息
发送单聊消息的接口如下：
```golang
    /**
     * @param[toAppAccount string] 接收者
     * @param[msgByte []byte] 开发者自定义消息体
     * @return 客户端生成的数据包Id
     **/
    func (this *MCUser) SendMessage(toAppAccount string, msgByte []byte) string
```
发送单聊消息示例如下：
```golang
    packetId := mcUser.SendMessage("MiFen", []byte("Are you OK?"))
```
## 发送群聊消息
发送群聊消息的接口如下：
```golang
    /**
     * @param[topicId *int64] 群Id
     * @param[msgByte []byte] 开发者自定义消息体
     * @return 客户端生成的数据包Id
     **/
    func (this *MCUser) SendGroupMessage(topicId *int64, msgByte []byte) string
```
发送群聊消息示例如下：
```golang
    groupId := int64(123456789) 
    packetId := mcUser.SendMessage(&groupid, []byte("Are you OK?"))
```
## 接收消息
### 接口概览
当MIMC-Go-SDK收到消息时，会通过**消息回调**来执行开发者的业务逻辑，主要包括5种回调：
> 1) 单聊消息回调
> 2) 群聊消息回调
> 3) 单聊消息超时回调
> 4) 群聊消息超时回调
> 5) 消息送达服务器响应回调

消息回调的接口如下：
```golang
type MessageHandlerDelegate interface {
	HandleMessage(packets *list.List)
	HandleGroupMessage(packets *list.List)
	HandleServerAck(packetId *string, sequence, timestamp *int64)
	HandleSendMessageTimeout(message *msg.P2PMessage)
	HandleSendGroupMessageTimeout(message *msg.P2TMessage)
}
```
在处理“接收单聊/群聊消息”、“单聊/群聊消息超时”以及“消息送达服务器”的业务逻辑时，开发者通过实现MessageHandlerDelegate接口的方法，并通过MCUser.RegisterMessageDelegate(msgDelegate MessageHandlerDelegate)方法实现**接收消息回调**的注册。
### 单聊消息回调
单聊消息回调的HandleMessage方法实现示例如下：
```golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleMessage(packets *list.List) {
        for ele := packets.Front(); ele != nil; ele = ele.Next() {
            p2pmsg := ele.Value.(*msg.P2PMessage)
            logger.Info("[handle p2p msg]%v -> %v, pcktId: %v, timestamp: %v.", *(p2pmsg.FromAccount()), string(p2pmsg.Payload()), *(p2pmsg.PacketId()), *(p2pmsg.Timestamp()))
        }
    }
```
> 注：HandleMessage方法中参数packets是P2PMessage的集合。

P2PMessage的结构如下：
```golang
type P2PMessage struct {
    packetId     *string 数据包Id
    sequence     *int64  消息序列号
    timestamp    *int64  时间戳
    fromAccount  *string 发送者账号
    fromResource *string 发送者设备标记
    payload      []byte  消息体
}
```
### 群聊消息回调
群聊消息回调的HandleGroupMessage方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleGroupMessage(packets *list.List) {
        for ele := packets.Front(); ele != nil; ele = ele.Next() {
            p2tmsg := ele.Value.(*msg.P2TMessage)
            logger.Info("[handle p2t msg]%v -> %v, pcktId: %v, timestamp: %v.", *(p2tmsg.FromAccount()), string(p2tmsg.Payload()), *(p2tmsg.PacketId()), *(p2tmsg.Timestamp()))
        }
    }
```
> 注：HandleGroupMessage方法中参数packets是P2TMessage的集合。

P2TMessage的结构如下：
```golang
type P2TMessage struct {
    packetId     *string 数据包Id
    sequence     *int64  消息序列号
    timestamp    *int64  时间戳
    fromAccount  *string 发送者账号
    fromResource *string 发送者设备
    groupId      *int64  群Id
    payload      []byte  消息体
}
```
### 单聊消息超时回调
单聊消息超时回调的HandleSendMessageTimeout方法的实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleSendMessageTimeout(message *msg.P2PMessage) {
        logger.Info("[handle p2pmsg timeout] packetId:%v, msg:%v, time: %v.", *(message.PacketId()), string(message.Payload()), time.Now())
    }
```
### 群聊消息超时回调
群聊消息超时回调的HandleSendGroupMessageTimeout方法的实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleSendGroupMessageTimeout(message *msg.P2TMessage) {
        logger.Info("[handle p2tmsg timeout] packetId:%v, msg:%v.", *(message.PacketId()), string(message.Payload()))
    }
```
### 消息送达服务器响应回调
消息送达服务器响应回调HandleServerAck方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleServerAck(packetId *string, sequence, timestamp *int64) {
        logger.Info("[handle server ack] packetId:%v, timestamp:%v.", *packetId, *timestamp)
    }
```
> 注：超时回调，保证消息比丢失；消息送达服务器时，服务会响应Ack，SDK收到Ack响应才说明消息送达。

## 退出
``` golang
    mcUser.Logout()
```

## 示例代码
在MIMC-Go-SDK中提供了一个demo，实现两个用户收发消息，该Demo实现可分为以下几步：
> 1) 实现Token接口

> 2) 实现StatusDelegate接口

> 3) 实现MessageHandlerDelegate接口

> 4) 创建回调接口实现类以及用户

> 5) 登录、收发消息、退出

demo在MIMC-Go-SDK中的目录如下：

    |-$GOPATH
        |-bin
        |-pkg
        |-src
            |-github.com
                |-Xiaomi-mimc
                    |-mimc-go-sdk
                        |-demo
                            |-handler
                                |-TokenHandler.go
                                |-StatusHandler.go
                                |-MsgHandler.go
                            |-MCDemo.go


### 实现Token接口
TokenHandler.go中定义了实现Token接口的struct：**TokenHandler**。该文件在$GOPATH/src/github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler目录。
### 实现StatusDelegate接口
StatusHandler.go中定义了实现StatusDelegate接口的struct：**StatusHandler**。该文件在$GOPATH/src/github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler
### 实现MessageHandlerDelegate接口
MsgHandler.go中定义了实现MessageDelegate接口的struct：**MsgHandler**。该文件在$GOPATH/src/github.com/Xiaomi-mimc/mimc-go-sdk/demo/handler目录
### 创建回调接口实现类以及用户
```golang
// 创建三个回调接口的实现类
func createDelegates(appAccount *string) (*handler.StatusHandler, *handler.TokenHandler, *handler.MsgHandler) {
    return handler.NewStatusHandler(), handler.NewTokenHandler(&httpUrl, &apppKey, &appSecurt, appAccount, &appId), handler.NewMsgHandler()
}
// 创建用户
func createUser(appAccount *string) *user.MCUser {
    mcUser := mimc.NewUser(*appAccount)
    // 创建Token接口，StatusDelegate接口，MessageDelegate接口的实现类
    statusDelegate, tokenDelegate, msgDelegate := createDelegates(appAccount)
    // 将三个实现类注册给用户
    mcUser.StatusDelegate(statusDelegate).TokenDelegate(tokenDelegate).MsgDelegate(msgDelegate).InitAndSetup()
    return mcUser
}
```
### 登录、收发消息、退出
```golang
// 创建用户
    leijun := createUser(&appId, &appAccount1)
    mifen := createUser(&appId, &appAccount2)

    // 用户登录
    leijun.Login()
    mifen.Login()
    mimc.Sleep(3000)

    // 互发消息
    leijun --> mifen
    leijun.SendMessage(appAccount2, []byte("Are you OK?"))
    leijun.SendMessage(appAccount2, []byte("Are you Okay?"))
    leijun.SendMessage(appAccount2, []byte("R U OK?"))
    
    mifen --> leijun
    mifen.SendMessage(appAccount2, []byte("I am Fine. Thanks!"))
    mifen.SendMessage(appAccount2, []byte("I'm Fine. Thanks!"))
    mifen.SendMessage(appAccount2, []byte("i m fine. thx!"))
    mimc.Sleep(3000)

    // 用户退出
    leijun.Logout()
    mifen.Logout()

    mimc.Sleep(3000)
```

[回到目录](#开发者文档)