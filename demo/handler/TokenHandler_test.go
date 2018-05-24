package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func TestFetchToken(t *testing.T) {
	httpUrl := "http://10.38.162.149/api/account/token"
	appId := int64(2882303761517479657)
	apppKey := "5221747911657"
	appSecurt := "PtfBeZyC+H8SIM/UXhZx1w=="
	appAccount := "yusimin"

	tokenHandler := NewTokenHandler(&httpUrl, &apppKey, &appSecurt, &appAccount, &appId)
	tokenResponse := tokenHandler.FetchToken()
	fmt.Printf("%s\n", *tokenResponse)

	// 解析token
	var tokenMap map[string]interface{}
	if err := json.Unmarshal([]byte(*tokenResponse), &tokenMap); err == nil {
		data := tokenMap["data"].(map[string]interface{})
		code := tokenMap["code"].(float64)
		if code != 200 {
			return
		}
		appPackage := data["appPackage"].(string)
		chid := data["miChid"].(float64)
		uuid, err := strconv.ParseInt(data["miUserId"].(string), 10, 64)
		if err != nil {
			return
		}
		securityKey := data["miUserSecurityKey"].(string)
		token, ok := data["token"]
		if ok {
			tokenStr := token.(string)
			fmt.Printf("appPackage:%s\nchid:%v\nuuid:%v\nsecretKey:%s\ntoken:%s\n", appPackage, chid, uuid, securityKey, tokenStr)

		} else {
			fmt.Printf("token parse fail.\n")
		}
	} else {
		fmt.Printf("josn parse fail.\n")
	}
}
