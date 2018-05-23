## 开发者文档
* [1.概述](#1-概述)
* [2.下载源代码](#2-源代码)
> * [2.1 引入包](#2-1-引入包)
> * [2.2 常用接口](#2-2-常用接口)
* [3.引入第三方库](#3-引入第三方库)
* [4.开发者需实现的接口](#4-开发者需实现的接口)
> * [4.1 安全认证之token获取与解析](#4-1-安全认证之token获取与解析)
> * [4.2 用户在线状态回调](#4-2-用户在线状态回调)
> * [4.3 消息回调](#4-3-消息回调)
> > * [4.3.1 单聊消息回调](#4-3-1-单聊消息回调)
> > * [4.3.2 群聊消息接收回调](#4-3-2-群聊消息接收回调)
> > * [4.3.3 单聊消息超时回调](#4-3-3-单聊消息超时回调)
> > * [4.3.4 群聊消息超时回调](#4-3-4-群聊消息超时回调)
* [5.用户创建及初始化](#5-用户创建及初始化)
* [6.登录](#6-登录)
* [7.发送单聊消息](#7-发送单聊消息)
* [8.发送群聊消息](#8-发送群聊消息)
* [9.接收单聊与群聊消息](#9-接收单聊与群聊消息)
* [10.退出](#10-退出)
* [11.简单的例子](#11-简单的例子)

## 1-概述

MIMC, 全称XiaoMi Instant Messaging Cloud（小米即时消息云），为企业、个人服务或App提供快速接入即时通讯的功能。MIMC-Go-SDK提供了企业Go服务接入MIMC的能力。
## 2-源代码
MIMC-Go-SDK的源码已经上传到github中，开发者可以下载至本地GOPATH下安装。
``` sbtshell
    go get github.com/Xiaomi-mimc/mimc-go-sdk
    cd $GOPATH/github.com/Xiaomi-mimc/mimc-go-sdk
    go build
    go install
```
### 2-1-引入包
``` golang
    // 引入包路径, 该路径下引入了mimc包
    import "github.com/Xiaomi-mimc/mimc-go-sdk"
```
### 2-2-常用接口
引入MIMC-Go-SDK的包后，就可以与MIMC交互。其中最常用的API如下:
``` golang
    // 创建用户
    func NewUser(appId int64, appAccount string) *MCUser
    // 用户登录
    func (this *MCUser) Login() bool
    // 发送单聊消息
    func (this *MCUser) SendMessage(toAppAccount string, msgByte []byte) string
    // 发送群聊消息
    func (this *MCUser) SendGroupMessage(topicId *int64, msgByte []byte) string
    // 用户退出
    func (this *MCUser) Logout() bool
```
下面会对这些API使用的具体细节展开。

## 3-引入第三方库
MIMC-Go-SDK依赖proto buffer进行序列化与反序列化数据，在使用时，确保已经install了proto buffer。如未安装，可参考如下进行安装操作。
``` sbtshell
    go get github.com/golang/protobuf/proto
    // 进入下载目录
    cd $GOPATH/src/github.com/golang/protobuf/proto 
    // 编译安装
    go build
    go install
 ```
## 4-开发者需实现的接口
开发者服务通过SDK与MIMC交互时，需要处理以下问题：
* 安全相关的token获取与解析

token是SDK与MIMC服务进行交互的鉴权凭证。
安全相关的敏感信息(小米开放平台下发的AppKey,AppSecret)不应存储在SDK中，应由开发者服务存储。开发者服务器用AppId,AppKey, AppSecret, AppAccount向MIMC服务为用户申请Token，并返回给SDK。
* 用户在线状态回调

用户在线改变会触发相应业务逻辑，为此MIMC-Go-SDK提供了回调接口，供开发者实现其业务逻辑。

* 消息回调

消息回调主要是：1）单聊/群聊消息到达时，通知开发者对消息进行处理（如显示消息）；2）单聊/群聊消息发送超时时，通知开发者对超时消息处理（如重发超时消息）
### 4-1-安全认证之Token获取与解析
token是用户通过SDK与MIMC服务交互的凭证。开发者需要根据小米平台提供的AppId, AppKey, AppSecret, AppAccount为用户向MIMC服务请求token。
开发者需要实现下面接口，在实现中，可根据[MIMC-Auth安全认证](03-auth.md)文档实现请求Token的细节。
``` golang
    type Token interface {
        FetchToken() *string
    }
```
开发者实现上面接口后，应在创建用户时，将实例通过MCUser.registerTokenDelegate(tokenDeleget Token)方法实现**token回调**的注册。

Token接口实现示例如下：
``` golang
    type TokenHandler struct {
        httpUrl    string                       // 请求MIMC-Account服务的url
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
### 4-2-用户在线状态回调
SDK在用户状态发生改变时，会通过**用户在线状态回调**来通知开发者进行处理，**用户在线状态回调**的接口如下：
``` golang
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

开发者实现了上面接口后，在创建用户时，通过MCUser.registerStatusDelegate(statusDelegate StatusDelegate)方法实现**用户在线状态回调**的注册。

**StatusDelegate**接口实现示例如下：
``` golang
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

### 4-3-消息回调
当MIMC-Go-SDK收到消息时，会通过**消息回调**来执行开发者的业务逻辑，主要包括5种回调：
> 1) 单聊消息回调
> 2) 群聊消息回调
> 3) 单聊消息超时回调
> 4) 群聊消息超时回调
> 5) MIMC Server Ack回调

消息回调的接口如下：
``` golang
    type MessageHandlerDelegate interface {
        // 单聊消息回调
        HandleMessage(packets *list.List)
        // 群聊消息回调
        HandleGroupMessage(packets *list.List)
        // 处理服务器Sequence Ack的回调
        HandleServerAck(packetId *string, sequence, timestamp *int64)
        // 单聊消息超时回调
        HandleSendMessageTimeout(message *msg.P2PMessage)
        // 群聊消息超时回调
        HandleSendGroupMessageTimeout(message *msg.P2TMessage)
    }
```
在此接口中，单聊消息是P2PMessage的List集合，群聊消息是P2TMessage的List集合，二者结构如下：
``` golang
    package msg
    // 单聊消息
    type P2PMessage struct {
        packetId     *string
        sequence     *int64
        timestamp    *int64     
        fromAccount  *string
        fromResource *string
        payload      []byte
    }
    // 群聊消息
    type P2TMessage struct {
        packetId     *string
        sequence     *int64
        timestamp    *int64
        fromAccount  *string
        fromResource *string
        groupId      *int64
        payload      []byte
    }
```
#### 4-3-1-单聊消息回调
单聊消息回调的方法接口如下：
``` golang
    /**
     * @param[packets *list.List] 单聊消息集
     * @note: P2PMessage 单聊消息
     *  P2PMessage.pakcetId     数据包Id
     *  P2PMessage.sequence     消息序列号
     *  P2PMessage.timestamp    时间戳
     *  P2PMessage.fromAccount  发送者账号
     *  P2PMessage.fromResource 发送者设备
     *  P2PMessage.payload      消息体
     *
     **/
    HandleMessage(packets *list.List)
```
**HandleMessage**方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleMessage(packets *list.List) {
        for ele := packets.Front(); ele != nil; ele = ele.Next() {
            p2pmsg := ele.Value.(*msg.P2PMessage)
            logger.Info("[handle p2p msg]%v -> %v, pcktId: %v, timestamp: %v.", *(p2pmsg.FromAccount()), string(p2pmsg.Payload()), *(p2pmsg.PacketId()), *(p2pmsg.Timestamp()))
        }
    }
```
#### 4-3-2-群聊消息接收回调
群聊消息回调的方法接口如下：
``` golang
    /**
     * @param[packets *list.List] 群聊消息集
     * @note: P2TMessage 群聊消息
     *  P2TMessage.pakcetId     数据包Id
     *  P2TMessage.sequence     消息序列号
     *  P2TMessage.timestamp    时间戳
     *  P2TMessage.fromAccount  发送者账号
     *  P2TMessage.fromResource 发送者设备
     *  P2TMessage.groupId      群Id
     *  P2TMessage.payload      消息体
     *
     **/
    HandleGroupMessage(packets *list.List)
```
**HandleGroupMessage**方法实现示例如下：
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
#### 4-3-3-单聊消息超时回调
单聊消息超时回调的方法接口如下：
``` golang
    /**
    * @param[message *msg.P2PMessage] 单聊消息
    *  P2PMessage.pakcetId     数据包Id
    *  P2PMessage.sequence     消息序列号
    *  P2PMessage.timestamp    时间戳
    *  P2PMessage.fromAccount  发送者账号
    *  P2PMessage.fromResource 发送者设备
    *  P2PMessage.payload      消息体
    *
    **/
    HandleSendMessageTimeout(message *msg.P2PMessage)
```
**HandleSendMessageTimeout**方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleSendMessageTimeout(message *msg.P2PMessage) {
        logger.Info("[handle p2pmsg timeout] packetId:%v, msg:%v, time: %v.", *(message.PacketId()), string(message.Payload()), time.Now())
    }
```
#### 4-3-4-群聊消息超时回调
群聊消息超时回调的方法接口如下：
``` golang
    /**
     * @param[message *msg.P2TMessage] 群聊消息
     *  P2TMessage.pakcetId     数据包Id
     *  P2TMessage.sequence     消息序列号
     *  P2TMessage.timestamp    时间戳
     *  P2TMessage.fromAccount  发送者账号
     *  P2TMessage.fromResource 发送者设备
     *  P2TMessage.groupId      群Id
     *  P2TMessage.payload      消息体
     *
     **/
    HandleSendGroupMessageTimeout(message *msg.P2TMessage)
```
**HandleSendGroupMessageTimeout**方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleSendGroupMessageTimeout(message *msg.P2TMessage) {
        logger.Info("[handle p2tmsg timeout] packetId:%v, msg:%v.", *(message.PacketId()), string(message.Payload()))
    }
```
#### 4-3-5-MIMC Server Ack回调
MIMC Server Ack回调的方法接口如下：
``` golang
    /**
     * @param[packetId *string] 数据包Id
     * @param[sequence *int64] 消息序列号
     * @param[timestamp *int64] 时间戳
     */
    HandleServerAck(packetId *string, sequence, timestamp *int64)
```
**HandleServerAck**方法实现示例如下：
``` golang
    type MsgHandler struct {
    }
    func (this MsgHandler) HandleServerAck(packetId *string, sequence, timestamp *int64) {
        logger.Info("[handle server ack] packetId:%v, timestamp:%v.", *packetId, *timestamp)
    }
```
## 5-用户创建及初始化
用户登录前需要：创建、注册回调、初始化。具体包括如下：
> **1) 创建用户**
``` golang
    import "github.com/XiaoMi-mimc/mimc-go-sdk"
    // 创建用户需要appId，appAccount
    var appId int64 = int64(2882303761517479657)
    var appAccount string = "leijun"
    mcUser := mimc.NewUser(appId, appAccount)
```
> **2) 注册回调接口**
``` golang
    // 1. 创建回调接口实现类的实例，开发者实现回调接口，以及实现类的创建
    var tokenDelegate Token = createTokenDelegate()   // 创建Token接口实现类的实例
    var statusDelegate StatusDelegate = createStatusDelegate() // 创建StatusDelegate接口实现类的实例
    var messageDelegate MessageDelegate = createMessageDelegate() // 创建MessageDelegate接口实现类的实例
    // 2.注册回调接口
    mcUser.RegisterStatusDelegate(StatusDelegate).RegisterTokenDelegate(tokenDelegate).RegisterMessageDelegate(MessageDelegate)
```
> **3) 初始化**

最后需要为MCUser初始化以及启动读、写协程。
``` golang
    mcUser.InitAndSetup()
```
至此，用户就可以进行登录操作了。

## 6-登录
``` golang
    mcUser.Login()
```
## 7-发送单聊消息
发送单聊消息的接口如下：
``` golang
    /**
     * @param[toAppAccount string] 接收者
     * @param[msgByte []byte] 开发者自定义消息体
     * @return 客户端生成的数据包Id
     **/
    func (this *MCUser) SendMessage(toAppAccount string, msgByte []byte) string
```
发送单聊消息示例如下：
``` golang
    packetId := mcUser.SendMessage("MiFen", []byte("Are you OK?"))
```
## 8-发送群聊消息
发送群聊消息的借口如下：
```golang
    /**
     * @param[topicId *int64] 群Id
     * @param[msgByte []byte] 开发者自定义消息体
     * @return 客户端生成的数据包Id
     **/
    func (this *MCUser) SendGroupMessage(topicId *int64, msgByte []byte) string
```
发送群聊消息示例如下：
``` golang
    groupId := int64(123456789) 
    packetId := mcUser.SendMessage(&groupid, []byte("Are you OK?"))
```
## 9-接收单聊与群聊消息
接收单聊消息与群聊消息的业务逻辑在[4.3.1-单聊消息回调](#4.3.1-单聊消息回调)和[4.3.2-群聊消息回调](#4.3.2-群聊消息接收回调)中实现。
## 10-退出
``` golang
    mcUser.Logout()
```

## 11-简单的例子
``` golang

```

[回到目录](#开发者文档)
