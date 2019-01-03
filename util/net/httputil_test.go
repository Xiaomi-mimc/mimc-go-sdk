package httputil

import (
	"bytes"
	//"container/list"
	"encoding/json"
	"fmt"
	"github.com/Xiaomi-mimc/mimc-go-sdk/common/constant"
	"reflect"
	"testing"
)

func TestGetAddressFromResolver(t *testing.T) {
	url := cnst.ONLINE_RESOLVER_URL
	list := "app.chat.xiaomi.net"
	response := GetAddressFromResolver(url, list)
	var jsonData = bytes.NewBufferString(*response).Bytes()
	domain := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &domain)
	if err != nil {
		fmt.Printf("error:%v", err)
		return
	}
	fmt.Printf("--%s---\n", *response)
	fmt.Printf("\nS:%s", domain["S"])
	fmt.Printf("\nR:%s", domain["R"])
	fmt.Printf("\nR:%s\n", domain["R"].(map[string]interface{})["wifi"])
	felist := domain["R"].(map[string]interface{})["wifi"].(map[string]interface{})[list].([]interface{})
	//fe := felist.([]interface{})
	fmt.Printf("\n%v", felist[0].(string))
	fmt.Printf("\n%v\n", reflect.TypeOf(felist))
}

func TestHelloWorld(t *testing.T) {
	fmt.Printf("--%s---\n", "hello world")
}

func TestGetFEAddress(t *testing.T) {
	url := cnst.ONLINE_RESOLVER_URL
	list := "app.chat.xiaomi.net"
	response := GetFEAddress(url, list)
	if response == nil {
		fmt.Printf("request resolover failed!\n")
		return
	}
	fmt.Printf("%v\n", response[0])
}
